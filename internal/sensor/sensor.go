package sensor

import (
	"sort"
	"sync"
	"time"

	"windturbine/internal/model"
)

type Sensor struct {
	ID         string
	Healthy    bool
	WindSpeed  float64
	RotorSpeed float64
	LastSeen   time.Time
}

type Collector struct {
	mu      sync.RWMutex
	sensors map[string]*Sensor
	now     func() time.Time
}

func NewCollector() *Collector {
	return &Collector{
		sensors: make(map[string]*Sensor),
		now:     time.Now,
	}
}

func (c *Collector) Register(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sensors[id] = &Sensor{ID: id, Healthy: true}
}

func (c *Collector) MarkUnhealthy(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if s, ok := c.sensors[id]; ok {
		s.Healthy = false
	}
}

func (c *Collector) Feed(id string, windSpeed, rotorSpeed float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	s, ok := c.sensors[id]
	if !ok {
		s = &Sensor{ID: id, Healthy: true}
		c.sensors[id] = s
	}
	s.WindSpeed = windSpeed
	s.RotorSpeed = rotorSpeed
	s.LastSeen = c.now().UTC()
	s.Healthy = true
}

func (c *Collector) LastSample(id string) *model.WindSample {
	c.mu.RLock()
	defer c.mu.RUnlock()
	s, ok := c.sensors[id]
	if !ok || !s.Healthy {
		return &model.WindSample{HasData: false}
	}
	return &model.WindSample{
		WindSpeed:  s.WindSpeed,
		RotorSpeed: s.RotorSpeed,
		Timestamp:  s.LastSeen,
		HasData:    true,
	}
}

func (c *Collector) HealthyIDs() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	ids := make([]string, 0, len(c.sensors))
	for id, s := range c.sensors {
		if s.Healthy {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	return ids
}

func (c *Collector) AverageWind() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	total := 0.0
	count := 0
	for _, s := range c.sensors {
		if !s.Healthy {
			continue
		}
		total += s.WindSpeed
		count++
	}
	if count == 0 {
		return 0
	}
	return total / float64(count)
}
