package safe_test

import (
	"testing"

	"windturbine/internal/alarm"
	"windturbine/internal/core"
	"windturbine/internal/model"
	"windturbine/internal/safe"
	"windturbine/internal/yaw"
)

func TestRecoveryUsesLatestSnapshot(t *testing.T) {
	state := core.NewState()
	alarmSvc := alarm.NewService()
	yawCtl := yaw.NewController(state, alarmSvc)
	_ = alarmSvc.SaveSnapshot(model.RecoverySnapshot{Protection: model.ProtectionOverspeed})
	safeCtl := safe.NewController(state, yawCtl, alarmSvc)
	state.SetProtection(model.ProtectionStop)
	_ = alarmSvc.SaveSnapshot(model.RecoverySnapshot{Protection: model.ProtectionNormal})
	if err := safeCtl.Recover(); err != nil {
		t.Fatal(err)
	}
	if state.Snapshot().Protection != model.ProtectionNormal {
		t.Fatalf("recovery used stale snapshot: %v", state.Snapshot().Protection)
	}
}
