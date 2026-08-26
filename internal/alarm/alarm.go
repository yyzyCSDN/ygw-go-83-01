package alarm

import (
	"sort"
	"sync"
	"time"

	"windturbine/internal/model"
)

type Service struct {
	mu              sync.Mutex
	alarms          map[string]model.Alarm
	snapshot        model.RecoverySnapshot
	snapshotVersion int
	closed          bool
	now             func() time.Time
}

func NewService() *Service {
	return &Service{
		alarms: make(map[string]model.Alarm),
		now:    time.Now,
	}
}

func (s *Service) Raise(id, level, message string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.alarms[id]
	if !ok {
		existing = model.Alarm{ID: id}
	}
	existing.Level = level
	existing.Message = message
	existing.Active = true
	existing.RaisedAt = s.now().UTC()
	s.alarms[id] = existing
	return nil
}

func (s *Service) Clear(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.alarms[id]
	if !ok {
		return nil
	}
	existing.Active = false
	existing.ClearedAt = s.now().UTC()
	s.alarms[id] = existing
	return nil
}

func (s *Service) ActiveCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	count := 0
	for _, a := range s.alarms {
		if a.Active {
			count++
		}
	}
	return count
}

func (s *Service) Alarms() []model.Alarm {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]model.Alarm, 0, len(s.alarms))
	for _, a := range s.alarms {
		result = append(result, a)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}

func (s *Service) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
}
