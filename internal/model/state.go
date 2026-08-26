package model

func (s PitchState) Valid() bool {
	switch s {
	case PitchRunning, PitchFeathering, PitchFeathered:
		return true
	default:
		return false
	}
}

func (s YawState) Valid() bool {
	switch s {
	case YawIdle, YawYawing, YawAligned:
		return true
	default:
		return false
	}
}

func (s ProtectionState) Valid() bool {
	switch s {
	case ProtectionNormal, ProtectionOverspeed, ProtectionStop:
		return true
	default:
		return false
	}
}

func (s ProtectionState) Severity() int {
	switch s {
	case ProtectionNormal:
		return 0
	case ProtectionOverspeed:
		return 1
	case ProtectionStop:
		return 2
	default:
		return -1
	}
}

func ClampAngle(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
