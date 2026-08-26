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
	return p.rotorFromSample(p.collect.LastSample(p.sensor))
}

// rotorFromSample 从采样读出转速。空读数（未注册/无数据）按无效采样处理，
// 返回 0 转速而非解引用 nil 字段。
func (p *Protector) rotorFromSample(sample *model.WindSample) float64 {
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

// Evaluate 执行一轮转速保护评估。空读数（测风通道未注册或尚无采样）视作
// 无效采样：跳过本轮越限判断，保留上一拍保护状态，协程继续运行，保护链路
// 不能因为一次空读数而崩溃。
func (p *Protector) Evaluate() model.ProtectionState {
	sample := p.collect.LastSample(p.sensor)
	if sample == nil || !sample.HasData {
		return p.state.Snapshot().Protection
	}
	rotor := sample.RotorSpeed
	p.mu.Lock()
	p.rotorSpeed = rotor
	p.commandedPitch = p.table.PitchLimit(rotor)
	p.generation++
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

func (p *Protector) PitchLimit(rotor float64) float64 {
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
