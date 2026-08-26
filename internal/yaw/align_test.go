package yaw

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"windturbine/internal/core"
	"windturbine/internal/model"
)

// neverAligns is an AlignmentSensor that never reports alignment, forcing the
// aligner to time out. It mirrors the field failure that leaves the nacelle
// stuck mid-alignment.
type neverAligns struct{}

func (neverAligns) Aligned(angle float64) bool { return false }

// stubAlarm records every Raise call so tests can assert the timeout alarm.
type stubAlarm struct {
	mu     sync.Mutex
	raise  []string
}

func (s *stubAlarm) Raise(id, level, message string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.raise = append(s.raise, id)
	return nil
}

func (s *stubAlarm) raised(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, r := range s.raise {
		if r == id {
			return true
		}
	}
	return false
}

// newTimeoutController wires a controller whose alignment probe never aligns,
// so YawTo always hits the 12s-equivalent timeout. The timeout is shortened
// only to keep the test fast.
func newTimeoutController() (*Controller, *stubAlarm) {
	state := core.NewState()
	alarm := &stubAlarm{}
	return NewControllerWithAlignment(state, alarm, neverAligns{}, 50*time.Millisecond), alarm
}

// TestYawToTimeoutResetsToIdle reproduces the field fault: when alignment times
// out the yaw state must fall back to idle so the turbine can re-align on the
// next tick instead of staying wedged in "yawing".
func TestYawToTimeoutResetsToIdle(t *testing.T) {
	c, alarm := newTimeoutController()

	err := c.YawTo(90)
	if !errors.Is(err, model.ErrYawTimeout) {
		t.Fatalf("expected ErrYawTimeout, got %v", err)
	}

	if c.IsYawing() {
		t.Fatalf("yawing flag still set after timeout; turbine cannot re-align")
	}
	if got := c.CurrentYawState(); got != model.YawIdle {
		t.Fatalf("YawState = %q, want %q (reset to re-alignable)", got, model.YawIdle)
	}
	if !alarm.raised("yaw-timeout") {
		t.Fatalf("expected yaw-timeout alarm to be raised")
	}

	// A subsequent alignment must be able to proceed from the reset state.
	go func() { _ = c.YawTo(120) }()
	time.Sleep(20 * time.Millisecond)
	if !c.IsYawing() {
		t.Fatalf("turbine could not re-enter yawing after timeout; stuck at idle permanently")
	}
	c.Abort()
}

// TestAbortCancelsInProgressAlignment verifies that an abort (safety stop)
// cancels the blocking alignment wait rather than letting it run to the full
// timeout — the wait that the field report says "was not cancelled".
func TestAbortCancelsInProgressAlignment(t *testing.T) {
	state := core.NewState()
	alarm := &stubAlarm{}
	// Long timeout so the only way out is the abort path.
	c := NewControllerWithAlignment(state, alarm, neverAligns{}, 30*time.Second)

	done := make(chan error, 1)
	go func() { done <- c.YawTo(45) }()

	// Give the goroutine time to enter the blocking wait.
	if !waitFor(func() bool { return c.IsYawing() }, time.Second) {
		t.Fatalf("turbine never entered yawing state")
	}

	start := time.Now()
	if err := c.Abort(); err != nil {
		t.Fatalf("Abort returned error: %v", err)
	}
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled after abort, got %v", err)
		}
		if elapsed := time.Since(start); elapsed > time.Second {
			t.Fatalf("abort did not cancel the wait; took %v", elapsed)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("abort did not cancel the blocking alignment wait")
	}

	if got := c.CurrentYawState(); got != model.YawIdle {
		t.Fatalf("YawState = %q after abort, want %q", got, model.YawIdle)
	}
}

func waitFor(cond func() bool, max time.Duration) bool {
	deadline := time.Now().Add(max)
	for time.Now().Before(deadline) {
		if cond() {
			return true
		}
		time.Sleep(time.Millisecond)
	}
	return cond()
}

// TestYawToAligns commits to aligned on success, guarding against the reset
// path running when nothing failed.
func TestYawToAligns(t *testing.T) {
	state := core.NewState()
	c := NewControllerWithAlignment(state, &stubAlarm{}, &alignProbe{aligned: true}, time.Second)

	if err := c.YawTo(60); err != nil {
		t.Fatalf("YawTo aligned failed: %v", err)
	}
	if got := c.CurrentYawState(); got != model.YawAligned {
		t.Fatalf("YawState = %q, want %q", got, model.YawAligned)
	}
}
