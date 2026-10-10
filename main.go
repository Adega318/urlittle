package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Adega318/urlittle/internal/config"
	"github.com/Adega318/urlittle/internal/handlers"
	"github.com/Adega318/urlittle/internal/store"
	"github.com/lmittmann/tint"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := slog.New(tint.NewTextHandler(os.Stdout, &tint.Options{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	st, err := store.NewStore(ctx, config.Store.DBURL, config.Store.CacheSize)
	if err != nil {
		slog.Error("store connection error", "err", err)
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
		slog.Info("server started", "port", config.Port)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		slog.Error("server error", "err", err)
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("failed to shutdown", "err", err)
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
				slog.Error("periodic task failed", "err", err)
			}
		}
	}
}
