package speed

type PowerPoint struct {
	WindSpeed float64
	Power     float64
}

type PowerCurve struct {
	points []PowerPoint
}

func DefaultPowerCurve() *PowerCurve {
	return &PowerCurve{
		points: []PowerPoint{
			{WindSpeed: 3.0, Power: 0},
			{WindSpeed: 4.0, Power: 40},
			{WindSpeed: 6.0, Power: 220},
			{WindSpeed: 8.0, Power: 620},
			{WindSpeed: 10.0, Power: 1250},
			{WindSpeed: 12.0, Power: 2100},
			{WindSpeed: 14.0, Power: 2900},
			{WindSpeed: 16.0, Power: 3300},
		},
	}
}

func (c *PowerCurve) PowerAt(windSpeed float64) float64 {
	if len(c.points) == 0 {
		return 0
	}
	if windSpeed <= c.points[0].WindSpeed {
		return 0
	}
	last := c.points[len(c.points)-1]
	if windSpeed >= last.WindSpeed {
		return last.Power
	}
	for i := 1; i < len(c.points); i++ {
		upper := c.points[i]
		lower := c.points[i-1]
		if windSpeed <= upper.WindSpeed {
			span := upper.WindSpeed - lower.WindSpeed
			ratio := (windSpeed - lower.WindSpeed) / span
			return lower.Power + ratio*(upper.Power-lower.Power)
		}
	}
	return last.Power
}

func (c *PowerCurve) CutInWindSpeed() float64 {
	if len(c.points) == 0 {
		return 0
	}
	return c.points[0].WindSpeed
}

func (c *PowerCurve) RatedWindSpeed() float64 {
	if len(c.points) == 0 {
		return 0
	}
	return c.points[len(c.points)-1].WindSpeed
}
