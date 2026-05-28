package repo

import (
	"context"
	"payment_integration/internal/uow"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Product struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *Product {
	return &Product{db: db}
}

func (r *Product) Exists(ctx context.Context, productID uuid.UUID) bool {
	session := uow.GetExecutor(ctx, r.db)
	var exists bool

	row := session.QueryRow(ctx,
		`SELECT 
			 EXISTS(
				SELECT 1
               FROM products
               WHERE products.id = $1
			)
		`,
		productID,
	)

	if err := row.Scan(&exists); err != nil {
		return false
	}

	return exists
}
