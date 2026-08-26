package alarm_test

import (
	"testing"

	"windturbine/internal/alarm"
	"windturbine/internal/core"
	"windturbine/internal/model"
	"windturbine/internal/safe"
	"windturbine/internal/yaw"
)

func TestRecoverWritebackErrorNotSwallowed(t *testing.T) {
	state := core.NewState()
	alarmSvc := alarm.NewService()
	yawCtl := yaw.NewController(state, alarmSvc)
	safeCtl := safe.NewController(state, yawCtl, alarmSvc)
	state.SetProtection(model.ProtectionStop)
	_ = alarmSvc.SaveSnapshot(model.RecoverySnapshot{Protection: model.ProtectionNormal})
	alarmSvc.Close()
	if err := safeCtl.Recover(); err == nil {
		t.Fatal("expected writeback error")
	}
}
