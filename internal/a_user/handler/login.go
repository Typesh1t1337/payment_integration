package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"payment_integration/internal/a_user"
	"payment_integration/internal/a_user/usecases/login"
	"payment_integration/internal/domain"
	"payment_integration/internal/transport/http_transport"
)

func (h *Handler) Login() http.HandlerFunc {
	return http_transport.WrapHandler(h.logger, func(w http.ResponseWriter, r *http.Request) error {
		var dto login.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			return &http_transport.ErrorResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Request body is invalid",
			}
		}
		if err := h.validate.Struct(dto); err != nil {
			return err
		}
		tokens, err := h.loginUseCase.Execute(r.Context(), dto)
		if err != nil {
			if errors.Is(err, a_user.ErrInvalidPassword) {
				return &http_transport.ErrorResponse{
					StatusCode: http.StatusUnauthorized,
					Message:    "Invalid email or password",
				}
			}
			if errors.Is(err, domain.ErrNotFound) {
				return &http_transport.ErrorResponse{
					StatusCode: http.StatusUnauthorized,
					Message:    "Invalid email or password",
				}
			}
			return err
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    tokens.RefreshToken,
			HttpOnly: true,
			Secure:   h.cfg.Env == "prod",
			SameSite: http.SameSiteStrictMode,
			MaxAge:   int(h.cfg.Jwt.RefreshTTL.Seconds()),
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"access_token": tokens.AccessToken,
		}); err != nil {
			h.logger.Error("Error encoding response", "error", err)
		}
		return nil
	})
}