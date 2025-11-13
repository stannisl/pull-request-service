package handlers

import (
	"encoding/json"
	"net/http"
)

type ErrorCode string

const (
	ErrorTeamExists   ErrorCode = "TEAM_EXISTS"
	ErrorPRExists     ErrorCode = "PR_EXISTS"
	ErrorPRMerged     ErrorCode = "PR_MERGED"
	ErrorNotAssigned  ErrorCode = "NOT_ASSIGNED"
	ErrorNoCandidate  ErrorCode = "NO_CANDIDATE"
	ErrorNotFound     ErrorCode = "NOT_FOUND"
	ErrorUnknown      ErrorCode = "UNKNOWN"
	ErrorUnauthorized ErrorCode = "UNAUTHORIZED"
)

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, code ErrorCode, message string) {
	writeJSON(w, status, errorEnvelope{
		Error: errorBody{
			Code:    code,
			Message: message,
		},
	})
}
