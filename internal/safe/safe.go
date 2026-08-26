package safe

import (
	"sync"

	"windturbine/internal/alarm"
	"windturbine/internal/core"
	"windturbine/internal/model"
)

type YawCoordinator interface {
	Abort() error
	IsYawing() bool
	ResetAbort()
}

type Controller struct {
	state   *core.State
	yaw     YawCoordinator
	alarm   *alarm.Service
	mu      sync.Mutex
	stopping bool
}

func NewController(state *core.State, yaw YawCoordinator, alarmSvc *alarm.Service) *Controller {
	return &Controller{state: state, yaw: yaw, alarm: alarmSvc}
}

func (c *Controller) IsStopped() bool {
	return c.state.Snapshot().Protection == model.ProtectionStop
}

func (c *Controller) IsStopping() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.stopping
}

func (c *Controller) State() *core.State {
	return c.state
}
