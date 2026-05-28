package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderItems struct {
	id        uuid.UUID  `db:"id"`
	OrderID   uuid.UUID  `db:"order_id"`
	ProductID uuid.UUID  `db:"product_id"`
	Quantity  int        `db:"quantity"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
