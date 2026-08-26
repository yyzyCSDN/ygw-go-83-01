package record_test

import (
	"testing"

	"windturbine/internal/record"
)

func TestEventHandleClosed(t *testing.T) {
	dir := t.TempDir()
	svc := record.NewService(dir, 1)
	defer func() { _ = svc.Close() }()
	for i := 0; i < 200; i++ {
		if err := svc.Append(record.Event{ID: "e", Kind: "control", Message: "m"}); err != nil {
			t.Fatal(err)
		}
	}
	if svc.OpenHandleCount() > 2 {
		t.Fatalf("handle leak: %d", svc.OpenHandleCount())
	}
}
