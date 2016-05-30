package main

import (
	"fmt"
	"github.com/rcrowley/go-metrics"
	"runtime"
	"strings"
	"sync"
	"time"
)

func NewStatsContainer() *StatsContainer {

	return &StatsContainer{
		instanceStats: &InstanceStats{
			Os:       fmt.Sprint(runtime.GOOS, "-", runtime.GOARCH),
			CpuCores: runtime.NumCPU(),
		},
		appStats:   make(map[string]*ApplicationStats),
		svcStats:   make(map[string]*ServiceStats),
		endpStats:  make(map[string]*EndpointStats),
		reg:        metrics.NewRegistry(),
		countCh:    make(chan struct{}),
		resourceCh: make(chan struct{}),
	}
}

type StatsContainer struct {
	instanceStats  *InstanceStats
	appStats       map[string]*ApplicationStats
	svcStats       map[string]*ServiceStats
	endpStats      map[string]*EndpointStats
	countTicker    *time.Ticker
	resourceTicker *time.Ticker
	countCh        chan struct{}
	resourceCh     chan struct{}
	reg            metrics.Registry
	sync.RWMutex
}

func (s *StatsContainer) InstanceStats() InstanceStats {
	s.RLock()
	defer s.RUnlock()
	i := *s.instanceStats
	return i
}

func (s *StatsContainer) UnregisterApplication(name string) {
	s.reg.Unregister("countApp__" + name)
	s.reg.Unregister("meterApp__" + name)
}

func (s *StatsContainer) UnregisterService(id string) {
	s.reg.Unregister("countService__" + id)
	s.reg.Unregister("meterService__" + id)
}

func (s *StatsContainer) UnregisterEndpoint(id string) {
	s.reg.Unregister("countEndpoint__" + id)
	s.reg.Unregister("meterEndpoint__" + id)
}

func (s *StatsContainer) IncrementApp(name string) {

	c := s.reg.GetOrRegister("countApp__"+name, func() metrics.Counter {
		//s.reg.Register("meterApp__"+name, metrics.NewMeter())
		return metrics.NewCounter()
	})

	c.(metrics.Counter).Inc(1)

	log.Println("Stats --- app hit: ", name)

}

func (s *StatsContainer) IncrementService(id string) {

	c := s.reg.GetOrRegister("countService__"+id, func() metrics.Counter {
		//s.reg.Register("meterService__"+id, metrics.NewMeter())
		return metrics.NewCounter()
	})

	c.(metrics.Counter).Inc(1)

	log.Println("Stats --- service hit: ", id)

}

func (s *StatsContainer) IncrementEndpoint(name string) {

	c := s.reg.GetOrRegister("countEndpoint__"+name, func() metrics.Counter {
		return metrics.NewCounter()
	})

	c.(metrics.Counter).Inc(1)

	log.Println("Stats --- endpoint hit: ", name)

}

func (s *StatsContainer) Start() {

	s.countTicker = time.NewTicker(time.Second)

	go func() {
		for {
			select {
			case <-s.countTicker.C:
				var totalMeter float64

				s.reg.Each(func(key string, value interface{}) {

					p := strings.Split(key, "__")

					if len(p) != 2 {
						log.Println("Stats naming error for key: ", s)
						return
					}

					if !strings.HasPrefix(p[0], "count") {
						return
					}

					counter := value.(metrics.Counter)
					count := counter.Count()
					counter.Clear()

					switch p[0] {
					case "countApp":
						meter := s.reg.GetOrRegister("meterApp__"+p[1], func() metrics.Meter {
							return metrics.NewMeter()
						}).(metrics.Meter)
						meter.Mark(count)
						s.Lock()

						st := s.appStats[p[1]]

						if st == nil {
							st = &ApplicationStats{
								Name: p[1],
							}
							s.appStats[p[1]] = st
						}

						s.appStats[p[1]].Rps = meter.Rate1()
						totalMeter = totalMeter + s.appStats[p[1]].Rps
						s.Unlock()
					case "countService":
						meter := s.reg.GetOrRegister("meterService__"+p[1], func() metrics.Meter {
							return metrics.NewMeter()
						}).(metrics.Meter)
						meter.Mark(count)
						s.Lock()

						st := s.svcStats[p[1]]

						if st == nil {
							st = &ServiceStats{
								Id: p[1],
							}
							s.svcStats[p[1]] = st
						}

						s.svcStats[p[1]].Rps = meter.Rate1()
						s.Unlock()
					case "countEndpoint":
						meter := s.reg.GetOrRegister("meterEndpoint__"+p[1], func() metrics.Meter {
							return metrics.NewMeter()
						}).(metrics.Meter)
						meter.Mark(count)
						s.Lock()

						st := s.endpStats[p[1]]

						if st == nil {
							st = &EndpointStats{
								Id: p[1],
							}
							s.endpStats[p[1]] = st
						}

						s.endpStats[p[1]].Rps = meter.Rate1()

						s.Unlock()
					default:
						log.Println("Unknown stats entry:", key)
						return
					}
				})

				s.Lock()
				s.instanceStats.Rps = totalMeter
				s.Unlock()

			case <-s.countCh:
				return
			}
		}
	}()

	s.resourceTicker = time.NewTicker(time.Second * 10)

	go func() {
		for {
			select {
			case <-s.resourceTicker.C:
				memory, memoryUnit := s.resolveMemoryUsage()
				s.Lock()
				s.instanceStats.Mem = fmt.Sprintf("%.2f %s", memory, memoryUnit)
				s.Unlock()

			case <-s.resourceCh:
				return
			}
		}

	}()
}

func (s *StatsContainer) resolveMemoryUsage() (memory float32, memoryUnit string) {

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	memory = float32(m.Alloc) / 1024
	memoryUnit = "kb"

	if memory > 1024 {
		memory = memory / float32(1024)
		memoryUnit = "mb"

		if memory > 1024 {
			memory = memory / float32(1024)
			memoryUnit = "gb"
		}
	}
	return
}

func (s *StatsContainer) Stop() {

	s.reg.UnregisterAll()

	s.countTicker.Stop()
	s.countCh <- struct{}{}
	s.countTicker = nil

	s.resourceTicker.Stop()
	s.resourceCh <- struct{}{}
	s.resourceTicker = nil
}

type InstanceStats struct {
	Os       string  `json:"os"`
	CpuCores int     `json:"cpuCores"`
	Mem      string  `json:"mem"`
	Rps      float64 `json:"rps"`
}

type ApplicationStats struct {
	Name string  `json:"name"`
	Rps  float64 `json:"rps"`
}

type ServiceStats struct {
	Id  string  `json:"id"`
	Rps float64 `json:"rps"`
}

type EndpointStats struct {
	Id  string  `json:"id"`
	Rps float64 `json:"rps"`
}
