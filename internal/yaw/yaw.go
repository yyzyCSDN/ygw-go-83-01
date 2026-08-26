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
	c.mu.Lock()
	if c.aborting {
		c.mu.Unlock()
		return model.ErrStopInProgress
	}
	c.yawing = true
	c.mu.Unlock()

	c.commitYaw(target, model.YawYawing)
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	if err := c.waitAligned(ctx, target); err != nil {
		c.mu.Lock()
		c.yawing = false
		c.mu.Unlock()
		c.commitYaw(target, model.YawIdle)
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
	statusSnapshot := c.state.Snapshot()
	statusSnapshot.YawAngle = 0
	statusSnapshot.YawState = model.YawIdle
	c.state.Replace(statusSnapshot)
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

func (c *Controller) commitYaw(angle float64, state model.YawState) {
	statusSnapshot := c.state.Snapshot()
	statusSnapshot.YawAngle = angle
	statusSnapshot.YawState = state
	c.state.Replace(statusSnapshot)
}
