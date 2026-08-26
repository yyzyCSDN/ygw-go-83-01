package alarm

import "windturbine/internal/model"

type Summary struct {
	Active   int            `json:"active"`
	Critical int            `json:"critical"`
	Warning  int            `json:"warning"`
	Info     int            `json:"info"`
	Items    []model.Alarm  `json:"items"`
}

func (s *Service) Summary() Summary {
	items := s.Alarms()
	summary := Summary{Items: items}
	for _, a := range items {
		if !a.Active {
			continue
		}
		summary.Active++
		switch a.Level {
		case LevelCritical:
			summary.Critical++
		case LevelWarning:
			summary.Warning++
		case LevelInfo:
			summary.Info++
		}
	}
	return summary
}

func (s *Service) Acknowledge(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.alarms[id]
	if !ok {
		return
	}
	a.Active = false
	a.ClearedAt = s.now().UTC()
	s.alarms[id] = a
}
