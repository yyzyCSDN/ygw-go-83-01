package main

import (
	"windturbine/internal/alarm"
	"windturbine/internal/core"
	"windturbine/internal/model"
	"windturbine/internal/pitch"
	"windturbine/internal/record"
	"windturbine/internal/safe"
	"windturbine/internal/sensor"
	"windturbine/internal/speed"
	"windturbine/internal/yaw"
)

type Controller struct {
	state     *core.State
	sensors   *sensor.Collector
	speed     *speed.Protector
	pitch     *pitch.Controller
	yaw       *yaw.Controller
	safe      *safe.Controller
	alarm     *alarm.Service
	journal   *record.Service
	sensorID  string
	powerCurve *speed.PowerCurve
	power      *speed.PowerController
	turbine    *speed.TurbineControl
	history    *historyBuffer
	series     *model.MetricSeries
	sequence   int64
}

func NewController(journalDir string) *Controller {
	state := core.NewState()
	collect := sensor.NewCollector()
	collect.Register("anemometer")
	table := speed.DefaultLimitTable()
	protector := speed.NewProtector(state, table, collect, "anemometer")
	alarmSvc := alarm.NewService()
	pitchCtl := pitch.NewController(state, protector, alarmSvc)
	yawCtl := yaw.NewController(state, alarmSvc)
	safeCtl := safe.NewController(state, yawCtl, alarmSvc)
	journal := record.NewService(journalDir, 64*1024)
	powerCtrl := speed.NewPowerController(speed.DefaultPowerCurve(), speed.DefaultPowerCurve().RatedWindSpeed())
	return &Controller{
		state:    state,
		sensors:  collect,
		speed:    protector,
		pitch:    pitchCtl,
		yaw:      yawCtl,
		safe:     safeCtl,
		alarm:    alarmSvc,
		journal:  journal,
		sensorID: "anemometer",
		powerCurve: speed.DefaultPowerCurve(),
		power:      powerCtrl,
		turbine:    speed.NewTurbineControl(protector, powerCtrl),
		history:    newHistoryBuffer(128),
		series:     model.NewMetricSeries("wind"),
	}
}

func (c *Controller) Status() model.NacelleStatus {
	status := c.state.Snapshot()
	return model.NacelleStatus{
		PitchState:     status.PitchState,
		PitchAngle:     status.PitchAngle,
		YawState:       status.YawState,
		YawAngle:       status.YawAngle,
		Protection:     status.Protection,
		RotorSpeed:     status.RotorSpeed,
		CommandedPitch: status.CommandedPitch,
		BladeAngles:    status.BladeAngles,
	}
}

func (c *Controller) FeedSample(windSpeed, rotorSpeed float64) {
	c.sensors.Feed(c.sensorID, windSpeed, rotorSpeed)
}

func (c *Controller) Evaluate() model.ProtectionState {
	protection := c.speed.Evaluate()
	if protection == model.ProtectionStop {
		_ = c.safe.EmergencyShutdown()
	}
	return protection
}

func (c *Controller) Tick(windDirection float64) model.NacelleStatus {
	sample := c.sensors.LastSample(c.sensorID)
	if sample == nil || !sample.HasData {
		return c.Status()
	}
	protection := c.Evaluate()
	c.turbine.Step(c.sensors.AverageWind(), sample.RotorSpeed)
	if protection == model.ProtectionStop {
		if !c.safe.IsStopped() {
			_ = c.safe.EmergencyShutdown()
		}
	} else if c.safe.IsStopped() {
		_ = c.safe.RecoverSequence()
	}
	if !c.safe.IsStopping() {
		_ = c.pitch.ApplyPitch()
	}
	target := c.yaw.TrackWind(windDirection)
	if !c.safe.IsStopping() {
		_ = c.yaw.Align(target)
	}
	c.sequence++
	_ = c.Record("control", "tick")
	c.series.Append(sample.WindSpeed, sample.Timestamp)
	status := c.Status()
	c.history.add(status)
	return status
}

func (c *Controller) History() []model.NacelleStatus {
	return c.history.snapshot()
}

func (c *Controller) ApplyPitch() error {
	return c.pitch.ApplyPitch()
}

func (c *Controller) YawTo(target float64) error {
	return c.yaw.YawTo(target)
}

func (c *Controller) Stop() error {
	return c.safe.StopSequence()
}

func (c *Controller) Recover() error {
	return c.safe.RecoverSequence()
}

func (c *Controller) Record(kind, message string) error {
	return c.journal.Append(record.Event{ID: kind, Kind: kind, Message: message})
}

func (c *Controller) Sequence() int64 {
	return c.sequence
}

func (c *Controller) PowerController() *speed.PowerController {
	return c.power
}

func (c *Controller) ResetMetrics() {
	c.series.Reset()
}
