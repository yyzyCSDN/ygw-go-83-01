package alarm

import "windturbine/internal/model"

func (s *Service) SaveSnapshot(snap model.RecoverySnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return model.ErrRecoverFailed
	}
	snap.Version = s.snapshotVersion + 1
	s.snapshot = snap
	s.snapshotVersion = snap.Version
	return nil
}

func (s *Service) LatestSnapshot() (model.RecoverySnapshot, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.snapshotVersion == 0 {
		return model.RecoverySnapshot{}, model.ErrSnapshotStale
	}
	return s.snapshot, nil
}
