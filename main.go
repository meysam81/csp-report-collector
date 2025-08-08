package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	defer stop()

	config, err := NewConfig()
	if err != nil {
		log.Fatalf("failed initialing config: %s", err)
	}

	ctxS, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	app, err := NewApp(ctx, config)
	if err != nil {
		log.Fatalf("failed initialing app: %s", err)
	}

	go func() {
		defer cancel()
		app.logger.Info().Msgf("listening on address: %s", app.handler.Addr)
		if err := app.handler.ListenAndServe(); err != http.ErrServerClosed {
			app.logger.Error().Err(err).Msg("failed shutting down the server")
		}
	}()

	<-ctx.Done()
	stop()

	app.logger.Info().Msg("shutdown signal received. waiting for server to close.")

	<-ctxS.Done()
	app.logger.Info().Msg("shutdown complete. bye bye.")
}
