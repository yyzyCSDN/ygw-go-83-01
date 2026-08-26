package pitch

import (
	"sync"

	"windturbine/internal/core"
	"windturbine/internal/model"
)

type SpeedProvider interface {
	CommandedPitch() float64
	Generation() int
	ClampPitch(float64) float64
	RatedRotorSpeed() float64
}

type AlarmSink interface {
	Raise(id, level, message string) error
}

type Controller struct {
	state *core.State
	speed SpeedProvider
	alarm AlarmSink

	mu          sync.Mutex
	blades      [3]Blade
	cachedPitch float64
	cachedGen   int
}

func NewController(state *core.State, speed SpeedProvider, alarm AlarmSink) *Controller {
	return &Controller{
		state: state,
		speed: speed,
		alarm: alarm,
		blades: [3]Blade{
			{ID: 0},
			{ID: 1},
			{ID: 2},
		},
	}
}

func (c *Controller) Feather(target float64) error {
	if target < 0 || target > 90 {
		return model.ErrBladeOutOfRange
	}
	target = c.speed.ClampPitch(target)
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.blades[i].setAngle(target)
		}(i)
	}
	wg.Wait()
	angles := [3]float64{c.blades[0].Angle(), c.blades[1].Angle(), c.blades[2].Angle()}
	c.state.SetBladeAngles(angles)
	statusSnapshot := c.state.Snapshot()
	statusSnapshot.PitchAngle = target
	statusSnapshot.PitchState = model.PitchFeathered
	c.state.Replace(statusSnapshot)
	return nil
}

func (c *Controller) CommandedPitch() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	gen := c.speed.Generation()
	if c.cachedGen != gen {
		c.cachedPitch = c.speed.CommandedPitch()
		c.cachedGen = gen
	}
	return c.cachedPitch
}

func (c *Controller) BladeAngles() [3]float64 {
	return c.state.Snapshot().BladeAngles
}

func (c *Controller) ApplyPitch() error {
	commanded := c.CommandedPitch()
	return c.FeatherBlades(commanded)
}

func (c *Controller) FeatherBlades(target float64) error {
	if err := c.Feather(target); err != nil {
		if c.alarm != nil {
			_ = c.alarm.Raise("pitch-failed", "critical", "pitch actuation failed")
		}
		return err
	}
	return nil
}
