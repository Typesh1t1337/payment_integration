package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"payment_integration/internal/a_user"
	"payment_integration/internal/a_user/usecases/refresh"
	"payment_integration/internal/config"
	"payment_integration/internal/transport/http_transport"
)

func (h *Handler) Refresh() http.HandlerFunc {
	return http_transport.WrapHandler(h.logger, func(w http.ResponseWriter, r *http.Request) error {
		var request refresh.RefreshRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil && !errors.Is(err, io.EOF) {
			return &http_transport.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Request body is invalid",
			}
		}

		if request.RefreshToken == "" {
			refreshToken, err := r.Cookie("refresh_token")
			if err == nil {
				request.RefreshToken = refreshToken.Value
			}
		}

		if request.RefreshToken == "" {
			return &http_transport.ErrorResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    "Unauthorized",
			}
		}
		response, err := h.refreshUseCase.Execute(r.Context(), &request)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    "",
				HttpOnly: true,
				Secure:   h.cfg.Env == config.EnvProd,
				SameSite: http.SameSiteStrictMode,
				MaxAge:   -1,
			})
			if errors.Is(err, a_user.ErrInvalidToken) {
				return &http_transport.ErrorResponse{
					StatusCode: http.StatusUnauthorized,
					Message:    "Unauthorized",
				}
			}
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			h.logger.Error("Error encoding response", "error", err)
		}
		return nil
	})
}