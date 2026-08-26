package record

import "sort"

func (s *Service) Query(kind string) []Event {
	events, err := s.Replay()
	if err != nil {
		return nil
	}
	result := make([]Event, 0, len(events))
	for _, e := range events {
		if kind == "" || e.Kind == kind {
			result = append(result, e)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].At.Before(result[j].At) })
	return result
}

func (s *Service) CountByKind() map[string]int {
	events, err := s.Replay()
	if err != nil {
		return map[string]int{}
	}
	counts := make(map[string]int)
	for _, e := range events {
		counts[e.Kind]++
	}
	return counts
}
