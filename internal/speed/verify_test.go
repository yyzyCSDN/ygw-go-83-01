package speed_test

import (
	"testing"

	"windturbine/internal/core"
	"windturbine/internal/pitch"
	"windturbine/internal/sensor"
	"windturbine/internal/speed"
)

func TestSpeedCacheInvalidatedOnState(t *testing.T) {
	state := core.NewState()
	collect := sensor.NewCollector()
	collect.Register("anemometer")
	table := speed.DefaultLimitTable()
	protector := speed.NewProtector(state, table, collect, "anemometer")
	p := pitch.NewController(state, protector, nil)
	collect.Feed("anemometer", 5, 5)
	protector.Evaluate()
	low := p.CommandedPitch()
	collect.Feed("anemometer", 18, 18)
	protector.Evaluate()
	high := p.CommandedPitch()
	if high <= low {
		t.Fatalf("cache not invalidated: low=%v high=%v", low, high)
	}
}
