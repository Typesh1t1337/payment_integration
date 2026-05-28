package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"payment_integration/internal/a_user"
	"payment_integration/internal/a_user/usecases/register"
	"payment_integration/internal/transport/http_transport"
)

func (h *Handler) Register() http.HandlerFunc {
	return http_transport.WrapHandler(h.logger, func(w http.ResponseWriter, r *http.Request) error {
		var dto register.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			return &http_transport.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Request body is invalid",
			}
		}
		if err := h.validate.Struct(dto); err != nil {
			return err
		}
		response, err := h.registerUseCase.Execute(r.Context(), dto)
		if err != nil {
			if errors.Is(err, a_user.ErrUserAlreadyExists) {
				return &http_transport.ErrorResponse{
					StatusCode: http.StatusConflict,
					Message:    "User already exists",
				}
			}
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			h.logger.Error("Error encoding response", "error", err)
		}
		return nil
	})
}