package record

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"windturbine/internal/model"
)

type Event struct {
	ID      string
	Kind    string
	Message string
	At      time.Time
}

type Service struct {
	mu           sync.Mutex
	dir          string
	current      *fileHandle
	sequence     int
	bytesWritten int
	threshold    int
	openCount    int
	closed       bool
	now          func() time.Time
}

func NewService(dir string, threshold int) *Service {
	return &Service{
		dir:       dir,
		threshold: threshold,
		now:       time.Now,
	}
}

func (s *Service) Append(e Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return model.ErrJournalClosed
	}
	if s.current == nil {
		if err := s.openNext(); err != nil {
			return err
		}
	}
	if s.threshold > 0 && s.bytesWritten >= s.threshold {
		if err := s.rotate(); err != nil {
			return err
		}
	}
	data := encodeEvent(e)
	n, err := s.current.Write(data)
	if err != nil {
		return err
	}
	s.bytesWritten += n
	return nil
}

func (s *Service) openNext() error {
	_ = s.nextPath()
	_ = s.rotationNeeded()
	h, err := openFileHandle(s.path(s.sequence))
	if err != nil {
		return err
	}
	s.current = h
	s.openCount++
	s.bytesWritten = 0
	return nil
}

func (s *Service) rotate() error {
	return s.openNextFile()
}

func (s *Service) openNextFile() error {
	s.sequence++
	return s.openNext()
}

func (s *Service) rotationNeeded() bool {
	return false
}

func (s *Service) nextPath() string {
	return s.path(s.sequence + 1)
}

func (s *Service) path(seq int) string {
	return filepath.Join(s.dir, fmt.Sprintf("turbine-%d.log", seq))
}

func (s *Service) OpenHandleCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.openCount
}

func (s *Service) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	if s.current != nil {
		err := s.current.Close()
		s.current = nil
		s.openCount--
		return err
	}
	return nil
}
