package alarm

import "windturbine/internal/model"

const (
	LevelInfo     = "info"
	LevelWarning  = "warning"
	LevelCritical = "critical"
)

func (s *Service) ActiveByLevel(level string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	count := 0
	for _, a := range s.alarms {
		if a.Active && a.Level == level {
			count++
		}
	}
	return count
}

func (s *Service) HasCritical() bool {
	return s.ActiveByLevel(LevelCritical) > 0
}

func (s *Service) Snapshot() []model.Alarm {
	return s.Alarms()
}

func (s *Service) ClearAll(level string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	cleared := 0
	for id, a := range s.alarms {
		if a.Active && a.Level == level {
			a.Active = false
			a.ClearedAt = s.now().UTC()
			s.alarms[id] = a
			cleared++
		}
	}
	return cleared
}
