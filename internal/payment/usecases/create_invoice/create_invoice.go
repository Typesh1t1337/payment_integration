package create_invoice

import (
	"context"
	order_model "payment_integration/internal/a_order/model"
	invoice_model "payment_integration/internal/payment/models"

	"github.com/google/uuid"
)

type InvoiceRepository interface {
	Create(ctx context.Context, invoice invoice_model.Invoice) (*invoice_model.Invoice, error)
}

type OrderRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*order_model.Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status order_model.OrderStatus) error
}

type CreateInvoiceUseCase struct {
	invoiceRepository InvoiceRepository
	orderRepository OrderRepository
}

func NewCreateInvoiceUseCase(invoiceRepository InvoiceRepository, orderRepository OrderRepository) *CreateInvoiceUseCase {
	return &CreateInvoiceUseCase{
		invoiceRepository: invoiceRepository,
		orderRepository: orderRepository,
	}
}

type CreateInvoiceRequest struct {
	OrderID uuid.UUID `json:"order_id"`
}

type CreateInvoiceResponse struct {}

func (uc *CreateInvoiceUseCase) Execute(ctx context.Context, request CreateInvoiceRequest) (*CreateInvoiceResponse, error) {
	order, err := uc.orderRepository.GetByID(ctx, request.OrderID)
	if err != nil {
		return nil, err
	}

	if order.Status

	invoice, err := uc.invoiceRepository.Create(ctx, invoice_model.Invoice{
		OrderID: request.OrderID,
	})
	if err != nil {
		return nil, err
	}
	return &CreateInvoiceResponse{
		InvoiceID: invoice.ID,
	}, nil
}