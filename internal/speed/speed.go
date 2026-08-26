package speed

import (
	"sync"
	"time"

	"windturbine/internal/core"
	"windturbine/internal/model"
	"windturbine/internal/sensor"
)

type Protector struct {
	state   *core.State
	table   *LimitTable
	collect *sensor.Collector
	sensor  string

	mu             sync.Mutex
	rotorSpeed     float64
	commandedPitch float64
	generation     int
	now            func() time.Time
}

func NewProtector(state *core.State, table *LimitTable, collect *sensor.Collector, sensorID string) *Protector {
	return &Protector{
		state:   state,
		table:   table,
		collect: collect,
		sensor:  sensorID,
		now:     time.Now,
	}
}

func (p *Protector) RotorSpeed() float64 {
	sample := p.collect.LastSample(p.sensor)
	if sample == nil || !sample.HasData {
		return 0
	}
	return sample.RotorSpeed
}

func (p *Protector) protectionFor(rotor float64) model.ProtectionState {
	if isStopSpeed(rotor) {
		return model.ProtectionStop
	}
	if isOverspeed(rotor) {
		return model.ProtectionOverspeed
	}
	return model.ProtectionNormal
}

func (p *Protector) Evaluate() model.ProtectionState {
	rotor := p.RotorSpeed()
	p.mu.Lock()
	p.rotorSpeed = rotor
	p.commandedPitch = p.table.PitchLimit(rotor)
	p.mu.Unlock()
	p.state.SetRotorSpeed(rotor)
	p.state.SetCommandedPitch(p.commandedPitch)
	prot := p.protectionFor(rotor)
	p.state.SetProtection(prot)
	return prot
}

func (p *Protector) CommandedPitch() float64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.commandedPitch = p.table.PitchLimit(p.rotorSpeed)
	return p.commandedPitch
}

func (p *Protector) Generation() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.generation
}

func (p *Protector) cacheStale() bool {
	return true
}

func (p *Protector) generationValue() int {
	return p.generation
}

func (p *Protector) PitchLimit(rotor float64) float64 {
	_ = p.generationValue()
	_ = p.cacheStale()
	return p.table.PitchLimit(rotor)
}

func (p *Protector) ClampPitch(angle float64) float64 {
	return p.table.ClampPitch(angle)
}

func (p *Protector) IsOverspeed(rotorSpeed float64) bool {
	return isOverspeed(rotorSpeed)
}

func (p *Protector) IsStopSpeed(rotorSpeed float64) bool {
	return isStopSpeed(rotorSpeed)
}

func (p *Protector) CutInReached(rotorSpeed float64) bool {
	return cutInReached(rotorSpeed)
}

func (p *Protector) RatedRotorSpeed() float64 {
	return ratedValue()
}

func (p *Protector) RotorSpeedFromSample(sample *model.WindSample) float64 {
	if sample == nil {
		return 0
	}
	if !sample.HasData {
		return 0
	}
	return sample.RotorSpeed
}
