package repo

import (
	"context"
	"payment_integration/internal/a_order"
	"payment_integration/internal/a_order/model"
	"payment_integration/internal/uow"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderItemRepository struct {
	db *pgxpool.Pool
}

func NewOrderItemRepository(db *pgxpool.Pool) *OrderItemRepository {
	return &OrderItemRepository{db: db}
}

func (r *OrderItemRepository) Create(ctx context.Context, value a_order.AddOrderItem) (*model.OrderItems, error) {
	session := uow.GetExecutor(ctx, r.db)
	var orderItem model.OrderItems

	row := session.QueryRow(
		ctx,
		`INSERT INTO order_items 
    		(order_id, product_id, quantity)
			VALUES ($1, $2, $3)
			ON CONFLICT (order_id, product_id) DO UPDATE
			SET quantity = quantity + $3
			RETURNING id, order_id, product_id, quantity, created_at, updated_at`,
		value.OrderID, value.ProductID, value.Quantity,
	)

	err := row.Scan(
		&orderItem.ID, &orderItem.OrderID, &orderItem.ProductID,
		&orderItem.Quantity, &orderItem.CreatedAt, &orderItem.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &orderItem, nil
}
