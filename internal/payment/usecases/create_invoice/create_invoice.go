package create_invoice

import (
	"context"
	"payment_integration/internal/a_order"
	"payment_integration/internal/a_order/model"
	"payment_integration/internal/payment"
	invoice_model "payment_integration/internal/payment/models"
	"payment_integration/internal/uow"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InvoiceRepository interface {
	Create(ctx context.Context, invoice *payment.CreateInvoiceDTO) (*invoice_model.Invoice, error)
}

type OrderRepository interface {
	UpdateStatus(ctx context.Context, id uuid.UUID, status a_order.OrderStatus, lockedUntil *time.Time) (*model.Order, error)
	GetTotalAmount(ctx context.Context, orderID uuid.UUID) (decimal.Decimal, error)	
}

type CreateInvoiceUseCase struct {
	invoiceRepository InvoiceRepository
	orderRepository OrderRepository
	uow uow.UoW
}

func NewCreateInvoiceUseCase(invoiceRepository InvoiceRepository, orderRepository OrderRepository, uow uow.UoW) *CreateInvoiceUseCase {
	return &CreateInvoiceUseCase{
		invoiceRepository: invoiceRepository,
		orderRepository: orderRepository,
		uow: uow,
	}
}

type CreateInvoiceRequest struct {
	OrderID uuid.UUID `json:"order_id"`
	PublicTerminalId string `json:"publicId"`
	AccountId string `json:"accountId"`
}

type CreateInvoiceResponse struct {
	PublicTerminalId string `json:"publicId"`
	Description string `json:"description"`
	PaymentSchema string `json:"paymentSchema"`
	Currency string `json:"currency"`
	Amount int64 `json:"amount"`
	ExternalId string `json:"externalId"`
	AccountId string `json:"accountId"`
	ApplePaySupport bool `json:"applePaySupport"`
	GooglePaySupport bool `json:"googlePaySupport"`
	Language string `json:"language"`
	RequireEmail bool `json:"requireEmail"`
	Data *map[string]any `json:"data"`
}

func (uc *CreateInvoiceUseCase) Execute(ctx context.Context, request CreateInvoiceRequest) (*CreateInvoiceResponse, error) {
	response, err := uow.Do(ctx, uc.uow, func(ctx context.Context) (*CreateInvoiceResponse, error) {
		lockUntil := time.Now().Add(time.Minute * 5)
		order, err := uc.orderRepository.UpdateStatus(ctx, request.OrderID, a_order.OrderStatusHandling, &lockUntil)
		if err != nil {
			return nil, err
		}
		totalAmount, err := uc.orderRepository.GetTotalAmount(ctx, request.OrderID)
		if err != nil {
			return nil, err
		}
		invoice, err := uc.invoiceRepository.Create(ctx, &payment.CreateInvoiceDTO{
			OrderID: order.ID,
			Status: payment.InvoiceStatusPending,
			TotalAmount: totalAmount,
			ExpiresAt: lockUntil,
		})
		if err != nil {
			return nil, err
		}
		return &CreateInvoiceResponse{
			PublicTerminalId: request.PublicTerminalId,
			Description: "",
			PaymentSchema: "Single",
			Currency: string(payment.CurrencyKZT),
			Amount: totalAmount.IntPart(),
			ExternalId: invoice.ID.String(),
			AccountId: request.AccountId,
			ApplePaySupport: false,
			GooglePaySupport: false,
			Language: "ru-RU",
			RequireEmail: false,
			Data: nil,
		}, nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}