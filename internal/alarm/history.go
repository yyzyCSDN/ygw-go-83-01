package alarm

import "windturbine/internal/model"

func (s *Service) ClearedAlarms() []model.Alarm {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]model.Alarm, 0)
	for _, a := range s.alarms {
		if !a.Active {
			result = append(result, a)
		}
	}
	return result
}

func (s *Service) AlarmCount(level string) int {
	if level == "" {
		return len(s.Alarms())
	}
	count := 0
	for _, a := range s.Alarms() {
		if a.Level == level {
			count++
		}
	}
	return count
}
