package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Invoice struct {
	ID          uuid.UUID       `db:"id"`
	UserID      uuid.UUID       `db:"user_id"`
	OrderID     uuid.UUID       `db:"order_id"`
	Status      string          `db:"status"`
	TotalAmount decimal.Decimal `db:"total_amount"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at"`
	PaidAt      time.Time       `db:"paid_at"`
	ExpiresAt   time.Time       `db:"expires_at"`
}
