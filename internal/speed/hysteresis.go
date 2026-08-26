package speed

import "windturbine/internal/model"

const (
	overspeedEnter = 19.0
	overspeedExit  = 17.5
	stopEnter      = 21.0
	stopExit       = 19.5
)

func (p *Protector) ProtectionWithHysteresis(current model.ProtectionState, rotorSpeed float64) model.ProtectionState {
	switch current {
	case model.ProtectionStop:
		if rotorSpeed < stopExit {
			return model.ProtectionOverspeed
		}
		return model.ProtectionStop
	case model.ProtectionOverspeed:
		if rotorSpeed >= stopEnter {
			return model.ProtectionStop
		}
		if rotorSpeed < overspeedExit {
			return model.ProtectionNormal
		}
		return model.ProtectionOverspeed
	default:
		if rotorSpeed >= overspeedEnter {
			return model.ProtectionOverspeed
		}
		return model.ProtectionNormal
	}
}
