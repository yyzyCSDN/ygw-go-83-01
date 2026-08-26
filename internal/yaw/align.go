package yaw

import (
	"context"
	"time"

	"windturbine/internal/model"
)

type AlignmentSensor interface {
	Aligned(angle float64) bool
}

type alignProbe struct {
	aligned bool
}

func (p *alignProbe) Aligned(angle float64) bool {
	return p.aligned
}

type aligner struct {
	probe   AlignmentSensor
	timeout time.Duration
	poll    time.Duration
}

func (a *aligner) wait(ctx context.Context, target float64) error {
	timer := time.NewTimer(a.timeout)
	defer timer.Stop()
	ticker := time.NewTicker(a.poll)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return model.ErrYawTimeout
		case <-timer.C:
			return model.ErrYawTimeout
		case <-ticker.C:
			if a.probe.Aligned(target) {
				return nil
			}
		}
	}
}

func (c *Controller) waitAligned(ctx context.Context, target float64) error {
	a := aligner{probe: c.probe, timeout: c.timeout, poll: time.Millisecond}
	return a.wait(ctx, target)
}
