package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"payment_integration/internal/a_order/usecase"
)

type Handler struct {
	AddItemUseCase usecase.AddItemUseCase

	logger *slog.Logger
}

func NewHandler(addItemUseCase usecase.AddItemUseCase, logger *slog.Logger) *Handler {
	return &Handler{
		AddItemUseCase: addItemUseCase,
		logger:         logger,
	}
}

func (h *Handler) RegisterRoutes(r *http.ServeMux, prefix string) {
	r.HandleFunc(fmt.Sprintf("POST %s", prefix), h.AddItemHandler())
}
