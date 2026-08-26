package main

import (
	"os"
	"strconv"
	"time"
)

type runtimeConfig struct {
	Addr        string
	WebRoot     string
	JournalDir  string
	TickInterval time.Duration
	HistoryCap  int
}

func loadConfig() runtimeConfig {
	cfg := runtimeConfig{
		Addr:         envString("TURBINE_ADDR", ":8080"),
		WebRoot:      envString("TURBINE_WEB_ROOT", "."),
		JournalDir:   envString("TURBINE_JOURNAL_DIR", os.TempDir()),
		TickInterval: time.Duration(envInt("TURBINE_TICK_MS", 1000)) * time.Millisecond,
		HistoryCap:   envInt("TURBINE_HISTORY_CAP", 128),
	}
	return cfg
}

func envString(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
