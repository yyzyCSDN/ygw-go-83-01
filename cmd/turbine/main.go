package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := loadConfig()
	controller := NewController(cfg.JournalDir)
	server := NewServer(controller, cfg.WebRoot)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		ticker := time.NewTicker(cfg.TickInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_ = controller.Tick(0)
			}
		}
	}()

	log.Printf("wind turbine control listening on %s", cfg.Addr)
	serverInstance := &http.Server{Addr: cfg.Addr, Handler: server.Handler()}
	go func() {
		<-ctx.Done()
		_ = serverInstance.Shutdown(context.Background())
	}()
	if err := serverInstance.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
