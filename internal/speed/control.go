package speed

import (
	"windturbine/internal/model"
)

type ControlDecision struct {
	Protection model.ProtectionState
	Pitch      float64
	ShouldStop bool
}

func (p *Protector) Decide(rotorSpeed float64) ControlDecision {
	protection := p.protectionFor(rotorSpeed)
	pitch := p.table.PitchLimit(rotorSpeed)
	return ControlDecision{
		Protection: protection,
		Pitch:      pitch,
		ShouldStop: protection == model.ProtectionStop,
	}
}

func (p *Protector) EmergencyFeatherAngle(rotorSpeed float64) float64 {
	if isStopSpeed(rotorSpeed) {
		return maxPitchAngle
	}
	if isOverspeed(rotorSpeed) {
		return 45
	}
	return 0
}

func (p *Protector) RotorFromWind(windSpeed float64) float64 {
	if windSpeed < cutInRotorSpeed {
		return 0
	}
	return windSpeed * 1.15
}

func (p *Protector) RecoveryRotorSpeed() float64 {
	return ratedRotorSpeed - 0.5
}
