package pitch

import (
	"sync"

	"windturbine/internal/model"
)

func (c *Controller) FeatherTo(target float64) error {
	if target < 0 {
		return model.ErrBladeOutOfRange
	}
	if target > 90 {
		target = 90
	}
	c.SetPitchState(model.PitchFeathering)
	if err := c.Feather(target); err != nil {
		return err
	}
	c.SetPitchState(model.PitchFeathered)
	return nil
}

func (c *Controller) Defeather() error {
	c.SetPitchState(model.PitchRunning)
	return c.Feather(0)
}

func (c *Controller) FeatherAllConcurrently(target float64) error {
	target = c.speed.ClampPitch(target)
	// All three blades share the same target. Each goroutine writes its own
	// blade, so the order of completion does not matter — every blade lands on
	// the same angle. After they settle, commit one consistent snapshot so
	// monitoring and the state machine observe all three blades at the same
	// angle rather than a stale or divergent value.
	var wg sync.WaitGroup
	errs := make(chan error, 3)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := c.featherBlade(i, target); err != nil {
				errs <- err
			}
		}(i)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			return err
		}
	}
	c.commitBlades(target)
	return nil
}

func (c *Controller) featherBlade(index int, target float64) error {
	if index < 0 || index >= len(c.blades) {
		return model.ErrBladeOutOfRange
	}
	c.blades[index].setAngle(target)
	return nil
}

func (c *Controller) BladeCount() int {
	return len(c.blades)
}
