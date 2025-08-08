package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/goccy/go-json"
)

const (
	ExpectedContentType  = "application/reports+json"
	ExpectedContentType2 = "application/json"
)

var (
	ErrBadContentType = []byte(`{"status":"failed","message":"Bad content-type provided. Only application/reports+json is acceptable."}`)
	ErrInvalidBody    = []byte(`{"status":"failed","message":"Invalid body provided. Only JSON-serializable strings are acceptable."}`)
)

func (a *AppState) respondWithInterface(w http.ResponseWriter, obj interface{}, statusCode int) {
	w.WriteHeader(statusCode)
	if obj == nil {
		return
	}
	w.Header().Set("content-type", "application/json")

	switch v := obj.(type) {
	case []byte:
		_, err := w.Write(v)
		if err != nil {
			a.logger.Error().Err(err).Msg("failed writing the response body")
		}
		return
	}

	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed writing the response body")
	}
}

func (a *AppState) ReceiverCSPViolation(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("content-type")
	if contentType != ExpectedContentType && contentType != ExpectedContentType2 {
		a.logger.Error().Str("content_type", contentType).Msg("invalid content type rejected.")
		a.respondWithInterface(w, ErrBadContentType, http.StatusBadRequest)
		return
	}

	csp := &CSPReport{}
	err := json.NewDecoder(r.Body).Decode(csp)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed decoding request body.")
		a.respondWithInterface(w, ErrInvalidBody, http.StatusBadRequest)
		return
	}
	defer func() {
		err = r.Body.Close()
		if err != nil {
			a.logger.Error().Err(err).Msg("failed closing request body.")
		}
	}()

	a.logger.Info().Interface("csp_report", csp).Msg("received a csp violation report")

	parsedData, err := json.Marshal(csp)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed encoding CSP report for save.")
		a.respondWithInterface(w, ErrInvalidBody, http.StatusBadRequest)
		return
	}

	now := fmt.Sprintf("%d", time.Now().Unix())
	_, err = a.redisClient.Set(r.Context(), now, parsedData, 0).Result()
	if err != nil {
		a.logger.Error().Err(err).Msg("failed saving body to redis.")
	}

	w.WriteHeader(http.StatusNoContent)
}
