package pitch

import "sync"

type Blade struct {
	ID int

	mu    sync.Mutex
	angle float64
}

func (b *Blade) Angle() float64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.angle
}

func (b *Blade) setAngle(angle float64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.angle = angle
}
