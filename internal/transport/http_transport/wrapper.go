package http_transport

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func WrapHandler(logger *slog.Logger, handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			var validationErr validator.ValidationErrors
			if errors.As(err, &validationErr) {
				fields := make(map[string]string)
				for _, err := range validationErr {
					fields[err.Field()] = err.Tag()
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"message":     "Validation error",
					"fields":      fields,
				})
				return
			}

			var errorResponse *ErrorResponse
			if errors.As(err, &errorResponse) {
				errorResponse.Write(w)
				return
			}

			logger.Error("Error handling request", "error", err)
			errorResponse = &ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal server error",
			}
			errorResponse.Write(w)
			return
		}
	}
}