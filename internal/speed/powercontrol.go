package speed

import "windturbine/internal/model"

type PowerController struct {
	curve    *PowerCurve
	limit    float64
	derated  bool
}

func NewPowerController(curve *PowerCurve, limit float64) *PowerController {
	return &PowerController{curve: curve, limit: limit}
}

func (p *PowerController) Derate(enabled bool) {
	p.derated = enabled
}

func (p *PowerController) ActivePower(windSpeed float64) float64 {
	power := p.curve.PowerAt(windSpeed)
	if p.derated {
		power *= 0.7
	}
	if p.limit > 0 && power > p.limit {
		return p.limit
	}
	return power
}

func (p *PowerController) IsDerated() bool {
	return p.derated
}

func (p *PowerController) PowerLimit() float64 {
	return p.limit
}

func (p *PowerController) Production(windSpeed float64, seconds float64) float64 {
	return p.ActivePower(windSpeed) * seconds / 3600.0
}

func (p *Protector) DerateDecision(protection model.ProtectionState) bool {
	return protection == model.ProtectionOverspeed
}
