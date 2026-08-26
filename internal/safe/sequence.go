package safe

import "windturbine/internal/model"

func (c *Controller) SavePreStopSnapshot() error {
	snap := c.state.SnapshotForRecovery()
	snap.Protection = model.ProtectionNormal
	return c.alarm.SaveSnapshot(snap)
}

func (c *Controller) StopSequence() error {
	if err := c.SavePreStopSnapshot(); err != nil {
		return err
	}
	return c.Stop()
}

func (c *Controller) RecoverSequence() error {
	if err := c.Recover(); err != nil {
		return err
	}
	return nil
}
