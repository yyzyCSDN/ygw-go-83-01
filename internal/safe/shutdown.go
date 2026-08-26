package safe

import "windturbine/internal/model"

func (c *Controller) EmergencyShutdown() error {
	if err := c.SavePreStopSnapshot(); err != nil {
		return err
	}
	if err := c.Stop(); err != nil && err != model.ErrStopInProgress {
		return err
	}
	if c.alarm != nil {
		_ = c.alarm.Raise("emergency-stop", "critical", "emergency shutdown")
	}
	return nil
}

func (c *Controller) NormalShutdown() error {
	if err := c.SavePreStopSnapshot(); err != nil {
		return err
	}
	return c.Stop()
}

func (c *Controller) RecoveryReady() bool {
	return c.state.Snapshot().Protection == model.ProtectionStop
}
