package pitch_test

import (
	"testing"

	"windturbine/internal/core"
	"windturbine/internal/pitch"
	"windturbine/internal/sensor"
	"windturbine/internal/speed"
)

func TestBladePitchConsistent(t *testing.T) {
	state := core.NewState()
	collect := sensor.NewCollector()
	collect.Feed("anemometer", 10, 15)
	table := speed.DefaultLimitTable()
	protector := speed.NewProtector(state, table, collect, "anemometer")
	protector.Evaluate()
	p := pitch.NewController(state, protector, nil)
	if err := p.Feather(45); err != nil {
		t.Fatal(err)
	}
	angles := p.BladeAngles()
	if angles[0] != angles[1] || angles[1] != angles[2] {
		t.Fatalf("blades inconsistent: %v", angles)
	}
}
