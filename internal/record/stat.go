package record

import (
	"os"
	"path/filepath"
	"time"
)

func (s *Service) TotalBytes() int {
	total := 0
	for _, name := range s.ListFiles() {
		info, err := os.Stat(filepath.Join(s.dir, name))
		if err != nil {
			continue
		}
		total += int(info.Size())
	}
	return total
}

func (s *Service) LastEventAt() time.Time {
	events, err := s.Replay()
	if err != nil || len(events) == 0 {
		return time.Time{}
	}
	last := events[0]
	for _, e := range events {
		if e.At.After(last.At) {
			last = e
		}
	}
	return last.At
}
