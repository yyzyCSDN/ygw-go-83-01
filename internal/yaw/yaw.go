package yaw

import (
	"context"
	"sync"
	"time"

	"windturbine/internal/core"
	"windturbine/internal/model"
)

type AlarmSink interface {
	Raise(id, level, message string) error
}

type Controller struct {
	state *core.State
	alarm AlarmSink

	mu          sync.Mutex
	yawing      bool
	aborting    bool
	probe       AlignmentSensor
	timeout     time.Duration
	cancelAlign context.CancelFunc
}

func NewController(state *core.State, alarm AlarmSink) *Controller {
	return NewControllerWithAlignment(state, alarm, &alignProbe{aligned: true}, 12*time.Second)
}

func NewControllerWithAlignment(state *core.State, alarm AlarmSink, probe AlignmentSensor, timeout time.Duration) *Controller {
	return &Controller{state: state, alarm: alarm, probe: probe, timeout: timeout}
}

func (c *Controller) YawTo(target float64) error {
	c.mu.Lock()
	if c.aborting {
		c.mu.Unlock()
		return model.ErrStopInProgress
	}
	c.yawing = true
	ctx, cancel := context.WithCancel(context.Background())
	c.cancelAlign = cancel
	c.mu.Unlock()
	defer cancel()

	current := c.state.Snapshot().YawAngle
	c.commitYaw(target, model.YawYawing)
	if err := c.waitForAlignment(ctx, target); err != nil {
		c.resetAlignment(current, err)
		return err
	}
	c.mu.Lock()
	c.yawing = false
	c.cancelAlign = nil
	c.mu.Unlock()
	c.commitYaw(target, model.YawAligned)
	return nil
}

func (c *Controller) Abort() error {
	c.mu.Lock()
	c.aborting = true
	if c.cancelAlign != nil {
		c.cancelAlign()
	}
	c.mu.Unlock()
	c.commitYaw(0, model.YawIdle)
	return nil
}

func (c *Controller) IsYawing() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.yawing
}

func (c *Controller) ResetAbort() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.aborting = false
}

func (c *Controller) waitForAlignment(ctx context.Context, target float64) error {
	return c.waitAligned(ctx, target)
}

// resetAlignment restores the turbine to a state where it can re-attempt yaw
// alignment after a wait failed (timeout, abort, or cancellation). The yaw
// state returns to idle, the pre-alignment yaw angle is kept so a subsequent
// Align can recompute the target from the actual nacelle position, and a
// warning alarm is raised when the failure was a timeout.
func (c *Controller) resetAlignment(currentAngle float64, err error) {
	c.mu.Lock()
	c.yawing = false
	c.cancelAlign = nil
	c.mu.Unlock()
	c.commitYaw(currentAngle, model.YawIdle)
	if err == model.ErrYawTimeout && c.alarm != nil {
		_ = c.alarm.Raise("yaw-timeout", "warning", "yaw alignment timed out")
	}
}

func (c *Controller) commitYaw(angle float64, state model.YawState) {
	c.state.Update(func(s *core.Status) {
		s.YawAngle = angle
		s.YawState = state
	})
}
