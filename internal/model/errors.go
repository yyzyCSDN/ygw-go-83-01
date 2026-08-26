package model

import "errors"

var (
	ErrNoData           = errors.New("sensor returned no data")
	ErrSensorUnhealthy  = errors.New("sensor is unhealthy")
	ErrPitchFailed      = errors.New("pitch actuation failed")
	ErrYawTimeout       = errors.New("yaw alignment timed out")
	ErrYawAborted       = errors.New("yaw alignment aborted")
	ErrStopInProgress   = errors.New("safety stop already in progress")
	ErrNotStopping      = errors.New("turbine is not in a stopping state")
	ErrRecoverFailed    = errors.New("recovery writeback failed")
	ErrJournalClosed    = errors.New("event journal is closed")
	ErrBladeOutOfRange  = errors.New("blade angle is out of range")
	ErrSnapshotStale    = errors.New("recovery snapshot is stale")
	ErrCommandMismatch  = errors.New("commanded pitch mismatch")
	ErrProtectionActive = errors.New("protection is active")
)
