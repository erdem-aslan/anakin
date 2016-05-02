package main

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"log"
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

func NewFsStore() Store {

	return &FsStore{
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
	log.Println("Initializing store...")

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

		_, err = tx.CreateBucketIfNotExists([]byte("apps"))

		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("services"))

		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("endpoints"))

		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("stats"))

		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("users"))

		if err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		return err
	}

	fs.state = Initialized

	log.Println("Initializing store, finished")

	return nil

}

// -- apps

func (fs *FsStore) CreateApplication(a *Application) error {

	err := fs.create(a, "apps")

	if err == nil {
		fs.ll.RLock()
		for lis, _ := range fs.l {
			go lis.ApplicationAdded(a)
		}
		fs.ll.RUnlock()
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

	err = fs.delete(id, "apps")

	if err == nil {
		fs.ll.RLock()
		for lis, _ := range fs.l {
			go lis.ApplicationRemoved(id)
		}
		fs.ll.RUnlock()
	}

	return err
}

func (fs *FsStore) UpdateApplication(a *Application) error {

	err := fs.update(a, "apps")

	if err == nil {
		fs.ll.RLock()
		for lis, _ := range fs.l {
			go lis.ApplicationUpdated(a)
		}
		fs.ll.RUnlock()
	}

	return err

}

func (fs *FsStore) GetApplication(id string) (*Application, error) {

	app := &Application{UniqueId: id}
	app.Init()
	err := fs.get(app, "apps")

	return app, err

}

func (fs *FsStore) GetApplications() ([]*Application, error) {
	return fs.getAllApps()
}

// -- services

func (fs *FsStore) CreateService(s *Service) error {

	err := fs.create(s, "services")

	if err == nil {
		fs.ll.RLock()
		for lis, _ := range fs.l {
			go lis.ServiceAdded(s)
		}
		fs.ll.RUnlock()
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

	err = fs.delete(id, "services")

	if err == nil {
		fs.ll.RLock()
		for lis, _ := range fs.l {
			go lis.ServiceRemoved(id)
		}
		fs.ll.RUnlock()
	}

	return err
}

func (fs *FsStore) UpdateService(s *Service) error {
	err := fs.update(s, "services")

	if err == nil {
		fs.ll.RLock()
		for lis, _ := range fs.l {
			go lis.ServiceUpdated(s)
		}
		fs.ll.RUnlock()
	}

	return err
}

func (fs *FsStore) GetService(id string) (*Service, error) {

	service := &Service{UniqueId: id}
	err := fs.get(service, "services")

	return service, err

}

func (fs *FsStore) GetServices() ([]*Service, error) {
	return fs.getAllServices()
}

// -- endpoints

func (fs *FsStore) CreateEndpoint(e *Endpoint) error {

	err := fs.create(e, "endpoints")

	if err == nil {
		fs.ll.RLock()
		for lis, _ := range fs.l {
			go lis.EndpointAdded(e)
		}
		fs.ll.RUnlock()
	}

	return err
}

func (fs *FsStore) DeleteEndpoint(id string) error {

	err := fs.delete(id, "endpoints")

	if err == nil {
		fs.ll.RLock()
		for lis, _ := range fs.l {
			go lis.EndpointRemoved(id)
		}
		fs.ll.RUnlock()
	}

	return err
}

func (fs *FsStore) UpdateEndpoint(e *Endpoint) error {
	err := fs.update(e, "endpoints")

	if err == nil {
		fs.ll.RLock()
		for lis, _ := range fs.l {
			go lis.EndpointUpdated(e)
		}
		fs.ll.RUnlock()
	}

	return err
}

func (fs *FsStore) GetEndpoint(id string) (*Endpoint, error) {

	endpoint := &Endpoint{UniqueId: id}
	err := fs.get(endpoint, "endpoints")

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

		b := tx.Bucket([]byte("apps"))

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

		b := tx.Bucket([]byte("services"))

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

		b := tx.Bucket([]byte("endpoints"))

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
	fs.l[listener] = true
	return nil
}

func (fs *FsStore) RemoveListener(listener StoreListener) {
	delete(fs.l, listener)
}
