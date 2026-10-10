package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Adega318/urlittle/internal/config"
	"github.com/Adega318/urlittle/internal/handlers"
	"github.com/Adega318/urlittle/internal/store"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	st, err := store.NewStore(ctx, config.Store.DBURL, config.Store.CacheSize)
	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()

	go runPeriodic(ctx, time.Minute, st.ClearExpired)

	h := handlers.NewHandler(st)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /{$}", h.Store)
	mux.HandleFunc("GET /{id}", h.Redirect)
	mux.HandleFunc("GET /health", h.Health)

	srv := &http.Server{
		Addr:              ":" + config.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	log.Printf("server listening on :%s", config.Port)

	select {
	case err := <-errCh:
		log.Fatal(err)
	case <-ctx.Done():
	}

	log.Printf("server stopping...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}

func runPeriodic(ctx context.Context, interval time.Duration, fn func(context.Context) error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := fn(ctx); err != nil {
				log.Printf("periodic task failed: %v", err)
			}
		}
	}
}
