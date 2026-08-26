package sensor

import (
	"math"

	"windturbine/internal/model"
)

const (
	maxPlausibleWindSpeed  = 60.0
	maxPlausibleRotorSpeed = 25.0
)

func ValidateSample(sample model.WindSample) bool {
	if !sample.HasData {
		return false
	}
	if math.IsNaN(sample.WindSpeed) || math.IsInf(sample.WindSpeed, 0) {
		return false
	}
	if math.IsNaN(sample.RotorSpeed) || math.IsInf(sample.RotorSpeed, 0) {
		return false
	}
	if sample.WindSpeed < 0 || sample.WindSpeed > maxPlausibleWindSpeed {
		return false
	}
	if sample.RotorSpeed < 0 || sample.RotorSpeed > maxPlausibleRotorSpeed {
		return false
	}
	return true
}

func (c *Collector) Plausible(id string) bool {
	sample := c.LastSample(id)
	if sample == nil || !sample.HasData {
		return false
	}
	return ValidateSample(*sample)
}

func (c *Collector) SensorCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.sensors)
}
