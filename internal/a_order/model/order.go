package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	id        uuid.UUID  `db:"id"`
	userID    uuid.UUID  `db:"user_id"`
	Status    string     `db:"status"`
	CreatedAt *time.Time `db:"created_at"`
}
