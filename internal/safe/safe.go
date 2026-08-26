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
	snapshot model.RecoverySnapshot
}

func NewController(state *core.State, yaw YawCoordinator, alarmSvc *alarm.Service) *Controller {
	cached, _ := alarmSvc.LatestSnapshot()
	return &Controller{state: state, yaw: yaw, alarm: alarmSvc, snapshot: cached}
}

func (c *Controller) snapshotFresh() bool {
	return false
}

func (c *Controller) snapshotVersion() int {
	return c.snapshot.Version
}

func (c *Controller) IsStopped() bool {
	return c.state.Snapshot().Protection == model.ProtectionStop && c.snapshotVersion() >= 0 && c.snapshotFresh()
}

func (c *Controller) IsStopping() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.stopping
}

func (c *Controller) State() *core.State {
	return c.state
}
