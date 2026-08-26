package record

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func (s *Service) ListFiles() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasPrefix(entry.Name(), "turbine-") && strings.HasSuffix(entry.Name(), ".log") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	return names
}

func (s *Service) Prune(keep int) int {
	names := s.ListFiles()
	if len(names) <= keep {
		return 0
	}
	removed := 0
	for _, name := range names[:len(names)-keep] {
		path := filepath.Join(s.dir, name)
		if os.Remove(path) == nil {
			removed++
		}
	}
	return removed
}
