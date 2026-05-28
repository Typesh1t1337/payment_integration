package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"payment_integration/internal/a_user/usecases/login"
	"payment_integration/internal/a_user/usecases/refresh"
	"payment_integration/internal/a_user/usecases/register"
	"payment_integration/internal/config"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	loginUseCase    login.LoginUseCase
	registerUseCase register.RegisterUseCase
	refreshUseCase  refresh.RefreshUseCase
	logger *slog.Logger
	validate *validator.Validate
	cfg *config.Config
}

func NewHandler(loginUseCase login.LoginUseCase, registerUseCase register.RegisterUseCase, refreshUseCase refresh.RefreshUseCase, logger *slog.Logger, validate *validator.Validate, cfg *config.Config) *Handler {
	return &Handler{
		loginUseCase: loginUseCase,
		registerUseCase: registerUseCase,
		refreshUseCase: refreshUseCase,
		logger: logger,
		validate: validate,
		cfg: cfg,
	}
}

func (h *Handler) RegisterRoutes(r *http.ServeMux, prefix string) {
	r.HandleFunc(fmt.Sprintf("POST %s/login", prefix), h.Login())
	r.HandleFunc(fmt.Sprintf("POST %s/register", prefix), h.Register())
	r.HandleFunc(fmt.Sprintf("POST %s/refresh", prefix), h.Refresh())
}