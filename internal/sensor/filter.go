package sensor

type rollingAverage struct {
	window  int
	values  []float64
	cursor  int
	filled  bool
}

func newRollingAverage(window int) *rollingAverage {
	if window <= 0 {
		window = 8
	}
	return &rollingAverage{window: window, values: make([]float64, window)}
}

func (r *rollingAverage) add(value float64) float64 {
	r.values[r.cursor] = value
	r.cursor++
	if r.cursor >= r.window {
		r.cursor = 0
		r.filled = true
	}
	return r.current()
}

func (r *rollingAverage) current() float64 {
	count := r.cursor
	if r.filled {
		count = r.window
	}
	if count == 0 {
		return 0
	}
	total := 0.0
	for i := 0; i < count; i++ {
		total += r.values[i]
	}
	return total / float64(count)
}

func (c *Collector) SmoothWind(id string, window int) float64 {
	avg := newRollingAverage(window)
	sample := c.LastSample(id)
	if sample == nil || !sample.HasData {
		return 0
	}
	return avg.add(sample.WindSpeed)
}
