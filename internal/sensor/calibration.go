package sensor

type calibration struct {
	offset float64
	scale  float64
}

func (c calibration) apply(value float64) float64 {
	return value*c.scale + c.offset
}

type CalibratedCollector struct {
	inner *Collector
	cal   calibration
}

func NewCalibratedCollector(inner *Collector, offset, scale float64) *CalibratedCollector {
	return &CalibratedCollector{inner: inner, cal: calibration{offset: offset, scale: scale}}
}

func (c *CalibratedCollector) CalibratedWind(id string) float64 {
	sample := c.inner.LastSample(id)
	if sample == nil || !sample.HasData {
		return 0
	}
	return c.cal.apply(sample.WindSpeed)
}
