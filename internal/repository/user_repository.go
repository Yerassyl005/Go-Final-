package repository

import (
	"context"
	"errors"

	"myapp-auth/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (first_name, last_name, phone, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	return r.DB.QueryRow(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.PasswordHash,
	).Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, phone, password_hash, created_at
		FROM users
		WHERE phone = $1
	`

	var user models.User
	err := r.DB.QueryRow(ctx, query, phone).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, phone, password_hash, created_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func IsDuplicateKey(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
