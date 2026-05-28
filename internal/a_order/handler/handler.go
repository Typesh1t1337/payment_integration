package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"payment_integration/internal/a_order/usecase"
	"payment_integration/internal/config"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	AddItemUseCase usecase.AddItemUseCase

	logger   *slog.Logger
	validate *validator.Validate
	cfg      *config.Config
}

func NewHandler(addItemUseCase usecase.AddItemUseCase, logger *slog.Logger, validate *validator.Validate, cfg *config.Config) *Handler {
	return &Handler{
		AddItemUseCase: addItemUseCase,
		logger:         logger,
		validate:       validate,
		cfg:            cfg,
	}
}

func (h *Handler) RegisterRoutes(r *http.ServeMux, prefix string) {
	r.HandleFunc(fmt.Sprintf("POST %s", prefix), h.AddItemHandler())
}
