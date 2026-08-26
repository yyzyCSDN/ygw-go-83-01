package pitch

import "math"

func rateLimit(current, target, maxStep float64) float64 {
	if maxStep <= 0 {
		return target
	}
	delta := target - current
	if math.Abs(delta) <= maxStep {
		return target
	}
	if delta > 0 {
		return current + maxStep
	}
	return current - maxStep
}

func (c *Controller) FeatherLimited(target float64, maxStep float64) error {
	current := c.state.Snapshot().PitchAngle
	limited := rateLimit(current, target, maxStep)
	return c.FeatherTo(limited)
}

func (c *Controller) MaxPitchAngle() float64 {
	return 90
}
