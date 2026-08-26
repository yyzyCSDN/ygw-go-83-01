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
