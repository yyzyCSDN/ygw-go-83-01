package yaw_test

import (
	"sync"
	"testing"

	"windturbine/internal/core"
	"windturbine/internal/pitch"
	"windturbine/internal/yaw"
)

type b01Speed struct{}

func (b01Speed) CommandedPitch() float64 { return 45 }
func (b01Speed) Generation() int         { return 1 }
func (b01Speed) ClampPitch(a float64) float64 { return a }
func (b01Speed) RatedRotorSpeed() float64 { return 16 }

type b01Alarm struct{}

func (b01Alarm) Raise(id, level, message string) error { return nil }

func TestPitchYawStateNoOverwrite(t *testing.T) {
	for i := 0; i < 500; i++ {
		state := core.NewState()
		p := pitch.NewController(state, b01Speed{}, b01Alarm{})
		y := yaw.NewController(state, b01Alarm{})
		var wg sync.WaitGroup
		wg.Add(2)
		go func() { defer wg.Done(); _ = p.Feather(45) }()
		go func() { defer wg.Done(); _ = y.YawTo(180) }()
		wg.Wait()
		st := state.Snapshot()
		if st.PitchAngle != 45 || st.YawAngle != 180 {
			t.Fatalf("lost update at iteration %d: pitch=%v yaw=%v", i, st.PitchAngle, st.YawAngle)
		}
	}
}
