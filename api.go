package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/goccy/go-json"
)

var (
	ErrBadContentType = []byte(`{"status":"failed","message":"Bad content-type provided. Only application/reports+json is acceptable."}`)
	ErrInvalidBody    = []byte(`{"status":"failed","message":"Invalid body provided. Only JSON-serializable strings are acceptable."}`)
	ErrInternalError  = []byte(`{"status":"failed","message":"Failed processing your request. Please try again or contacty administrator."}`)
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
	var cspreport []byte

	switch contentType {
	case "application/reports+json":
		reportto := &ReportTo{}
		err := json.NewDecoder(r.Body).Decode(reportto)
		if err != nil {
			a.logger.Error().Err(err).Msg("failed decoding request body")
			a.respondWithInterface(w, ErrInvalidBody, http.StatusBadRequest)
			return
		}

		cspreport, err = json.Marshal(reportto)
		if err != nil {
			a.logger.Error().Err(err).Msg("failed encoding the request body for save")
			a.respondWithInterface(w, ErrInternalError, http.StatusInternalServerError)
			return
		}

	case "application/csp-report":
	case "application/json":
		reporturi := &ReportURI{}
		err := json.NewDecoder(r.Body).Decode(reporturi)
		if err != nil {
			a.logger.Error().Err(err).Msg("failed decoding request body")
			a.respondWithInterface(w, ErrInvalidBody, http.StatusBadRequest)
			return
		}

		cspreport, err = json.Marshal(reporturi)
		if err != nil {
			a.logger.Error().Err(err).Msg("failed encoding the request body for save")
			a.respondWithInterface(w, ErrInternalError, http.StatusInternalServerError)
			return
		}
	default:
		a.logger.Error().Str("content_type", contentType).Msg("invalid content type rejected.")
		a.respondWithInterface(w, ErrBadContentType, http.StatusBadRequest)
		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			a.logger.Error().Err(err).Msg("failed closing request body.")
		}
	}()

	a.logger.Info().Bytes("csp_report", cspreport).Msg("received a csp violation report")

	now := fmt.Sprintf("%d", time.Now().Unix())
	_, err := a.redisClient.Set(r.Context(), now, cspreport, 0).Result()
	if err != nil {
		a.logger.Error().Err(err).Msg("failed saving body to redis.")
	}

	w.WriteHeader(http.StatusNoContent)
}
