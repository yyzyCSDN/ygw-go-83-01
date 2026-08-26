package safe

import "windturbine/internal/model"

func (c *Controller) ProtectionState() model.ProtectionState {
	return c.state.ProtectionState()
}

func (c *Controller) ShouldStop(protection model.ProtectionState) bool {
	return protection == model.ProtectionStop || protection == model.ProtectionOverspeed
}

func (c *Controller) StopReason() FaultKind {
	return c.Fault()
}

func (c *Controller) YawEngaged() bool {
	if c.yaw == nil {
		return false
	}
	return c.yaw.IsYawing()
}
