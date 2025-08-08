package main

import (
	"net/http"

	"github.com/meysam81/x/logging"
	"github.com/redis/go-redis/v9"
)

type CSPReport struct {
	Age       int    `json:"age"`
	Body      *Body  `json:"body"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	UserAgent string `json:"user_agent"`
}

type Body struct {
	BlockedURL         string `json:"blockedURL"`
	ColumnNumber       int    `json:"columnNumber"`
	Disposition        string `json:"disposition"`
	DocumentURL        string `json:"documentURL"`
	EffectiveDirective string `json:"effectiveDirective"`
	LineNumber         int    `json:"lineNumber"`
	OriginalPolicy     string `json:"originalPolicy"`
	Referrer           string `json:"referrer"`
	Sample             string `json:"sample"`
	SourceFile         string `json:"sourceFile"`
	StatusCode         int    `json:"statusCode"`
}

type AppState struct {
	redisClient *redis.Client
	logger      *logging.Logger
	handler     *http.Server
	config      *Config
}
