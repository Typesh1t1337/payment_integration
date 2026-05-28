package model

import (
	"payment_integration/internal/a_order"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Status    a_order.OrderStatus
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	LockedUntil *time.Time `db:"locked_until"`
}
