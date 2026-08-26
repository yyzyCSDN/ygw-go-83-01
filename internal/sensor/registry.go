package sensor

import "sort"

type SensorStatus struct {
	ID       string
	Healthy  bool
	WindSpeed float64
}

func (c *Collector) Registry() []SensorStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]SensorStatus, 0, len(c.sensors))
	for _, s := range c.sensors {
		result = append(result, SensorStatus{
			ID:        s.ID,
			Healthy:   s.Healthy,
			WindSpeed: s.WindSpeed,
		})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}

func (c *Collector) HealthySensorCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	count := 0
	for _, s := range c.sensors {
		if s.Healthy {
			count++
		}
	}
	return count
}
