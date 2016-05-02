package main

import (
	"encoding/json"
	"errors"
	"github.com/satori/go.uuid"
	"sync"
)

type Application struct {
	UniqueId string          `json:"id"`
	Name     string          `json:"name"`
	BaseUrl  string          `json:"baseUrl"`
	Services map[string]bool `json:"services"`
	State    State           `json:"state"`
	sync.RWMutex
}

func (a *Application) Id() string {
	return a.UniqueId
}

func (a *Application) ServicesSet() map[string]bool {

	a.RLock()
	defer a.RUnlock()

	result := make(map[string]bool, len(a.Services))

	for k, v := range a.Services {
		result[k] = v
	}

	return result

}

func (a *Application) ContainsService(id string) bool {
	a.RLock()
	defer a.RUnlock()
	return a.Services[id]
}

func (a *Application) AddService(service *Service) error {

	a.Lock()
	defer a.Unlock()

	if a.Services[service.UniqueId] {
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

	a.Services[service.UniqueId] = true

	return nil
}

func (a *Application) RemoveServiceId(id string) {

	a.Lock()
	defer a.Unlock()

	delete(a.Services, id)
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
		a.Services = make(map[string]bool)
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
