package yaw

import "windturbine/internal/model"

func (c *Controller) CanYaw() bool {
	return c.state.Snapshot().Protection != model.ProtectionStop
}

func (c *Controller) YawIdle() bool {
	return c.state.Snapshot().YawState == model.YawIdle
}
