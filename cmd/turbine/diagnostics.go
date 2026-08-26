package main

import (
	"time"

	"windturbine/internal/alarm"
	"windturbine/internal/model"
	"windturbine/internal/sensor"
	"windturbine/internal/speed"
	"windturbine/internal/yaw"
)

type diagnostics struct {
	Sensors     sensorDiagnostics     `json:"sensors"`
	Alarms      alarm.Summary         `json:"alarms"`
	AlarmCounts map[string]int        `json:"alarm_counts"`
	Cleared     int                   `json:"cleared_alarms"`
	Critical    bool                  `json:"critical"`
	Journal     journalDiagnostics    `json:"journal"`
	Speed       speedDiagnostics      `json:"speed"`
	Pitch       pitchDiagnostics      `json:"pitch"`
	Yaw         yawDiagnostics        `json:"yaw"`
	Safe        safeDiagnostics       `json:"safe"`
	Model       modelDiagnostics      `json:"model"`
}

type sensorDiagnostics struct {
	Registry        []sensor.SensorStatus `json:"registry"`
	HealthyCount    int                   `json:"healthy_count"`
	Count           int                   `json:"count"`
	Plausible       bool                  `json:"plausible"`
	SmoothWind      float64               `json:"smooth_wind"`
	CalibratedWind  float64               `json:"calibrated_wind"`
}

type journalDiagnostics struct {
	CountByKind map[string]int `json:"count_by_kind"`
	TotalBytes  int            `json:"total_bytes"`
	LastEvent   string         `json:"last_event"`
	Files       int            `json:"files"`
	Recent      int            `json:"recent_events"`
}

type speedDiagnostics struct {
	Telemetry            speed.SpeedTelemetry `json:"telemetry"`
	CachedRotor          float64              `json:"cached_rotor"`
	CachedPitch          float64              `json:"cached_pitch"`
	Protection           string               `json:"protection"`
	EmergencyFeather     float64              `json:"emergency_feather"`
	RotorFromWind        float64              `json:"rotor_from_wind"`
	RecoveryRotor        float64              `json:"recovery_rotor"`
	Overspeed            bool                 `json:"overspeed"`
	StopSpeed            bool                 `json:"stop_speed"`
	CutIn                bool                 `json:"cut_in"`
	Hysteresis           string               `json:"hysteresis"`
	Power                float64              `json:"power"`
	PowerLimit           float64              `json:"power_limit"`
	Derated              bool                 `json:"derated"`
	CutInWindSpeed       float64              `json:"cut_in_wind_speed"`
	RotorFromSample      float64              `json:"rotor_from_sample"`
}

type pitchDiagnostics struct {
	Imbalance      float64 `json:"imbalance"`
	Consistent     bool    `json:"consistent"`
	AverageAngle   float64 `json:"average_angle"`
	BladeCount     int     `json:"blade_count"`
	MaxAngle       float64 `json:"max_angle"`
	TransitionTo   string  `json:"transition_to"`
	CurrentState   string  `json:"current_state"`
	BladesFeathered bool   `json:"blades_feathered"`
}

type yawDiagnostics struct {
	State       string  `json:"state"`
	Angle       float64 `json:"angle"`
	Target      float64 `json:"target"`
	Deadband    float64 `json:"deadband"`
	NeedsYaw    bool    `json:"needs_yaw"`
	CanYaw      bool    `json:"can_yaw"`
	Idle        bool    `json:"idle"`
	Dominant    float64 `json:"dominant_direction"`
	Total       float64 `json:"rose_total"`
}

type safeDiagnostics struct {
	Protection    string `json:"protection"`
	ShouldStop    bool   `json:"should_stop"`
	StopReason    string `json:"stop_reason"`
	YawEngaged    bool   `json:"yaw_engaged"`
	RecoveryReady bool   `json:"recovery_ready"`
}

type modelDiagnostics struct {
	WindValid    bool    `json:"wind_valid"`
	RotorValid   bool    `json:"rotor_valid"`
	PitchValid   bool    `json:"pitch_valid"`
	YawValid     bool    `json:"yaw_valid"`
	MeanWind     float64 `json:"mean_wind"`
	FrameHealthy bool    `json:"frame_healthy"`
	Production   float64 `json:"production_factor"`
	LastPower    float64 `json:"last_power"`
	Severity     int     `json:"severity"`
	Reading      float64 `json:"reading_wind"`
}

func (c *Controller) Diagnostics() diagnostics {
	status := c.Status()
	rotor := status.RotorSpeed
	wind := c.sensors.AverageWind()
	calibrated := sensor.NewCalibratedCollector(c.sensors, 0, 1)
	windRose := yaw.NewWindRose(8)
	windRose.Add(0, 1)

	speedTelemetry := c.speed.Telemetry()
	powerCtrl := speed.NewPowerController(c.powerCurve, c.powerCurve.RatedWindSpeed())
	powerCtrl.Derate(c.speed.DerateDecision(status.Protection))
	turbine := speed.NewTurbineControl(c.speed, powerCtrl)

	series := model.NewMetricSeries("wind")
	series.Append(wind, time.Now().UTC())
	mean, _ := series.Last()
	_ = mean
	frame := model.NewTelemetryFrame(1, c.TurbineMetrics(), status.Protection, status.PitchState, status.YawState)
	powerReading := model.PowerReading{ActivePower: c.TurbineMetrics().ActivePower}
	sample := c.sensors.LastSample(c.sensorID)
	var readingWind float64
	if sample != nil {
		readingWind = sample.WindSpeed
	}
	reading := model.SensorReading{ID: c.sensorID, WindSpeed: readingWind}
	readingTotal := windRose.Total()
	_ = readingTotal

	recent := c.journal.Query("")
	cleared := c.alarm.ClearedAlarms()

	return diagnostics{
		Sensors: sensorDiagnostics{
			Registry:       c.sensors.Registry(),
			HealthyCount:   c.sensors.HealthySensorCount(),
			Count:          c.sensors.SensorCount(),
			Plausible:      c.sensors.Plausible(c.sensorID),
			SmoothWind:     c.sensors.SmoothWind(c.sensorID, 8),
			CalibratedWind: calibrated.CalibratedWind(c.sensorID),
		},
		Alarms:      c.alarm.Summary(),
		AlarmCounts: map[string]int{"all": c.alarm.AlarmCount(""), "critical": c.alarm.ActiveByLevel(alarm.LevelCritical)},
		Cleared:     len(cleared),
		Critical:    c.alarm.HasCritical(),
		Journal: journalDiagnostics{
			CountByKind: c.journal.CountByKind(),
			TotalBytes:  c.journal.TotalBytes(),
			LastEvent:   c.journal.LastEventAt().String(),
			Files:       len(c.journal.ListFiles()),
			Recent:      len(recent),
		},
		Speed: speedDiagnostics{
			Telemetry:        speedTelemetry,
			CachedRotor:      c.speed.CachedRotorSpeed(),
			CachedPitch:      c.speed.CachedCommandedPitch(),
			Protection:       c.speed.Protection(),
			EmergencyFeather: c.speed.EmergencyFeatherAngle(rotor),
			RotorFromWind:    c.speed.RotorFromWind(wind),
			RecoveryRotor:    c.speed.RecoveryRotorSpeed(),
			Overspeed:        c.speed.IsOverspeed(rotor),
			StopSpeed:        c.speed.IsStopSpeed(rotor),
			CutIn:            c.speed.CutInReached(rotor),
			Hysteresis:       string(c.speed.ProtectionWithHysteresis(status.Protection, rotor)),
			Power:            turbine.ActivePower(wind),
			PowerLimit:       turbine.PowerLimit(),
			Derated:          turbine.IsDerated(),
			CutInWindSpeed:   c.powerCurve.CutInWindSpeed(),
			RotorFromSample:  c.speed.RotorSpeedFromSample(sample),
		},
		Pitch: pitchDiagnostics{
			Imbalance:       c.pitch.BladeImbalance(),
			Consistent:      c.pitch.BladesConsistent(),
			AverageAngle:    c.pitch.AverageBladeAngle(),
			BladeCount:      c.pitch.BladeCount(),
			MaxAngle:        c.pitch.MaxPitchAngle(),
			TransitionTo:    string(c.pitch.Transition(rotor)),
			CurrentState:    string(c.pitch.CurrentPitchState()),
			BladesFeathered: c.pitch.BladesFeathered(),
		},
		Yaw: yawDiagnostics{
			State:    string(c.yaw.CurrentYawState()),
			Angle:    c.yaw.CurrentYawAngle(),
			Target:   c.yaw.YawTarget(0),
			Deadband: c.yaw.Deadband(),
			NeedsYaw: c.yaw.NeedsYaw(0),
			CanYaw:   c.yaw.CanYaw(),
			Idle:     c.yaw.YawIdle(),
			Dominant: windRose.DominantDirection(),
			Total:    windRose.Total(),
		},
		Safe: safeDiagnostics{
			Protection:    string(c.safe.ProtectionState()),
			ShouldStop:    c.safe.ShouldStop(status.Protection),
			StopReason:    string(c.safe.StopReason()),
			YawEngaged:    c.safe.YawEngaged(),
			RecoveryReady: c.safe.RecoveryReady(),
		},
		Model: modelDiagnostics{
			WindValid:    model.ValidateWindSpeed(wind),
			RotorValid:   model.ValidateRotorSpeed(rotor),
			PitchValid:   model.ValidatePitchAngle(status.PitchAngle),
			YawValid:     model.ValidateYawAngle(status.YawAngle),
			MeanWind:     series.Mean(),
			FrameHealthy: frame.IsHealthy(),
			Production:   c.TurbineMetrics().ProductionFactor(),
			LastPower:    powerReading.ActivePower,
			Severity:     int(status.Protection.Severity()),
			Reading:      reading.WindSpeed,
		},
	}
}
