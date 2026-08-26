package yaw

import "windturbine/internal/model"

func (c *Controller) AutoAlign(windDirection float64, deadband float64) error {
	current := c.state.Snapshot().YawAngle
	target := c.TrackWind(windDirection)
	if withinDeadband(current, target, deadband) {
		c.commitYaw(target, model.YawAligned)
		return nil
	}
	return c.YawTo(target)
}

func (c *Controller) CurrentYawState() model.YawState {
	return c.state.Snapshot().YawState
}

func (c *Controller) CurrentYawAngle() float64 {
	return c.state.Snapshot().YawAngle
}

func (c *Controller) YawTarget(windDirection float64) float64 {
	return c.TrackWind(windDirection)
}
