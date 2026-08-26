package main

import "windturbine/internal/model"

type metrics struct {
	WindSpeed      float64 `json:"wind_speed"`
	RotorSpeed     float64 `json:"rotor_speed"`
	ActivePower    float64 `json:"active_power"`
	PitchAngle     float64 `json:"pitch_angle"`
	YawAngle       float64 `json:"yaw_angle"`
	Protection     string  `json:"protection"`
	ActiveAlarms   int     `json:"active_alarms"`
	JournalHandles int     `json:"journal_handles"`
}

func (c *Controller) Metrics() metrics {
	status := c.Status()
	return metrics{
		WindSpeed:      c.sensors.AverageWind(),
		RotorSpeed:     status.RotorSpeed,
		ActivePower:    c.powerCurve.PowerAt(c.sensors.AverageWind()),
		PitchAngle:     status.PitchAngle,
		YawAngle:       status.YawAngle,
		Protection:     string(status.Protection),
		ActiveAlarms:   c.alarm.ActiveCount(),
		JournalHandles: c.journal.OpenHandleCount(),
	}
}

func (c *Controller) TurbineMetrics() model.TurbineMetrics {
	status := c.Status()
	return model.TurbineMetrics{
		WindSpeed:   c.sensors.AverageWind(),
		RotorSpeed:  status.RotorSpeed,
		ActivePower: c.powerCurve.PowerAt(c.sensors.AverageWind()),
		PitchAngle:  status.PitchAngle,
		YawAngle:    status.YawAngle,
	}
}
