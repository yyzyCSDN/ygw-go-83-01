package safe_test

import (
	"sync"
	"testing"

	"windturbine/internal/alarm"
	"windturbine/internal/core"
	"windturbine/internal/model"
	"windturbine/internal/safe"
	"windturbine/internal/yaw"
)

func TestSafeStopExcludesYaw(t *testing.T) {
	for i := 0; i < 500; i++ {
		state := core.NewState()
		alarmSvc := alarm.NewService()
		yawCtl := yaw.NewController(state, alarmSvc)
		safeCtl := safe.NewController(state, yawCtl, alarmSvc)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() { defer wg.Done(); _ = safeCtl.Stop() }()
		go func() { defer wg.Done(); _ = yawCtl.YawTo(90) }()
		wg.Wait()
		if state.Snapshot().Protection != model.ProtectionStop {
			t.Fatalf("protection lost at iteration %d: %v", i, state.Snapshot().Protection)
		}
	}
}
