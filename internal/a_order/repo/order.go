package repo

import (
	"context"
	"errors"
	"payment_integration/internal/a_order"
	"payment_integration/internal/a_order/model"
	"payment_integration/internal/uow"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetOrCreate(ctx context.Context, userId uuid.UUID) (*model.Order, error) {
	session := uow.GetExecutor(ctx, r.db)
	var orderModel model.Order

	row := session.QueryRow(ctx,
		`INSERT INTO 
			orders (user_id, status) 
			VALUES ($1, $2)
			ON CONFLICT (user_id) WHERE status IN ('created', 'handling')
			RETURNING id, user_id, status, created_at, updated_at
		`,
		userId, a_order.OrderStatusCreated,
	)

	err := row.Scan(&orderModel.ID, &orderModel.UserID, &orderModel.Status, &orderModel.CreatedAt)

	if err == nil {
		return &orderModel, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, a_order.UnexpectedOrderError
	}

	getOrderRow := session.QueryRow(
		ctx,
		`SELECT 
					id, 
					user_id,
					status,
					created_at,
					updated_at
				FROM orders
				user_id = $1 AND status IN ('created', 'handling')
				LIMIT 1
				FOR UPDATE
			`,
		userId,
	)

	err = getOrderRow.Scan(
		&orderModel.ID,
		&orderModel.UserID,
		&orderModel.Status,
		&orderModel.CreatedAt,
		*orderModel.UpdatedAt,
	)

	if err != nil {
		return nil, a_order.UnexpectedOrderError
	}

	if orderModel.Status == a_order.OrderStatusHandling {
		return nil, a_order.OrderBeingProcessedError
	}

	return &orderModel, nil
}
