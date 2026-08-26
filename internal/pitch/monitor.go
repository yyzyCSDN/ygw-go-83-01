package pitch

import "math"

func (c *Controller) BladeImbalance() float64 {
	angles := c.BladeAngles()
	minAngle := angles[0]
	maxAngle := angles[0]
	for _, angle := range angles {
		minAngle = math.Min(minAngle, angle)
		maxAngle = math.Max(maxAngle, angle)
	}
	return maxAngle - minAngle
}

func (c *Controller) BladesConsistent() bool {
	return c.BladeImbalance() < 0.01
}

func (c *Controller) AverageBladeAngle() float64 {
	angles := c.BladeAngles()
	total := 0.0
	for _, angle := range angles {
		total += angle
	}
	return total / float64(len(angles))
}
