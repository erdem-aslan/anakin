package main

import (
	"encoding/json"
	"errors"
	"github.com/satori/go.uuid"
	"sort"
	"sync"
)

type Application struct {
	UniqueId       string            `json:"id" bson:"_id,omitempty"`
	Name           string            `json:"name" bson:"name"`
	BaseUrl        string            `json:"baseUrl" bson:"baseUrl"`
	Services       map[string]string `json:"services" bson:"services"`
	State          State             `json:"state" bson:"state"`
	servicesSorted []string
	sync.RWMutex   `json:"-" bson:"-"`
}

func (a *Application) Id() string {
	return a.UniqueId
}

func (a *Application) ServicesCopy() map[string]string {

	a.RLock()
	defer a.RUnlock()

	result := make(map[string]string, len(a.Services))

	for k, v := range a.Services {
		result[k] = v
	}

	return result

}

func (a *Application) ServicesSorted() []string {

	a.RLock()

	if a.servicesSorted != nil {
		a.RUnlock()
		return a.servicesSorted
	}

	a.Lock()
	defer a.Unlock()
	a.sortServices()

	return a.servicesSorted

}

func (a *Application) ContainsService(id string) bool {
	a.RLock()
	defer a.RUnlock()
	return a.Services[id] != ""
}

func (a *Application) AddService(service *Service) error {

	a.Lock()
	defer a.Unlock()

	if a.Services[service.UniqueId] != "" {
		return AlreadyPresentError
	}

	for s, _ := range a.Services {
		svc, err := store.GetService(s)

		if err != nil {
			return InvalidStoreStateError
		}

		if svc.ServiceUrl == service.ServiceUrl {
			return AlreadyPresentError
		}
	}

	a.Services[service.UniqueId] = service.ServiceUrl
	a.sortServices()

	return nil
}

func (a *Application) RemoveServiceId(id string) {

	a.Lock()
	defer a.Unlock()

	delete(a.Services, id)
	a.sortServices()
}

func (a *Application) sortServices() {
	a.servicesSorted = make([]string, len(a.Services))

	for _, v := range a.Services {
		a.servicesSorted = append(a.servicesSorted, v)
	}

	sort.Sort(SortByDESCLength(a.servicesSorted))

}

func (a *Application) SetState(state State) {
	a.Lock()
	defer a.Unlock()
	a.State = state
}

func (a *Application) Init() error {

	if a.Name == "" || a.BaseUrl == "" {
		return errors.New("Missing name and/or baseUrl definition")
	}

	if a.UniqueId == "" {
		a.UniqueId = uuid.NewV4().String()
	}

	if a.State == "" {
		a.SetState(Active)
	}

	if a.Services == nil {
		a.Services = make(map[string]string)
	}

	return nil
}

func (a *Application) String() string {
	j, err := json.Marshal(a)

	if err != nil {
		return err.Error()
	}
	return string(j)
}
