package pitch

import (
	"windturbine/internal/core"
	"windturbine/internal/model"
)

func (c *Controller) Transition(rotorSpeed float64) model.PitchState {
	current := c.state.Snapshot().PitchState
	switch current {
	case model.PitchRunning:
		if rotorSpeed > c.speed.RatedRotorSpeed() {
			return model.PitchFeathering
		}
		return model.PitchRunning
	case model.PitchFeathering:
		if c.BladesFeathered() {
			return model.PitchFeathered
		}
		return model.PitchFeathering
	case model.PitchFeathered:
		if rotorSpeed < c.speed.RatedRotorSpeed() {
			return model.PitchRunning
		}
		return model.PitchFeathered
	default:
		return model.PitchRunning
	}
}

func (c *Controller) BladesFeathered() bool {
	angles := c.BladeAngles()
	commanded := c.state.Snapshot().CommandedPitch
	for _, angle := range angles {
		if angle < commanded-0.01 {
			return false
		}
	}
	return true
}

func (c *Controller) SetPitchState(state model.PitchState) {
	c.state.Update(func(s *core.Status) {
		s.PitchState = state
	})
}

func (c *Controller) CurrentPitchState() model.PitchState {
	return c.state.Snapshot().PitchState
}
