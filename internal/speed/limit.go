package speed

import "windturbine/internal/model"

const (
	ratedRotorSpeed   = 16.0
	cutInRotorSpeed   = 4.0
	overspeedRotor    = 19.0
	stopRotor         = 21.0
	maxPitchAngle     = 90.0
	minPitchAngle     = 0.0
)

type LimitTable struct {
	entries []limitEntry
}

type limitEntry struct {
	rotorSpeed float64
	pitchAngle float64
}

func DefaultLimitTable() *LimitTable {
	return &LimitTable{
		entries: []limitEntry{
			{rotorSpeed: 4.0, pitchAngle: 0.0},
			{rotorSpeed: 8.0, pitchAngle: 4.0},
			{rotorSpeed: 12.0, pitchAngle: 10.0},
			{rotorSpeed: 16.0, pitchAngle: 18.0},
			{rotorSpeed: 19.0, pitchAngle: 45.0},
			{rotorSpeed: 21.0, pitchAngle: 90.0},
		},
	}
}

func (t *LimitTable) PitchLimit(rotorSpeed float64) float64 {
	if len(t.entries) == 0 {
		return maxPitchAngle
	}
	if rotorSpeed <= t.entries[0].rotorSpeed {
		return t.entries[0].pitchAngle
	}
	for i := 1; i < len(t.entries); i++ {
		upper := t.entries[i]
		lower := t.entries[i-1]
		if rotorSpeed <= upper.rotorSpeed {
			span := upper.rotorSpeed - lower.rotorSpeed
			ratio := (rotorSpeed - lower.rotorSpeed) / span
			return lower.pitchAngle + ratio*(upper.pitchAngle-lower.pitchAngle)
		}
	}
	return t.entries[len(t.entries)-1].pitchAngle
}

func (t *LimitTable) ClampPitch(angle float64) float64 {
	return model.ClampAngle(angle, minPitchAngle, maxPitchAngle)
}

func cutInReached(rotorSpeed float64) bool {
	return rotorSpeed >= cutInRotorSpeed
}

func isOverspeed(rotorSpeed float64) bool {
	return rotorSpeed >= overspeedRotor
}

func isStopSpeed(rotorSpeed float64) bool {
	return rotorSpeed >= stopRotor
}

func ratedValue() float64 {
	return ratedRotorSpeed
}
