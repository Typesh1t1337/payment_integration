package handler

import (
	"log/slog"
	"payment_integration/internal/config"
	"payment_integration/internal/payment/usecases/create_invoice"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	createInvoiceUseCase create_invoice.CreateInvoiceUseCase
	logger *slog.Logger
	validate *validator.Validate
	cfg *config.Config
}

func NewHandler(createInvoiceUseCase create_invoice.CreateInvoiceUseCase, logger *slog.Logger, validate *validator.Validate, cfg *config.Config) *Handler {
	return &Handler{
		createInvoiceUseCase: createInvoiceUseCase,
		logger: logger,
		validate: validate,
		cfg: cfg,
	}
}