package model

const (
	MaxWindSpeed     = 60.0
	MaxRotorSpeed    = 25.0
	MaxPitchAngle    = 90.0
	MinPitchAngle    = 0.0
	MaxYawAngle      = 359.9
	RatedPowerKW     = 3300.0
)

func ValidateWindSpeed(value float64) bool {
	return value >= 0 && value <= MaxWindSpeed
}

func ValidateRotorSpeed(value float64) bool {
	return value >= 0 && value <= MaxRotorSpeed
}

func ValidatePitchAngle(value float64) bool {
	return value >= MinPitchAngle && value <= MaxPitchAngle
}

func ValidateYawAngle(value float64) bool {
	return value >= 0 && value <= MaxYawAngle
}
