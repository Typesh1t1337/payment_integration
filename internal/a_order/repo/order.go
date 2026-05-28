package repo

import (
	"context"
	"errors"
	"fmt"
	"payment_integration/internal/a_order"
	"payment_integration/internal/a_order/model"
	"payment_integration/internal/domain"
	"payment_integration/internal/uow"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
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
			RETURNING id, user_id, status, created_at, updated_at, locked_until
		`,
		userId, a_order.OrderStatusCreated,
	)

	err := row.Scan(&orderModel.ID, &orderModel.UserID, &orderModel.Status, &orderModel.CreatedAt, &orderModel.UpdatedAt, &orderModel.LockedUntil)

	if err == nil {
		return &orderModel, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, a_order.UnexpectedOrderError
	}

	getOrderRow := session.QueryRow(
		ctx,
		`UPDATE orders 
				SET locked_until = NULL,
					status = CASE 
						WHEN status = 'handling' AND locked_until < now() THEN 'created'
						ELSE status
					END
				WHERE user_id = $1
				  AND status IN ('handling', 'created')
				  AND (locked_until IS NULL OR locked_until < now())
				RETURNING id, user_id, status, created_at, updated_at, locked_until
			`,
		userId,
	)

	err = getOrderRow.Scan(
		&orderModel.ID,
		&orderModel.UserID,
		&orderModel.Status,
		&orderModel.CreatedAt,
		&orderModel.UpdatedAt,
		&orderModel.LockedUntil,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, a_order.OrderBeingProcessedError
	}

	if err != nil {
		return nil, a_order.UnexpectedOrderError
	}

	return &orderModel, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status a_order.OrderStatus, lockedUntil *time.Time) (*model.Order, error) {
	session := uow.GetExecutor(ctx, r.db)
	row := session.QueryRow(ctx, `
	UPDATE orders 
	SET status = $1, locked_until = $2 WHERE id = $3 
	returning id, user_id, status, created_at, updated_at, locked_until
	`, status, lockedUntil, id)
	
	var orderModel model.Order
	err := row.Scan(&orderModel.ID, &orderModel.UserID, &orderModel.Status, &orderModel.CreatedAt, &orderModel.UpdatedAt, &orderModel.LockedUntil)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}
	return &orderModel, nil
}

func (r *OrderRepository) GetTotalAmount(ctx context.Context, orderID uuid.UUID) (decimal.Decimal, error) {
	session := uow.GetExecutor(ctx, r.db)
	row := session.QueryRow(ctx, `
	SELECT COALESCE(SUM(p.price * oi.quantity), 0)
	FROM order_items oi
	JOIN products p ON oi.product_id = p.id
	WHERE oi.order_id = $1
	`, orderID)
	var totalAmount decimal.Decimal
	err := row.Scan(&totalAmount)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get total amount: %w", err)
	}
	if totalAmount.IsZero() {
		return decimal.Zero, domain.ErrNotFound
	}
	return totalAmount, nil
}