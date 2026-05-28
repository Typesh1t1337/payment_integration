package repository

import (
	"context"
	"errors"
	"fmt"
	"payment_integration/internal/a_user"
	"payment_integration/internal/a_user/model"
	"payment_integration/internal/domain"
	"payment_integration/internal/uow"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, createUser a_user.CreateUser) (*model.User, error) {
	executor := uow.GetExecutor(ctx, r.db)
	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, password, created_at, updated_at
	`
	row := executor.QueryRow(ctx, query, createUser.Name, createUser.Email, createUser.HashedPassword)
	var user model.User
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, domain.ErrAlreadyExists
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &user, nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	executor := uow.GetExecutor(ctx, r.db)
	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	row := executor.QueryRow(ctx, query, id)
	return r.parseUser(row, "GetByID")
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	executor := uow.GetExecutor(ctx, r.db)
	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	row := executor.QueryRow(ctx, query, email)
	return r.parseUser(row, "GetByEmail")
}

func (r *PostgresUserRepository) parseUser(row pgx.Row, method string) (*model.User, error) {
	var user model.User
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, fmt.Errorf("failed to parse user in %s: %w", method, err)
	}
	return &user, nil
}