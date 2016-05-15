package main

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"sync"
	"time"
)

type Endpoint struct {
	UniqueId     string      `json:"id" bson:"_id,omitempty"`
	Host         string      `json:"host"`
	Port         string      `json:"port"`
	Scheme       string      `json:"scheme"`
	CurrentState State       `json:"state" bson:"state"`
	failedCount  int         `json:"-"`
	failedTimer  *time.Timer `json:"-"`
	sync.RWMutex `json:"-" bson:"-"`
}

func (e *Endpoint) String() string {
	j, err := json.Marshal(e)

	if err != nil {
		return err.Error()
	}
	return string(j)
}

func (e *Endpoint) Id() string {
	return e.UniqueId
}

func (e *Endpoint) Address() string {
	return e.Host + ":" + e.Port
}

func (e *Endpoint) State() State {
	e.RLock()
	defer e.RUnlock()
	return e.CurrentState
}

func (e *Endpoint) SetState(state State) {
	e.Lock()
	defer e.Unlock()

	if e.CurrentState != state && state == Active {
		log.Println("Endpoint is active back again: ", e)
	}

	e.CurrentState = state

	if state == Failing {

		time.AfterFunc(30*time.Second, func() {

			e.Lock()

			if e.failedCount >= 5 {
				e.CurrentState = Suspended
				log.Println("Endpoint is suspended until further notice: ", e)
			} else {
				e.CurrentState = Trying
				e.failedCount = 0
			}

			e.failedCount++

			e.Unlock()

			store.UpdateEndpoint(e)

		})

		log.Println("Scheduled state update for endpointId: ", e)
	}

}

func (e *Endpoint) Init() {
	e.UniqueId = uuid.NewV4().String()

	if e.CurrentState == "" {
		e.SetState(Active)
	}
}
