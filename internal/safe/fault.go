package safe

import "windturbine/internal/model"

type FaultKind string

const (
	FaultNone      FaultKind = "none"
	FaultOverspeed FaultKind = "overspeed"
	FaultEmergency FaultKind = "emergency"
	FaultVibration FaultKind = "vibration"
)

func classifyFault(protection model.ProtectionState, activeAlarms int) FaultKind {
	if protection == model.ProtectionStop {
		return FaultEmergency
	}
	if protection == model.ProtectionOverspeed {
		return FaultOverspeed
	}
	if activeAlarms > 2 {
		return FaultVibration
	}
	return FaultNone
}

func (c *Controller) Fault() FaultKind {
	protection := c.state.Snapshot().Protection
	active := 0
	if c.alarm != nil {
		active = c.alarm.ActiveCount()
	}
	return classifyFault(protection, active)
}
