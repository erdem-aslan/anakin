package main

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"gopkg.in/mgo.v2"
	"sync"
	"time"
)

var (
	AlreadyPresentError    error = errors.New("Already present")
	MissingEntryError      error = errors.New("Missing entry")
	InvalidStoreStateError error = errors.New("Invalid store state")
)

type StoreState int

const (
	Created StoreState = iota
	Initialized
	Shutdown
)

const DB_NAME string = "Anakin"
const APPS_COL string = "apps"
const SERVICES_COL string = "services"
const ENDPOINTS_COL string = "endpoints"
const STATS_COL string = "stats"
const USERS_COL string = "users"

type Store interface {
	State() StoreState
	Initialize(dbFilePath string) (err error)

	AddListener(listener StoreListener) error
	RemoveListener(listener StoreListener)

	CreateApplication(a *Application) (err error)
	DeleteApplication(id string) error
	GetApplication(id string) (*Application, error)
	GetApplications() ([]*Application, error)
	UpdateApplication(a *Application) error

	CreateService(s *Service) (err error)
	DeleteService(id string) error
	GetService(id string) (*Service, error)
	GetServices() ([]*Service, error)
	UpdateService(s *Service) error

	CreateEndpoint(e *Endpoint) (err error)
	DeleteEndpoint(id string) error
	GetEndpoint(id string) (*Endpoint, error)
	GetEndpoints() ([]*Endpoint, error)
	UpdateEndpoint(e *Endpoint) error

	Shutdown() (bool, error)
}

type Entity interface {
	Id() string
}

type StoreListener interface {
	ApplicationAdded(a *Application)
	ApplicationUpdated(a *Application)
	ApplicationRemoved(id string)
	ServiceAdded(s *Service)
	ServiceUpdated(s *Service)
	ServiceRemoved(id string)
	EndpointAdded(e *Endpoint)
	EndpointUpdated(e *Endpoint)
	EndpointRemoved(id string)
}

type MongoStore struct {
	state   StoreState
	l       map[StoreListener]bool
	ll      sync.RWMutex
	session *mgo.Session
	sync.RWMutex
}

func (ms *MongoStore) State() StoreState {
	return ms.state
}

func (ms *MongoStore) Initialize(ignored string) (err error) {

	ms.Lock()
	defer ms.Unlock()
	log.Println("Initializing mongo store...")

	if ms.state == Initialized {
		return InvalidStoreStateError
	}

	s, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    config.MongoServers,
		Direct:   false,
		Timeout:  5 * time.Second,
		FailFast: true,
		Database: DB_NAME,
	})

	if err != nil {
		return err
	}

	ms.session = s

	s.SetMode(mgo.Monotonic, true)

	ms.state = Initialized

	log.Println("Initializing mongo store, finished")

	return nil

}

func (ms *MongoStore) CreateApplication(a *Application) error {

	err := ms.create(a, APPS_COL)

	if err == nil {
		ms.ll.RLock()
		for lis, _ := range ms.l {
			go lis.ApplicationAdded(a)
		}
		ms.ll.RUnlock()

		broadcastEvent(AppCreated, a.Id())

	}

	return err
}

func (ms *MongoStore) GetApplication(id string) (*Application, error) {
	app := &Application{UniqueId: id}
	app.Init()
	err := ms.get(app, APPS_COL)

	return app, err

}

func (ms *MongoStore) UpdateApplication(a *Application) error {

	err := ms.update(a, APPS_COL)

	if err == nil {
		ms.ll.RLock()
		for lis, _ := range ms.l {
			go lis.ApplicationUpdated(a)
		}
		ms.ll.RUnlock()

		broadcastEvent(AppUpdated, a.Id())
	}

	return err
}

func (ms *MongoStore) DeleteApplication(id string) error {

	a, err := ms.GetApplication(id)

	if a == nil {
		return MissingEntryError
	}

	if err != nil {
		return err
	}

	for serviceId, _ := range a.Services {
		ms.DeleteService(serviceId)
	}

	err = ms.delete(id, APPS_COL)

	if err == nil {
		ms.ll.RLock()
		for lis, _ := range ms.l {
			go lis.ApplicationRemoved(id)
		}
		ms.ll.RUnlock()

		broadcastEvent(AppDeleted, id)
	}

	return err
}

func (ms *MongoStore) GetApplications() ([]*Application, error) {

	apps := make([]*Application, 10)

	scopy := ms.session.Copy()
	defer scopy.Close()
	err := scopy.DB(DB_NAME).C(APPS_COL).Find(nil).All(&apps)

	return apps, err
}

func (ms *MongoStore) CreateService(s *Service) (err error) {

	err = ms.create(s, SERVICES_COL)

	if err == nil {
		ms.ll.RLock()
		for lis, _ := range ms.l {
			go lis.ServiceAdded(s)
		}
		ms.ll.RUnlock()

		broadcastEvent(SvcCreated, s.Id())
	}

	return err

}

func (ms *MongoStore) DeleteService(id string) error {

	s, err := ms.GetService(id)

	if s == nil {
		return MissingEntryError
	}

	for endpointId, _ := range s.Endpoints {
		ms.DeleteEndpoint(endpointId)
	}

	err = ms.delete(id, SERVICES_COL)

	if err == nil {
		ms.ll.RLock()
		for lis, _ := range ms.l {
			go lis.ServiceRemoved(id)
		}
		ms.ll.RUnlock()

		broadcastEvent(SvcDeleted, id)

	}

	return err

}

func (ms *MongoStore) GetService(id string) (*Service, error) {
	service := &Service{UniqueId: id}
	err := ms.get(service, SERVICES_COL)

	return service, err
}

func (ms *MongoStore) GetServices() ([]*Service, error) {
	services := make([]*Service, 10)

	scopy := ms.session.Copy()
	defer scopy.Close()
	err := scopy.DB(DB_NAME).C(SERVICES_COL).Find(nil).All(&services)

	return services, err

}

func (ms *MongoStore) UpdateService(s *Service) error {
	err := ms.update(s, SERVICES_COL)

	if err == nil {
		ms.ll.RLock()
		for lis, _ := range ms.l {
			go lis.ServiceUpdated(s)
		}
		ms.ll.RUnlock()

		broadcastEvent(SvcUpdated, s.Id())
	}

	return err

}

func (ms *MongoStore) CreateEndpoint(e *Endpoint) (err error) {

	err = ms.create(e, ENDPOINTS_COL)

	if err == nil {
		ms.ll.RLock()
		for lis, _ := range ms.l {
			go lis.EndpointAdded(e)
		}
		ms.ll.RUnlock()

		broadcastEvent(EndpCreated, e.Id())

	}

	return err

}
func (ms *MongoStore) DeleteEndpoint(id string) error {

	err := ms.delete(id, ENDPOINTS_COL)

	if err == nil {
		ms.ll.RLock()
		for lis, _ := range ms.l {
			go lis.EndpointRemoved(id)
		}
		ms.ll.RUnlock()

		broadcastEvent(EndpDeleted, id)
	}

	return err

}
func (ms *MongoStore) GetEndpoint(id string) (*Endpoint, error) {
	endpoint := &Endpoint{UniqueId: id}
	err := ms.get(endpoint, ENDPOINTS_COL)

	return endpoint, err

}
func (ms *MongoStore) GetEndpoints() ([]*Endpoint, error) {

	endpoints := make([]*Endpoint, 10)

	scopy := ms.session.Copy()
	defer scopy.Close()
	err := scopy.DB(DB_NAME).C(ENDPOINTS_COL).Find(nil).All(&endpoints)

	return endpoints, err

}
func (ms *MongoStore) UpdateEndpoint(e *Endpoint) error {

	err := ms.update(e, ENDPOINTS_COL)

	if err == nil {
		ms.ll.RLock()
		for lis, _ := range ms.l {
			go lis.EndpointUpdated(e)
		}
		ms.ll.RUnlock()

		broadcastEvent(EndpUpdated, e.Id())
	}

	return err
}

func (ms *MongoStore) AddListener(listener StoreListener) error {

	if ms.l[listener] {
		return errors.New("Listener already present")
	}

	ms.l[listener] = true
	return nil
}

func (ms *MongoStore) RemoveListener(listener StoreListener) {
	delete(ms.l, listener)
}

func (ms *MongoStore) Shutdown() (bool, error) {

	if ms.state != Initialized {
		return false, errors.New("Store is not initialized")
	}

	ms.session.Fsync(false)
	ms.session.Close()

	ms.state = Shutdown

	return true, nil
}

func (ms *MongoStore) create(entity Entity, collection string) (err error) {
	scopy := ms.session.Copy()
	defer scopy.Close()
	err = scopy.DB(DB_NAME).C(collection).Insert(entity)
	return
}

func (ms *MongoStore) delete(entityId string, collection string) (err error) {
	scopy := ms.session.Copy()
	defer scopy.Close()
	err = scopy.DB(DB_NAME).C(collection).RemoveId(entityId)
	return
}

func (ms *MongoStore) get(entity Entity, collection string) (err error) {
	scopy := ms.session.Copy()
	defer scopy.Close()
	err = scopy.DB(DB_NAME).C(collection).FindId(entity.Id()).One(entity)
	return
}

func (ms *MongoStore) update(entity Entity, collection string) (err error) {
	scopy := ms.session.Copy()
	defer scopy.Close()
	info, err := scopy.DB(DB_NAME).C(collection).Upsert(entity.Id(), entity)

	if info.Updated != 1 {
		log.Println("Store update has failed for entity:", entity)
	}

	return

}

func NewFsStore() Store {

	return &FsStore{
		state: Created,
		l:     make(map[StoreListener]bool),
	}
}

func NewMongoStore() Store {

	return &MongoStore{
		state: Created,
		l:     make(map[StoreListener]bool),
	}
}

type FsStore struct {
	db    *bolt.DB
	state StoreState
	l     map[StoreListener]bool
	ll    sync.RWMutex
	sync.RWMutex
}

func (fs *FsStore) State() StoreState {
	return fs.state
}

func (fs *FsStore) Initialize(dbFilePath string) (err error) {

	fs.Lock()
	defer fs.Unlock()
	log.Println("Initializing embedded store...")

	if fs.state == Initialized {
		return InvalidStoreStateError
	}

	fs.db, err = bolt.Open(dbFilePath,
		0600,
		&bolt.Options{Timeout: 2 * time.Second})

	if err != nil {
		return err
	}

	err = fs.db.Update(func(tx *bolt.Tx) error {

		_, err = tx.CreateBucketIfNotExists([]byte(APPS_COL))

		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(SERVICES_COL))

		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(ENDPOINTS_COL))

		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(STATS_COL))

		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(USERS_COL))

		if err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		return err
	}

	fs.state = Initialized

	log.Println("Initializing embedded store, finished")

	return nil

}

// -- apps

func (fs *FsStore) CreateApplication(a *Application) error {

	err := fs.create(a, APPS_COL)

	if err == nil {
		fs.ll.RLock()
		for lis, _ := range fs.l {
			go lis.ApplicationAdded(a)
		}
		fs.ll.RUnlock()

		broadcastEvent(AppCreated, a.Id())
	}

	return err

}

func (fs *FsStore) DeleteApplication(id string) error {

	a, err := fs.GetApplication(id)

	if a == nil {
		return MissingEntryError
	}

	if err != nil {
		return err
	}

	for serviceId, _ := range a.Services {
		fs.DeleteService(serviceId)
	}

	err = fs.delete(id, APPS_COL)

	if err == nil {
		fs.ll.RLock()
		for lis, _ := range fs.l {
			go lis.ApplicationRemoved(id)
		}
		fs.ll.RUnlock()

		broadcastEvent(AppDeleted, id)
	}

	return err
}

func (fs *FsStore) UpdateApplication(a *Application) error {

	err := fs.update(a, APPS_COL)

	if err == nil {
		fs.ll.RLock()
		for lis, _ := range fs.l {
			go lis.ApplicationUpdated(a)
		}
		fs.ll.RUnlock()

		broadcastEvent(AppUpdated, a.Id())
	}

	return err

}

func (fs *FsStore) GetApplication(id string) (*Application, error) {

	app := &Application{UniqueId: id}
	app.Init()
	err := fs.get(app, APPS_COL)

	return app, err

}

func (fs *FsStore) GetApplications() ([]*Application, error) {
	return fs.getAllApps()
}

// -- services

func (fs *FsStore) CreateService(s *Service) error {

	err := fs.create(s, SERVICES_COL)

	if err == nil {
		fs.ll.RLock()
		for lis, _ := range fs.l {
			go lis.ServiceAdded(s)
		}
		fs.ll.RUnlock()

		broadcastEvent(SvcCreated, s.Id())
	}

	return err
}

func (fs *FsStore) DeleteService(id string) error {

	s, err := fs.GetService(id)

	if s == nil {
		return MissingEntryError
	}

	for endpointId, _ := range s.Endpoints {
		fs.DeleteEndpoint(endpointId)
	}

	err = fs.delete(id, SERVICES_COL)

	if err == nil {
		fs.ll.RLock()
		for lis, _ := range fs.l {
			go lis.ServiceRemoved(id)
		}
		fs.ll.RUnlock()

		broadcastEvent(SvcDeleted, id)
	}

	return err
}

func (fs *FsStore) UpdateService(s *Service) error {
	err := fs.update(s, SERVICES_COL)

	if err == nil {
		fs.ll.RLock()
		for lis, _ := range fs.l {
			go lis.ServiceUpdated(s)
		}
		fs.ll.RUnlock()

		broadcastEvent(SvcUpdated, s.Id())
	}

	return err
}

func (fs *FsStore) GetService(id string) (*Service, error) {

	service := &Service{UniqueId: id}
	err := fs.get(service, SERVICES_COL)

	return service, err

}

func (fs *FsStore) GetServices() ([]*Service, error) {
	return fs.getAllServices()
}

// -- endpoints

func (fs *FsStore) CreateEndpoint(e *Endpoint) error {

	err := fs.create(e, ENDPOINTS_COL)

	if err == nil {
		fs.ll.RLock()
		for lis, _ := range fs.l {
			go lis.EndpointAdded(e)
		}
		fs.ll.RUnlock()

		broadcastEvent(EndpCreated, e.Id())
	}

	return err
}

func (fs *FsStore) DeleteEndpoint(id string) error {

	err := fs.delete(id, ENDPOINTS_COL)

	if err == nil {
		fs.ll.RLock()
		for lis, _ := range fs.l {
			go lis.EndpointRemoved(id)
		}
		fs.ll.RUnlock()

		broadcastEvent(EndpDeleted, id)
	}

	return err
}

func (fs *FsStore) UpdateEndpoint(e *Endpoint) error {
	err := fs.update(e, ENDPOINTS_COL)

	if err == nil {
		fs.ll.RLock()
		for lis, _ := range fs.l {
			go lis.EndpointUpdated(e)
		}
		fs.ll.RUnlock()

		broadcastEvent(EndpUpdated, e.Id())
	}

	return err
}

func (fs *FsStore) GetEndpoint(id string) (*Endpoint, error) {

	endpoint := &Endpoint{UniqueId: id}
	err := fs.get(endpoint, ENDPOINTS_COL)

	return endpoint, err

}

func (fs *FsStore) GetEndpoints() ([]*Endpoint, error) {

	return fs.getAllEndpoints()
}

func (fs *FsStore) create(entity Entity, bucket string) (err error) {

	err = fs.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(bucket))

		buf, err := json.Marshal(entity)

		if err != nil {
			return err
		}

		if buf == nil {
			return errors.New("Failed marshalling application")
		}

		err = b.Put([]byte(entity.Id()), buf)

		return err
	})

	return nil
}

func (fs *FsStore) delete(entityId string, bucket string) (err error) {

	fs.Lock()
	defer fs.Unlock()

	err = fs.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(bucket))

		return b.Delete([]byte(entityId))
	})

	return err

}

func (fs *FsStore) get(entity Entity, bucket string) (err error) {

	err = fs.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(bucket))

		appBytes := b.Get([]byte(entity.Id()))

		if appBytes != nil {

			err := json.Unmarshal(appBytes, entity)

			if err != nil {
				return err
			}
		}

		return nil
	})

	return err

}

func (fs *FsStore) getAllApps() (apps []*Application, err error) {

	apps = make([]*Application, 0, 10)

	err = fs.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(APPS_COL))

		if b.Stats().KeyN == 0 {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			var temp *Application = &Application{}

			err := json.Unmarshal(v, temp)

			if err != nil {
				return err
			}
			apps = append(apps, temp)
		}

		return nil
	})

	return

}

func (fs *FsStore) getAllServices() (services []*Service, err error) {

	services = make([]*Service, 0, 10)

	err = fs.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(SERVICES_COL))

		if b.Stats().KeyN == 0 {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			var temp *Service = &Service{}

			err := json.Unmarshal(v, temp)

			if err != nil {
				return err
			}
			services = append(services, temp)
		}

		return nil
	})

	return

}

func (fs *FsStore) getAllEndpoints() (endpoints []*Endpoint, err error) {

	endpoints = make([]*Endpoint, 0, 10)

	err = fs.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(ENDPOINTS_COL))

		if b.Stats().KeyN == 0 {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			var temp *Endpoint = &Endpoint{}

			err := json.Unmarshal(v, temp)

			if err != nil {
				return err
			}
			endpoints = append(endpoints, temp)
		}

		return nil
	})

	return

}

func (fs *FsStore) update(entity Entity, bucket string) error {

	err := fs.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(bucket))

		v := b.Get([]byte(entity.Id()))

		if v == nil {
			return MissingEntryError
		}

		v, err := json.Marshal(entity)

		if err != nil {
			return err
		}

		err = b.Put([]byte(entity.Id()), v)

		if err != nil {
			return err
		}

		return nil
	})

	return err

}

func (fs *FsStore) Shutdown() (bool, error) {

	if fs.state != Initialized {
		return false, errors.New("Store is not initialized")
	}

	err := fs.db.Sync()

	if err != nil {
		return false, err
	}

	err = fs.db.Close()

	if err != nil {
		return false, err
	}

	fs.state = Shutdown

	return true, nil
}

func (fs *FsStore) AddListener(listener StoreListener) error {

	if fs.l[listener] {
		return errors.New("Listener already present")
	}

	fs.l[listener] = true
	return nil
}

func (fs *FsStore) RemoveListener(listener StoreListener) {
	delete(fs.l, listener)
}

func broadcastEvent(eventType EventType, id string) {

	err := anakinCluster.BroadcastAnakinEvent(&ClusterEvent{
		EventType: eventType,
		Payload:   id,
	})

	if err != nil {
		log.Println("Failed notifiying the cluster, error: ", err)
	}
}
