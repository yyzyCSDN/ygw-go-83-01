package model

import "time"

type PitchState string

const (
	PitchRunning    PitchState = "running"
	PitchFeathering PitchState = "feathering"
	PitchFeathered  PitchState = "feathered"
)

type YawState string

const (
	YawIdle    YawState = "idle"
	YawYawing  YawState = "yawing"
	YawAligned YawState = "aligned"
)

type ProtectionState string

const (
	ProtectionNormal    ProtectionState = "normal"
	ProtectionOverspeed ProtectionState = "overspeed"
	ProtectionStop      ProtectionState = "stop"
)

type WindSample struct {
	WindSpeed  float64
	RotorSpeed float64
	Timestamp  time.Time
	HasData    bool
}

type NacelleStatus struct {
	PitchState     PitchState
	PitchAngle     float64
	YawState       YawState
	YawAngle       float64
	Protection     ProtectionState
	RotorSpeed     float64
	CommandedPitch float64
	BladeAngles    [3]float64
}

type SensorReading struct {
	ID         string
	WindSpeed  float64
	RotorSpeed float64
	Timestamp  time.Time
	Healthy    bool
}

type Alarm struct {
	ID        string
	Level     string
	Message   string
	RaisedAt  time.Time
	ClearedAt time.Time
	Active    bool
}

type RecoverySnapshot struct {
	Protection ProtectionState
	PitchState PitchState
	YawState   YawState
	PitchAngle float64
	YawAngle   float64
	Version    int
}
