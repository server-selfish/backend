package pkg

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Response struct {
	Message any    `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(v); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func ReturnError(w http.ResponseWriter, err error) {
	var httpStatus int
	msg := err.Error()

	switch {
	case errors.Is(err, ErrAlreadyExist):
		httpStatus = http.StatusConflict
	case errors.Is(err, ErrNotFound):
		httpStatus = http.StatusNotFound
	default:
		httpStatus = http.StatusInternalServerError
		msg = "internal server error"
	}

	WriteJSON(w, httpStatus, Response{Message: msg})
}

func ReturnSuccess(w http.ResponseWriter, statusCode int, message string, data any) {
	WriteJSON(w, statusCode, Response{Message: message, Data: data})
}
