package main

import (
	"net/http"

	"github.com/meysam81/x/logging"
	"github.com/redis/go-redis/v9"
)

// modern CSP report structure with content-type application/reports+json
type ReportTo struct {
	Age  int `json:"age"`
	Body struct {
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
	} `json:"body"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	UserAgent string `json:"user_agent"`
}

// legacy report structure with content-type application/csp-report
type ReportURI struct {
	CSPReport struct {
		BlockedURI         string `json:"blocked-uri"`
		Disposition        string `json:"disposition"`
		DocumentURI        string `json:"document-uri"`
		EffectiveDirective string `json:"effective-directive"`
		OriginalPolicy     string `json:"original-policy"`
		Referrer           string `json:"referrer"`
		StatusCode         int    `json:"status-code"`
		ViolatedDirective  string `json:"violated-directive"`
	} `json:"csp-report"`
}

type AppState struct {
	redisClient *redis.Client
	logger      *logging.Logger
	handler     *http.Server
	config      *Config
}
