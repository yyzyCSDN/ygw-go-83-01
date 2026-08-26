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

	mu       sync.Mutex
	yawing   bool
	aborting bool
	probe    AlignmentSensor
	timeout  time.Duration
}

func NewController(state *core.State, alarm AlarmSink) *Controller {
	return NewControllerWithAlignment(state, alarm, &alignProbe{aligned: true}, 12*time.Second)
}

func NewControllerWithAlignment(state *core.State, alarm AlarmSink, probe AlignmentSensor, timeout time.Duration) *Controller {
	return &Controller{state: state, alarm: alarm, probe: probe, timeout: timeout}
}

func (c *Controller) YawTo(target float64) error {
	// During a safety stop the yaw must stand down: neither start a new
	// alignment nor keep one going, so the stop state is never overwritten.
	if c.state.Snapshot().Protection == model.ProtectionStop {
		return model.ErrYawAborted
	}

	c.mu.Lock()
	c.yawing = true
	c.mu.Unlock()

	c.commitYaw(target, model.YawYawing)
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	if err := c.waitAligned(ctx, target); err != nil {
		c.mu.Lock()
		c.yawing = false
		c.mu.Unlock()
		// An aborted yaw leaves the nacelle idle; other errors (e.g. timeout)
		// also release the actuator. commitYaw is skipped if a stop engaged
		// while we were waiting, so we never clobber the stop state.
		if !c.aborting {
			c.commitYaw(target, model.YawIdle)
		}
		return err
	}
	c.mu.Lock()
	c.yawing = false
	c.mu.Unlock()
	c.commitYaw(target, model.YawAligned)
	return nil
}

func (c *Controller) Abort() error {
	c.mu.Lock()
	c.aborting = true
	c.mu.Unlock()
	// Mark the nacelle idle without touching Protection: commitYaw only ever
	// updates the yaw fields, so a concurrent stop stays intact.
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

func (c *Controller) isAborting() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.aborting
}

// commitYaw updates only the yaw fields of the nacelle status. It never writes
// the whole status back, so it cannot overwrite a Protection=stop that a
// concurrent safety stop just committed. During a stop the yaw is dormant and
// this is a no-op.
func (c *Controller) commitYaw(angle float64, state model.YawState) {
	if c.state.Snapshot().Protection == model.ProtectionStop {
		return
	}
	c.state.Update(func(s *core.Status) {
		s.YawAngle = angle
		s.YawState = state
	})
}
