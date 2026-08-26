package speed

import "windturbine/internal/model"

type TurbineControl struct {
	protector *Protector
	power     *PowerController
}

func NewTurbineControl(protector *Protector, power *PowerController) *TurbineControl {
	return &TurbineControl{protector: protector, power: power}
}

func (t *TurbineControl) Step(windSpeed, rotorSpeed float64) model.ProtectionState {
	protection := t.protector.Decide(rotorSpeed).Protection
	t.power.Derate(t.protector.DerateDecision(protection))
	return protection
}

func (t *TurbineControl) ActivePower(windSpeed float64) float64 {
	return t.power.ActivePower(windSpeed)
}

func (t *TurbineControl) IsDerated() bool {
	return t.power.IsDerated()
}

func (t *TurbineControl) PowerLimit() float64 {
	return t.power.PowerLimit()
}
