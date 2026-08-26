package safe

import "windturbine/internal/model"

// commitStop sets Protection=stop atomically on the protection field only.
// It deliberately avoids the Snapshot/Replace read-modify-write pattern: that
// would read the whole status (possibly with a stale Protection), mutate one
// field, and write the whole struct back, racing a concurrent yaw that could
// overwrite the stop. SetProtection locks and touches only this field.
func (c *Controller) commitStop() {
	c.state.SetProtection(model.ProtectionStop)
}

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
	c.commitStop()
	if c.alarm != nil {
		_ = c.alarm.Raise("safety-stop", "critical", "safety stop engaged")
	}
	return nil
}
