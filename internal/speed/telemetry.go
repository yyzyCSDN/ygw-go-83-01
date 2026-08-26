package speed

type SpeedTelemetry struct {
	RotorSpeed     float64
	CommandedPitch float64
	Protection     string
	Generation     int
}

func (p *Protector) Telemetry() SpeedTelemetry {
	p.mu.Lock()
	defer p.mu.Unlock()
	return SpeedTelemetry{
		RotorSpeed:     p.rotorSpeed,
		CommandedPitch: p.commandedPitch,
		Protection:     string(p.state.Snapshot().Protection),
		Generation:     p.generation,
	}
}

func (p *Protector) CachedRotorSpeed() float64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.rotorSpeed
}

func (p *Protector) CachedCommandedPitch() float64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.commandedPitch
}

func (p *Protector) Protection() string {
	return string(p.state.Snapshot().Protection)
}
