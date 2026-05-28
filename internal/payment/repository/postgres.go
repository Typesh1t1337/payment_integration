package repository

import (
	"context"
	"errors"
	"fmt"
	"payment_integration/internal/a_payment/model"
	"payment_integration/internal/domain"
	"payment_integration/internal/payment"
	"payment_integration/internal/uow"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresInvoiceRepository struct {
	db *pgxpool.Pool
}

func NewPostgresInvoiceRepository(db *pgxpool.Pool) *PostgresInvoiceRepository {
	return &PostgresInvoiceRepository{db: db}
}

func (r *PostgresInvoiceRepository) Create(ctx context.Context, invoice *payment.CreateInvoiceDTO) (*model.Invoice, error) {
	executor := uow.GetExecutor(ctx, r.db)
	query := `
	INSERT INTO invoices (order_id, status, total_amount, expires_at)
	VALUES ($1, $2, $3, $4)
	RETURNING id, order_id, status, total_amount, expires_at
	`
	row := executor.QueryRow(ctx, query, invoice.OrderID, invoice.Status, invoice.TotalAmount, invoice.ExpiresAt)
	var invoiceModel model.Invoice
	err := row.Scan(&invoiceModel.ID, &invoiceModel.OrderID, &invoiceModel.Status, &invoiceModel.TotalAmount, &invoiceModel.ExpiresAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, domain.ErrAlreadyExists
		}
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}
	return &invoiceModel, nil
}