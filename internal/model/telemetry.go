package model

import (
	"encoding/json"
	"time"
)

type TelemetryFrame struct {
	Sequence    int64           `json:"sequence"`
	Timestamp   time.Time       `json:"timestamp"`
	Metrics     TurbineMetrics  `json:"metrics"`
	Protection  ProtectionState `json:"protection"`
	PitchState  PitchState      `json:"pitch_state"`
	YawState    YawState        `json:"yaw_state"`
}

func NewTelemetryFrame(sequence int64, metrics TurbineMetrics, protection ProtectionState, pitchState PitchState, yawState YawState) TelemetryFrame {
	return TelemetryFrame{
		Sequence:   sequence,
		Timestamp:  time.Now().UTC(),
		Metrics:    metrics,
		Protection: protection,
		PitchState: pitchState,
		YawState:   yawState,
	}
}

func (f TelemetryFrame) Marshal() ([]byte, error) {
	return json.Marshal(f)
}

func UnmarshalTelemetryFrame(data []byte) (TelemetryFrame, error) {
	var frame TelemetryFrame
	err := json.Unmarshal(data, &frame)
	return frame, err
}

func (f TelemetryFrame) IsHealthy() bool {
	return f.Protection == ProtectionNormal && f.Metrics.Available()
}
