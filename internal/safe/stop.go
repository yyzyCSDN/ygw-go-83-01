package safe

import "windturbine/internal/model"

func (c *Controller) Stop() error {
	c.mu.Lock()
	if c.stopping {
		c.mu.Unlock()
		return model.ErrStopInProgress
	}
	c.stopping = true
	c.mu.Unlock()

	if err := c.yaw.Abort(); err != nil {
		return err
	}
	c.state.SetProtection(model.ProtectionStop)
	if c.alarm != nil {
		_ = c.alarm.Raise("safety-stop", "critical", "safety stop engaged")
	}
	return nil
}
