// Copyright (c) Illumio, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// TokenRequest is a struct to hold the request parameters for the authenticateHandler
// following the OAuth2.0 specification.
type TokenRequest struct {
	GrantType    string
	ClientID     string
	ClientSecret string
}

// TokenResponse is a struct to hold the response parameters for the authenticateHandler
// following the OAuth2.0 specification.
type TokenResponse struct {
	AccessToken string `json:"access_token,omitempty"` //nolint:tagliatelle
}

const (
	AllowedGrantType   = "client_credentials"
	InvalidGrantError  = "invalid_grant"
	UnauthorizedClient = "unauthorized_client"
)

type AuthService struct {
	logger       *zap.Logger
	clientID     string
	clientSecret string
	token        string
}

func (a *AuthService) authenticateHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Info(
		"received request",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
	)

	if r.Method != http.MethodPost {
		a.logger.Error("Invalid request method, method not allowed", zap.String("method", r.Method))
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)

		return
	}

	var req TokenRequest

	err := r.ParseForm()
	if err != nil {
		a.logger.Error("Invalid request, unable to parse form", zap.Error(err))
		http.Error(w, "Invalid request", http.StatusBadRequest)

		return
	}

	req.GrantType = r.FormValue("grant_type")
	req.ClientID = r.FormValue("client_id")
	req.ClientSecret = r.FormValue("client_secret")

	if req.GrantType != AllowedGrantType {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": InvalidGrantError})

		return
	}

	if req.ClientID == a.clientID && req.ClientSecret == a.clientSecret {
		response := TokenResponse{AccessToken: a.token}
		jsonResponse(w, http.StatusOK, response)
	} else {
		jsonResponse(w, http.StatusUnauthorized, map[string]string{"error": UnauthorizedClient})
	}
}

func jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if encodeError := json.NewEncoder(w).Encode(data); encodeError != nil {
		http.Error(w, encodeError.Error(), http.StatusInternalServerError)

		return
	}
}

func startHTTPServer(address string, loggerToUse *zap.Logger, clientID string, clientSecret string, token string) {
	authService := &AuthService{
		logger:       loggerToUse,
		clientID:     clientID,
		clientSecret: clientSecret,
		token:        token,
	}
	http.HandleFunc("/api/v1/authenticate", authService.authenticateHandler)

	server := &http.Server{
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
