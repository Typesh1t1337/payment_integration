package payment

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateInvoiceDTO struct {
	UserID      uuid.UUID       `json:"user_id"`
	OrderID     uuid.UUID       `json:"order_id"`
	Status      InvoiceStatus   `json:"status"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	ExpiresAt   time.Time       `json:"expires_at"`
}