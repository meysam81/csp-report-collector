package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/meysam81/x/logging"
	"github.com/redis/go-redis/v9"
)

var (
	router      *chi.Mux
	redisClient *redis.Client
	appstate    *AppState
	logger      = logging.NewLogger()

	recorder = httptest.NewRecorder()
)

func TestMain(m *testing.M) {
	router = chi.NewRouter()
	redisClient, _ = NewRedis(context.TODO(), &Config{
		RedisHost: "localhost",
		RedisPort: 6379,
	})
	appstate = &AppState{
		redisClient: redisClient,
		logger:      &logger,
	}

	os.Exit(m.Run())
}

func TestPostCSPViolationReportURI(t *testing.T) {
	router.Post("/", appstate.ReceiverCSPViolation)

	var body bytes.Buffer
	body.Write([]byte(`
{
  "csp-report": {
    "blocked-uri": "http://example.com/css/style.css",
    "disposition": "report",
    "document-uri": "http://example.com/signup.html",
    "effective-directive": "style-src-elem",
    "original-policy": "default-src 'none'; style-src cdn.example.com; report-uri /_/csp-reports",
    "referrer": "",
    "status-code": 200,
    "violated-directive": "style-src-elem"
  }
}
	`))
	req, _ := http.NewRequest(http.MethodPost, "/", &body)
	req.Header.Set("content-type", "application/csp-report; charset=utf-8")

	router.ServeHTTP(recorder, req)

	if recorder.Result().StatusCode != http.StatusNoContent {
		r, _ := io.ReadAll(recorder.Body)
		fmt.Println(string(r))
		t.Fatalf("expected %d, but got %s", http.StatusNoContent, recorder.Result().Status)
	}
}

func TestPostCSPViolationReportTo(t *testing.T) {
	router.Post("/", appstate.ReceiverCSPViolation)

	var body bytes.Buffer
	body.Write([]byte(`
{
  "age": 53531,
  "body": {
    "blockedURL": "inline",
    "columnNumber": 39,
    "disposition": "enforce",
    "documentURL": "https://example.com/csp-report",
    "effectiveDirective": "script-src-elem",
    "lineNumber": 121,
    "originalPolicy": "default-src 'self'; report-to csp-endpoint-name",
    "referrer": "https://www.google.com/",
    "sample": "console.log(\"lo\")",
    "sourceFile": "https://example.com/csp-report",
    "statusCode": 200
  },
  "type": "csp-violation",
  "url": "https://example.com/csp-report",
  "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36"
}
	`))
	req, _ := http.NewRequest(http.MethodPost, "/", &body)
	req.Header.Set("content-type", "application/reports+json; charset=utf-8")

	router.ServeHTTP(recorder, req)

	if recorder.Result().StatusCode != http.StatusNoContent {
		r, _ := io.ReadAll(recorder.Body)
		fmt.Println(string(r))
		t.Fatalf("expected %d, but got %s", http.StatusNoContent, recorder.Result().Status)
	}
}
