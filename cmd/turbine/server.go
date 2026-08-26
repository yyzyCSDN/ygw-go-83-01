package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"windturbine/internal/alarm"
	"windturbine/internal/model"
)

type Server struct {
	controller *Controller
	webRoot    string
}

func NewServer(controller *Controller, webRoot string) *Server {
	return &Server{controller: controller, webRoot: webRoot}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/api/status", s.handleStatus)
	mux.HandleFunc("/api/metrics", s.handleMetrics)
	mux.HandleFunc("/api/history", s.handleHistory)
	mux.HandleFunc("/api/alarms", s.handleAlarms)
	mux.HandleFunc("/api/sensors", s.handleSensors)
	mux.HandleFunc("/api/journal", s.handleJournal)
	mux.HandleFunc("/api/fault", s.handleFault)
	mux.HandleFunc("/api/diagnostics", s.handleDiagnostics)
	mux.HandleFunc("/api/ack", s.handleAck)
	mux.HandleFunc("/api/normal-shutdown", s.handleNormalShutdown)
	mux.HandleFunc("/api/feather-concurrent", s.handleFeatherConcurrent)
	mux.HandleFunc("/api/feather-limited", s.handleFeatherLimited)
	mux.HandleFunc("/api/mark-unhealthy", s.handleMarkUnhealthy)
	mux.HandleFunc("/api/clear-warnings", s.handleClearWarnings)
	mux.HandleFunc("/api/telemetry", s.handleTelemetry)
	mux.HandleFunc("/api/reset-metrics", s.handleResetMetrics)
	mux.HandleFunc("/api/sample", s.handleSample)
	mux.HandleFunc("/api/evaluate", s.handleEvaluate)
	mux.HandleFunc("/api/pitch", s.handlePitch)
	mux.HandleFunc("/api/feather", s.handleFeather)
	mux.HandleFunc("/api/defeather", s.handleDefeather)
	mux.HandleFunc("/api/yaw", s.handleYaw)
	mux.HandleFunc("/api/stop", s.handleStop)
	mux.HandleFunc("/api/recover", s.handleRecover)
	mux.HandleFunc("/api/tick", s.handleTick)
	mux.HandleFunc("/", s.handleIndex)
	return mux
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.controller.Status())
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.controller.Metrics())
}

func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.controller.History())
}

func (s *Server) handleAlarms(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.controller.alarm.Summary())
}

func (s *Server) handleSensors(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"count":     s.controller.sensors.SensorCount(),
		"healthy":   s.controller.sensors.HealthyIDs(),
		"avg_wind":  s.controller.sensors.AverageWind(),
		"smooth":    s.controller.sensors.SmoothWind(s.controller.sensorID, 8),
	})
}

func (s *Server) handleJournal(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		_ = s.controller.journal.Prune(16)
	}
	events, err := s.controller.journal.Replay()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"events": events,
		"handles": s.controller.journal.OpenHandleCount(),
		"files":   s.controller.journal.ListFiles(),
	})
}

func (s *Server) handleFault(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"fault": string(s.controller.safe.Fault())})
}

func (s *Server) handleDiagnostics(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.controller.Diagnostics())
}

func (s *Server) handleAck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	s.controller.alarm.Acknowledge(r.FormValue("id"))
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleNormalShutdown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	if err := s.controller.safe.NormalShutdown(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleFeatherConcurrent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	target := parseFloat(r.FormValue("target"))
	if err := s.controller.pitch.FeatherAllConcurrently(target); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleFeatherLimited(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	target := parseFloat(r.FormValue("target"))
	step := parseFloat(r.FormValue("step"))
	if err := s.controller.pitch.FeatherLimited(target, step); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleMarkUnhealthy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	s.controller.sensors.MarkUnhealthy(r.FormValue("id"))
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleClearWarnings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	cleared := s.controller.alarm.ClearAll(alarm.LevelWarning)
	writeJSON(w, http.StatusOK, map[string]int{"cleared": cleared})
}

func (s *Server) handleTelemetry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	status := s.controller.Status()
	frame := model.NewTelemetryFrame(s.controller.Sequence(), s.controller.TurbineMetrics(), status.Protection, status.PitchState, status.YawState)
	data, err := frame.Marshal()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	decoded, err := model.UnmarshalTelemetryFrame(data)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, decoded)
}

func (s *Server) handleResetMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	s.controller.ResetMetrics()
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleSample(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	wind := parseFloat(r.FormValue("wind"))
	rotor := parseFloat(r.FormValue("rotor"))
	s.controller.FeedSample(wind, rotor)
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleEvaluate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	protection := s.controller.Evaluate()
	writeJSON(w, http.StatusOK, map[string]string{"protection": string(protection)})
}

func (s *Server) handlePitch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	if err := s.controller.ApplyPitch(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleFeather(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	target := parseFloat(r.FormValue("target"))
	if err := s.controller.pitch.FeatherTo(target); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleDefeather(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	if err := s.controller.pitch.Defeather(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleTick(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	direction := parseFloat(r.FormValue("direction"))
	writeJSON(w, http.StatusOK, s.controller.Tick(direction))
}

func (s *Server) handleYaw(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	target := parseFloat(r.FormValue("target"))
	if err := s.controller.YawTo(target); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	if err := s.controller.Stop(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleRecover(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	if err := s.controller.Recover(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, s.webRoot+"/web/console.html")
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func parseFloat(raw string) float64 {
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0
	}
	return value
}
