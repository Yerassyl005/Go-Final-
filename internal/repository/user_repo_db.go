package repository

import (
	"database/sql"
	"errors"
	"smartqueue/internal/models"
)

type UserPostgresRepository struct {
	db *sql.DB
}

func NewUserPostgresRepository(db *sql.DB) *UserPostgresRepository {
	return &UserPostgresRepository{db: db}
}

func (r *UserPostgresRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (first_name, last_name, phone, password_hash, priority_category)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return r.db.QueryRow(query, user.FirstName, user.LastName, user.Phone, user.PasswordHash, user.PriorityCategory).
		Scan(&user.ID, &user.CreatedAt)
}

func (r *UserPostgresRepository) GetByPhone(phone string) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, phone, password_hash, priority_category, created_at
		FROM users
		WHERE phone = $1
	`

	var user models.User
	err := r.db.QueryRow(query, phone).
		Scan(&user.ID, &user.FirstName, &user.LastName, &user.Phone, &user.PasswordHash, &user.PriorityCategory, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserPostgresRepository) GetByID(id int) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, phone, password_hash, priority_category, created_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(query, id).
		Scan(&user.ID, &user.FirstName, &user.LastName, &user.Phone, &user.PasswordHash, &user.PriorityCategory, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
