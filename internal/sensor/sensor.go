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

// LastSample 取回指定测风通道的最新采样。通道未注册、不健康或尚无数据时
// 返回空读数（HasData=false 的无效采样）而非 nil，调用方据此跳过本轮评估，
// 保护协程不能因为一次空读数解引用 nil 而崩溃。
func (c *Collector) LastSample(id string) *model.WindSample {
	c.mu.RLock()
	defer c.mu.RUnlock()
	s, ok := c.sensors[id]
	if !ok {
		return &model.WindSample{}
	}
	if !s.Healthy {
		return &model.WindSample{}
	}
	sample := &model.WindSample{
		WindSpeed:  s.WindSpeed,
		RotorSpeed: s.RotorSpeed,
		Timestamp:  s.LastSeen,
		HasData:    true,
	}
	return sample
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
