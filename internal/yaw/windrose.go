package yaw

type WindRose struct {
	bins []float64
	size int
}

func NewWindRose(size int) *WindRose {
	if size < 2 {
		size = 2
	}
	return &WindRose{bins: make([]float64, size), size: size}
}

func (w *WindRose) Add(direction float64, weight float64) {
	index := int(normalizeDirection(direction) / 360.0 * float64(w.size))
	if index >= w.size {
		index = w.size - 1
	}
	w.bins[index] += weight
}

func (w *WindRose) DominantDirection() float64 {
	best := 0
	bestWeight := -1.0
	for i, weight := range w.bins {
		if weight > bestWeight {
			bestWeight = weight
			best = i
		}
	}
	return float64(best) * 360.0 / float64(w.size)
}

func (w *WindRose) Total() float64 {
	total := 0.0
	for _, weight := range w.bins {
		total += weight
	}
	return total
}
