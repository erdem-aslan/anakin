package main

import (
	"bytes"
	"encoding/gob"
	"github.com/gladmir/anakin/remote"
	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/serf/serf"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"os"
	"strconv"
	"time"
	"sync"
)

func initCluster() error {
	anakinCluster = newAnakinCluster()
	return anakinCluster.Start(false)
}

func newAnakinCluster() *AnakinCluster {
	return &AnakinCluster{
		ec: make(chan serf.Event),
		rc: make(map[string]remote.AnakinClient),
	}
}

type AnakinCluster struct {
	sf       *serf.Serf
	ec       chan serf.Event
	rc map[string]remote.AnakinClient
	rcl sync.RWMutex
	Name     string
	nameHash int
	started  time.Time
}

func (ac *AnakinCluster) Start(randomNodeName bool) error {

	ac.started = time.Now()
	sc := serf.DefaultConfig()
	sc.LogOutput = filter
	sc.MemberlistConfig.LogOutput = filter

	sc.EventCh = ac.ec

	if randomNodeName {
		sc.NodeName = uuid.NewV4().String()
	} else {
		h, err := os.Hostname()

		if err != nil {
			h = "unknown"
		}

		sc.NodeName = h + ":" + strconv.Itoa(config.ClusterPort)
	}

	ac.Name = sc.NodeName
	ac.nameHash = hash(sc.NodeName)

	if sc.Tags == nil {
		sc.Tags = make(map[string]string)
	}

	var grpcIp string

	if config.ClusterIp == "" {
		grpcIp = GetLocalIP()
	} else {
		grpcIp = config.ClusterIp
	}

	addr, err := ac.startGrpcService(grpcIp)

	if err != nil {
		return err
	}

	sc.Tags["grpcAddress"] = addr.String()

	sc.MemberlistConfig.AdvertiseAddr = config.ClusterIp
	sc.MemberlistConfig.AdvertisePort = config.ClusterPort

	sc.MemberlistConfig.BindAddr = config.ClusterIp
	sc.MemberlistConfig.BindPort = config.ClusterPort

	go ac.handleClusterEvents()


	s, err := serf.Create(sc)

	if err != nil {
		return err
	}

	ac.sf = s

	if len(config.ClusterMembers) != 0 {
		n, err := ac.sf.Join(config.ClusterMembers, true)

		if err != nil {
			log.Println("All configured member(s) is/are out of our reach, we are alone...")
		}

		if n == 1 {
			log.Println("Started anakin cluster, awaiting additional instances ...")
		} else {
			log.Println("Cluster join was succesful, number of anakin instances: ", n)
		}
	}

	return nil
}

func (ac *AnakinCluster) startGrpcService(grpcIp string) (net.Addr, error) {

	lis, err := net.Listen("tcp", grpcIp+":"+strconv.Itoa(config.ClusterServicesPort))

	if err != nil {
		return nil, err
	}

	s := grpc.NewServer()
	remote.RegisterAnakinServer(s, ac)

	go s.Serve(lis)

	log.Println("Started Anakin Grpc Service on ", lis.Addr())

	return lis.Addr(), nil
}

func (ac *AnakinCluster) GetInstance(ctx context.Context, r *remote.RpcRequest) (*remote.Instance, error) {
	return ac.fetchRemoteInstance()
}

func (ac *AnakinCluster) fetchRemoteInstance() (*remote.Instance, error) {
	li := ac.LocalInstance()

	var state remote.State

	switch li.State {
	case Active:
		state = remote.State_Active
	case Passive:
		state = remote.State_Passive
	case Failing:
		state = remote.State_Failing
	case Suspended:
		state = remote.State_Suspended
	case Trying:
		state = remote.State_Trying
	}

	log.Println("State: ", state)

	ss, err := ptypes.TimestampProto(li.Started)

	if err != nil {
		return nil, err
	}

	i := &remote.Instance{
		Id : li.Id,
		Version : li.Version,
		AdminPort: li.AdminPort,
		AdminIp: li.AdminIp,
		ProxyIp: li.ProxyIp,
		ProxyPort:li.ProxyPort,
		Started:ss,
		State:state,
		Stats: &remote.InstanceStats{
			li.Stats.Os,
			int32(li.Stats.CpuCores),
			li.Stats.Mem,
			li.Stats.Rps,
		},
	}

	return i, nil;

}

func (ac *AnakinCluster) Instances() []*remote.Instance {

	members := ac.sf.Members()

	instances := make([]*remote.Instance, 0, len(members))

	for _, member := range members {

		if member.Name == ac.Name {
			local, err := ac.fetchRemoteInstance()

			if err != nil {
				log.Println(err)
			}

			instances = append(instances, local)
			continue
		}

		ac.rcl.Lock()

		c := ac.rc[member.Name]

		if c == nil {
			conn, err := grpc.Dial(member.Tags["grpcAddress"], grpc.WithInsecure())

			if err != nil {
				log.Println(err)
			}

			c = remote.NewAnakinClient(conn)
		}

		ac.rc[member.Name] = c
		ac.rcl.Unlock()

		now, err := ptypes.TimestampProto(time.Now())

		if err != nil {
			log.Println(err)
			continue
		}

		i, err := c.GetInstance(context.Background(), &remote.RpcRequest{
			now,
		})

		if err != nil {
			log.Println(err)
		}
		instances = append(instances, i)
	}

	return instances
}

func (ac *AnakinCluster) LocalInstance() *Instance {

	return &Instance{
		ac.Name,
		Version,
		strconv.Itoa(config.AdminPort),
		config.AdminIp,
		config.ProxyIp,
		strconv.Itoa(config.ProxyPort),
		ac.started,
		Active,
		stats.InstanceStats(),
	}

}

func (ac *AnakinCluster) BroadcastAnakinEvent(e *ClusterEvent) error {

	e.Sender = ac.nameHash
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	err := enc.Encode(e)

	if err != nil {
		return err
	}

	err = ac.sf.UserEvent("ace", buffer.Bytes(), true)

	if err != nil {
		return err
	}

	return nil
}

func (ac *AnakinCluster) Shutdown(waitForIt bool) error {

	var err error

	err = ac.sf.Leave()

	if err != nil {
		return err
	}

	err = ac.sf.Shutdown()

	if err != nil {
		return err
	}

	if waitForIt {
		<-ac.sf.ShutdownCh()
	}

	return nil

}

func (ac *AnakinCluster) handleClusterEvents() {

loop:
	for {
		select {
		case event, ok := <-ac.ec:

			if event != nil {
				switch event.EventType() {

				case serf.EventMemberJoin,
					serf.EventMemberLeave,
					serf.EventMemberFailed,
					serf.EventMemberUpdate,
					serf.EventMemberReap:
					go ac.handleMemberEvent(event.(serf.MemberEvent))

				case serf.EventUser:
					go ac.handleUserEvent(event.(serf.UserEvent))

				case serf.EventQuery:
					go ac.handleQuery(event.(*serf.Query))

				}
			}

			if !ok {
				log.Println("ClusterEventHandler exiting...")
				break loop
			}

		}
	}

}

func (ac *AnakinCluster) handleMemberEvent(e serf.MemberEvent) {

	for _, m := range e.Members {
		if m.Name == ac.Name {
			return
		}
	}

	switch e.EventType() {
	case serf.EventMemberJoin:
		log.Println("Anakin instance(s) joined the cluster, ", e.Members)
	case serf.EventMemberLeave:
		log.Println("Anakin instance(s) has left the cluster, ", e.Members)
	case serf.EventMemberFailed:
		log.Println("Anakin instance(s) failing: ", e.Members)
	case serf.EventMemberUpdate:
		log.Println("Anakin instance(s) updated: ", e.Members)
	case serf.EventMemberReap:
		log.Println("Anakin instance(s) reaped: ", e.Members)
	}

}

func (ac *AnakinCluster) handleUserEvent(e serf.UserEvent) {

	if e.Name != "ace" {
		log.Println("Received non anakin related event, ignoring..., name: ", e.Name)
		return
	}

	dec := gob.NewDecoder(bytes.NewReader(e.Payload))

	var m ClusterEvent

	err := dec.Decode(&m)

	if err != nil {
		log.Println("Failed decoding anakin event, error: ", err)
		return
	}

	if m.Sender == ac.nameHash {
		return
	}

	registry.RemoteRegistryEvent(m)
}

func (ac *AnakinCluster) handleQuery(e *serf.Query) {
}

type EventType int

const (
	AppCreated EventType = iota
	AppDeleted
	AppUpdated
	SvcCreated
	SvcDeleted
	SvcUpdated
	EndpCreated
	EndpUpdated
	EndpDeleted
)

type ClusterEvent struct {
	Sender    int
	EventType EventType
	Payload   string
}

type Instance struct {
	Id        string        `json:"id"`
	Version   string        `json:"version"`
	AdminPort string        `json:"adminPort"`
	AdminIp   string        `json:"adminIp"`
	ProxyIp   string        `json:"proxyIp"`
	ProxyPort string        `json:"proxyPort"`
	Started   time.Time     `json:"started"`
	State     State         `json:"state"`
	Stats     InstanceStats `json:"stats,omitempty"`
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
