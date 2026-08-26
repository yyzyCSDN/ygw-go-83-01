package yaw

import "math"

func normalizeDirection(angle float64) float64 {
	for angle < 0 {
		angle += 360
	}
	for angle >= 360 {
		angle -= 360
	}
	return angle
}

func shortestYaw(from, to float64) float64 {
	delta := normalizeDirection(to - from)
	if delta > 180 {
		delta -= 360
	}
	return delta
}

func (c *Controller) TrackWind(windDirection float64) float64 {
	current := c.state.Snapshot().YawAngle
	target := normalizeDirection(windDirection)
	return normalizeDirection(current + shortestYaw(current, target))
}

func withinDeadband(from, to, deadband float64) bool {
	return math.Abs(shortestYaw(from, to)) <= deadband
}
