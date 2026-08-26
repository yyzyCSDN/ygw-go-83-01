package sensor_test

import (
	"testing"

	"windturbine/internal/core"
	"windturbine/internal/model"
	"windturbine/internal/sensor"
	"windturbine/internal/speed"
)

func TestEmptySensorNoNilPanic(t *testing.T) {
	collect := sensor.NewCollector()
	collect.Register("anemometer")
	collect.MarkUnhealthy("anemometer")
	state := core.NewState()
	table := speed.DefaultLimitTable()
	protector := speed.NewProtector(state, table, collect, "anemometer")
	rotor := protector.RotorSpeed()
	if rotor != 0 {
		t.Fatalf("expected 0 for empty sensor, got %v", rotor)
	}
	if protection := protector.Evaluate(); protection != model.ProtectionNormal {
		t.Fatalf("expected normal, got %v", protection)
	}
}
