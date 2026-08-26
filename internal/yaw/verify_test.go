package yaw_test

import (
	"testing"
	"time"

	"windturbine/internal/core"
	"windturbine/internal/model"
	"windturbine/internal/yaw"
)

type b06Alarm struct{}

func (b06Alarm) Raise(id, level, message string) error { return nil }

type b06NeverAlign struct{}

func (b06NeverAlign) Aligned(angle float64) bool { return false }

func TestYawTimeoutRecovers(t *testing.T) {
	state := core.NewState()
	y := yaw.NewControllerWithAlignment(state, b06Alarm{}, b06NeverAlign{}, 50*time.Millisecond)
	if err := y.YawTo(90); err == nil {
		t.Fatal("expected timeout error")
	}
	if y.IsYawing() {
		t.Fatal("yaw should not be stuck yawing")
	}
	if y.CurrentYawState() != model.YawIdle {
		t.Fatalf("yaw should be idle after timeout, got %v", y.CurrentYawState())
	}
}
