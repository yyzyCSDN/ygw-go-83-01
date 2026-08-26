package model

type PowerReading struct {
	ActivePower float64
	ReactivePower float64
	Timestamp   int64
}

type TurbineMetrics struct {
	WindSpeed   float64
	RotorSpeed  float64
	ActivePower float64
	PitchAngle  float64
	YawAngle    float64
}

func (m TurbineMetrics) Available() bool {
	return m.WindSpeed >= 0 && m.RotorSpeed >= 0
}

func (m TurbineMetrics) ProductionFactor() float64 {
	if m.ActivePower <= 0 {
		return 0
	}
	return m.ActivePower / 3300.0
}
