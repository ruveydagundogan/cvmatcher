package response

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	apperrors "github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/errors"
)

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("json encode error", "error", err, "status", status)
	}
}

func Success(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

func SuccessWithMessage(w http.ResponseWriter, message string) {
	JSON(w, http.StatusOK, SuccessResponse{
		Success: true,
		Message: message,
	})
}

func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, err error) {
	status := apperrors.HTTPStatusCode(err)

	var domainErr *apperrors.DomainError
	message := "an internal error occurred"
	if errors.As(err, &domainErr) {
		message = domainErr.Message
		slog.Error("request error", "kind", domainErr.Kind, "message", domainErr.Message, "inner_error", domainErr.Err, "status", status)
	} else {
		slog.Error("request error", "error", err, "status", status)
	}

	JSON(w, status, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    status,
	})
}

func BadRequest(w http.ResponseWriter, message string) {
	JSON(w, http.StatusBadRequest, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    http.StatusBadRequest,
	})
}

func NotFound(w http.ResponseWriter, message string) {
	JSON(w, http.StatusNotFound, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    http.StatusNotFound,
	})
}

func Unauthorized(w http.ResponseWriter, message string) {
	JSON(w, http.StatusUnauthorized, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    http.StatusUnauthorized,
	})
}

func ForbiddenError(message string) error {
	return apperrors.Forbidden(message)
}

func InternalError(message string) error {
	return apperrors.Internal(message, nil)
}
