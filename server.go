package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/meysam81/x/chimux"
	"github.com/meysam81/x/logging"
)

func NewLogger(logLevel string) *logging.Logger {
	l := logging.NewLogger(logging.WithLogLevel(logLevel))
	return &l
}

func NewApp(ctx context.Context, c *Config) (*AppState, error) {
	root := chimux.NewChi()
	mw := chimux.NewChi(chimux.WithLoggingMiddleware())
	api := chimux.NewChi(chimux.WithHealthz(), chimux.WithMetrics())

	root.Mount("/", mw)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      root,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	redisClient, err := NewRedis(ctx, c)
	if err != nil {
		return nil, err
	}

	logger := NewLogger(c.LogLevel)

	app := &AppState{
		redisClient: redisClient,
		handler:     s,
		logger:      logger,
		config:      c,
	}

	mw.Use(app.RateLimitMiddleware)
	mw.Mount("/", api)
	api.Post("/", app.ReceiverCSPViolation)

	return app, nil
}
