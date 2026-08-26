package yaw

const (
	defaultDeadband = 2.0
	minDeadband     = 0.5
	maxDeadband     = 10.0
)

func clampDeadband(value float64) float64 {
	if value < minDeadband {
		return minDeadband
	}
	if value > maxDeadband {
		return maxDeadband
	}
	return value
}

func (c *Controller) Deadband() float64 {
	return defaultDeadband
}

func (c *Controller) Align(windDirection float64) error {
	return c.AutoAlign(windDirection, clampDeadband(defaultDeadband))
}

func (c *Controller) NeedsYaw(windDirection float64) bool {
	current := c.state.Snapshot().YawAngle
	target := c.TrackWind(windDirection)
	return !withinDeadband(current, target, defaultDeadband)
}
