package core

import (
	"sync"

	"windturbine/internal/model"
)

type Status struct {
	PitchState     model.PitchState
	PitchAngle     float64
	YawState       model.YawState
	YawAngle       float64
	Protection     model.ProtectionState
	RotorSpeed     float64
	CommandedPitch float64
	BladeAngles    [3]float64
}

type State struct {
	mu     sync.RWMutex
	status Status
}

func NewState() *State {
	return &State{
		status: Status{
			PitchState: model.PitchRunning,
			YawState:   model.YawIdle,
			Protection: model.ProtectionNormal,
		},
	}
}

func (s *State) Snapshot() Status {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.status
}

// Update mutates the shared status under the state lock. Callers receive a
// pointer to the current status and must touch only the fields they own; this
// keeps concurrent writers (pitch, yaw, ...) from overwriting each other's
// changes the way a whole-structure Replace would.
func (s *State) Update(fn func(*Status)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fn(&s.status)
}

func (s *State) SetProtection(p model.ProtectionState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status.Protection = p
}

func (s *State) SetRotorSpeed(v float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status.RotorSpeed = v
}

func (s *State) SetCommandedPitch(v float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status.CommandedPitch = v
}

func (s *State) SetBladeAngles(angles [3]float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status.BladeAngles = angles
}
