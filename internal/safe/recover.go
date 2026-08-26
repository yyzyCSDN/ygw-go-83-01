package safe

import (
	"windturbine/internal/model"
)

func (c *Controller) Recover() error {
	if !c.IsStopped() {
		return model.ErrNotStopping
	}
	snap, err := c.alarm.LatestSnapshot()
	if err != nil {
		return err
	}
	c.state.RestoreProtection(snap.Protection)
	if c.yaw != nil {
		c.yaw.ResetAbort()
	}
	if err := c.alarm.SaveSnapshot(model.RecoverySnapshot{
		Protection: model.ProtectionNormal,
	}); err != nil {
		return err
	}
	c.mu.Lock()
	c.stopping = false
	c.mu.Unlock()
	if c.alarm != nil {
		_ = c.alarm.Clear("safety-stop")
	}
	return nil
}
