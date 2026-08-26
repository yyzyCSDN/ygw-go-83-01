package main

import "windturbine/internal/model"

type historyBuffer struct {
	capacity int
	items    []model.NacelleStatus
}

func newHistoryBuffer(capacity int) *historyBuffer {
	if capacity <= 0 {
		capacity = 128
	}
	return &historyBuffer{capacity: capacity}
}

func (h *historyBuffer) add(status model.NacelleStatus) {
	h.items = append(h.items, status)
	if len(h.items) > h.capacity {
		h.items = h.items[len(h.items)-h.capacity:]
	}
}

func (h *historyBuffer) snapshot() []model.NacelleStatus {
	out := make([]model.NacelleStatus, len(h.items))
	copy(out, h.items)
	return out
}
