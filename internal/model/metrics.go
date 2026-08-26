package model

import "time"

type MetricSeries struct {
	Name      string
	Values    []float64
	Timestamps []time.Time
}

func NewMetricSeries(name string) *MetricSeries {
	return &MetricSeries{Name: name}
}

func (m *MetricSeries) Append(value float64, at time.Time) {
	m.Values = append(m.Values, value)
	m.Timestamps = append(m.Timestamps, at)
}

func (m *MetricSeries) Last() (float64, bool) {
	if len(m.Values) == 0 {
		return 0, false
	}
	return m.Values[len(m.Values)-1], true
}

func (m *MetricSeries) Mean() float64 {
	if len(m.Values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range m.Values {
		total += value
	}
	return total / float64(len(m.Values))
}

func (m *MetricSeries) Reset() {
	m.Values = m.Values[:0]
	m.Timestamps = m.Timestamps[:0]
}
