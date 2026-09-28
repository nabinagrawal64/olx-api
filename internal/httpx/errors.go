package httpx

import (
	"encoding/json"
	"net/http"
)

type Code string

const (
	CodeInvaidID Code = "invalid_id"
	CodeNotFound Code = "not_found"
	CodeInternalError Code = "internal_server_error"
	CodeInvalidRequest Code = "invalid_request"
	CodeUnauthorized Code = "unauthorized"
	CodeForbidden Code = "forbidden"
	CodeConflict Code = "conflict"
	CodeRateLimitExceeded Code = "rate_limit_exceeded"
	CodeBadRequest Code = "bad_request"
	CodeServiceUnavailable Code = "service_unavailable"
	CodeTimeout Code = "timeout"
)

type errorEnvelope struct {
	Error errorPayload `json:"error"`
}

type errorPayload struct {
	Code    Code `json:"code"`
	Message string `json:"message"`
}

func Error(w http.ResponseWriter, statusCode int, message string, code Code) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(errorEnvelope{
		Error: errorPayload{
			Code:    code,
			Message: message,
		},
	})
}
