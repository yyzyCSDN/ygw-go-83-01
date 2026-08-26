package core

import "windturbine/internal/model"

func (s *State) SnapshotForRecovery() model.RecoverySnapshot {
	status := s.Snapshot()
	return model.RecoverySnapshot{
		Protection: status.Protection,
		PitchState: status.PitchState,
		YawState:   status.YawState,
		PitchAngle: status.PitchAngle,
		YawAngle:   status.YawAngle,
	}
}

func (s *State) RestoreProtection(p model.ProtectionState) {
	s.SetProtection(p)
}

func (s *State) ProtectionState() model.ProtectionState {
	return s.Snapshot().Protection
}
