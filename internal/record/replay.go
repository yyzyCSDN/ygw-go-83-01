package record

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func (s *Service) Replay() ([]Event, error) {
	names := s.ListFiles()
	sort.Strings(names)
	var events []Event
	for _, name := range names {
		path := filepath.Join(s.dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		for _, line := range strings.Split(string(data), "\n") {
			if line == "" {
				continue
			}
			if len(line) < 8 {
				continue
			}
			e, err := decodeEvent([]byte(line))
			if err != nil {
				continue
			}
			events = append(events, e)
		}
	}
	return events, nil
}
