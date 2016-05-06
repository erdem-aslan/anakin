package main

import (
	"github.com/hashicorp/logutils"
	"hash/fnv"
	"io"
	logger "log"
	"os"
	"os/signal"
	"net"
)

type State string

const (
	Active    State = "active"
	Passive   State = "passive"
	Suspended State = "suspended"
	Failing   State = "failing"
	Trying    State = "trying"

	Version           string = "1.0"
	DefaultCfgPath    string = "cfg"
	DefaultCfgFile    string = "anakin.toml"
	DefaultDbPath     string = "db"
	DefaultDbFileName string = "anakin.db"

	DefaultAdminIp   string = ""
	DefaultAdminPort int    = 16016

	DefaultProxyIp       string = ""
	DefaultProxyPort     int    = 16015
	DefaultProxyRootPath string = "/"
)

var (
	config        *Configuration
	anakinCluster *AnakinCluster
	store         Store
	registry      *Registry
	log           *logger.Logger
	filter        *logutils.LevelFilter
	stats         *StatsContainer
)

func init() {
	filter = &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("WARN"),
		Writer:   os.Stderr,
	}

	log = logger.New(io.MultiWriter(os.Stdout), "<anakin>", logger.Ldate|logger.Ltime|logger.Lshortfile)

}

func main() {

	attachShutdownHook()

	// Config init does not return error, just exits the runtime.
	initConfig()

	err := initStore()

	if err != nil {
		log.Fatal("Store access has failed with error: ", err)
	}

	// init clustering
	err = initCluster()

	if err != nil {
		log.Fatal("Clustering has failed: ", err)
	}

	registry, err = initRegistry()

	if err != nil {
		log.Fatal("Registry initialization has failed with error: ", err)
	}

	store.AddListener(registry)

	initStats()

	go serveAdminBackend()

	err = serveProxy(registry)

	if err != nil {
		log.Fatal("Failed serving reverse proxy, error: ", err)
	}
}

func initStore() error {

	if config.MongoServers == nil {
		store = NewFsStore()
		return store.Initialize(config.DbPath + SEPARATOR + config.DbFileName)
	}

	store = NewMongoStore()
	return store.Initialize("")

}

func initRegistry() (r *Registry, err error) {
	r = NewRegistry()
	err = r.Init(store)

	return
}

func initStats() {
	stats = NewStatsContainer()
	stats.Start()
}

func attachShutdownHook() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	go func() {
		<-c
		if anakinCluster != nil {
			anakinCluster.Shutdown(true)
		}

		os.Exit(0)
	}()
}

func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32())
}

// Helper wrapper type which implements sort.Interface
type SortByDESCLength []string

func (a SortByDESCLength) Len() int {
	return len(a)
}
func (a SortByDESCLength) Less(i, j int) bool {
	return len(a[i]) > len(a[j])
}
func (a SortByDESCLength) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

