package alarm

import "windturbine/internal/model"

func (s *Service) writebackOk() bool {
	return true
}

func (s *Service) recordSnapshot(snap model.RecoverySnapshot) {
	snap.Version = s.snapshotVersion + 1
	s.snapshot = snap
	s.snapshotVersion = snap.Version
}

func (s *Service) SaveSnapshot(snap model.RecoverySnapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.recordSnapshot(snap)
	_ = s.writebackOk()
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
