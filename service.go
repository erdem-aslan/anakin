package main

import (
	"container/ring"
	"encoding/json"
	"errors"
	"github.com/satori/go.uuid"
	"sync"
)

type BalanceStrategy string

const (
	Round_Robin    BalanceStrategy = "round-robin"
	Source_Hashing BalanceStrategy = "source-hashing"
)

type Service struct {
	UniqueId     string          `json:"id"`
	ServiceUrl   string          `json:"serviceUrl"`
	Endpoints    map[string]bool `json:"endpoints"`
	Tps          int             `json:"tps"`
	CurrentState State           `json:"state"`

	sync.RWMutex `json:"-"`

	BalanceStrategy BalanceStrategy `json:"balanceStrategy"`
	serviceEPRing   *ring.Ring
	serviceEPList   []*Endpoint
	slb             sync.RWMutex
}

func (s *Service) Id() string {
	return s.UniqueId
}

func (s *Service) EndpointsSet() map[string]bool {

	s.RLock()
	defer s.RUnlock()

	result := make(map[string]bool, len(s.Endpoints))

	for k, v := range s.Endpoints {
		result[k] = v
	}

	return result

}

func (s *Service) AddEndpoint(id string) error {

	s.Lock()
	defer s.Unlock()

	if s.Endpoints[id] {
		return AlreadyPresentError
	}

	s.Endpoints[id] = true

	return nil
}

func (s *Service) RemoveEndpoint(id string) {

	s.Lock()
	defer s.Unlock()

	delete(s.Endpoints, id)
}

func (s *Service) State() State {
	s.RLock()
	defer s.RUnlock()
	return s.CurrentState
}

func (s *Service) SetState(state State) {
	s.Lock()
	defer s.Unlock()
	s.CurrentState = state
}

func (s *Service) Init() error {

	s.UniqueId = uuid.NewV4().String()

	if s.CurrentState == "" {
		s.SetState(Active)
	}

	if s.BalanceStrategy == "" {
		s.BalanceStrategy = Round_Robin
	}

	if s.Endpoints == nil {
		s.Endpoints = make(map[string]bool)
	}

	if s.ServiceUrl == "" {
		return errors.New("Missing serviceUrl definition")
	}

	return nil

}

func (s *Service) String() string {
	ss, err := json.Marshal(s)

	if err != nil {
		return err.Error()
	}

	return string(ss)

}
