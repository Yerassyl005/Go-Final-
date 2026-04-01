package repository

import (
	"database/sql"
	"smartqueue/internal/models"
)

type ServicePointPostgresRepository struct {
	db *sql.DB
}

func NewServicePointPostgresRepository(db *sql.DB) *ServicePointPostgresRepository {
	return &ServicePointPostgresRepository{db: db}
}

func (r *ServicePointPostgresRepository) Create(sp models.ServicePoint) (models.ServicePoint, error) {
	query := `INSERT INTO service_points (name, description) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(query, sp.Name, sp.Description).Scan(&sp.ID)
	if err != nil {
		return sp, err
	}
	return sp, nil
}

func (r *ServicePointPostgresRepository) GetAll() ([]models.ServicePoint, error) {
	rows, err := r.db.Query(`SELECT id, name, description FROM service_points ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []models.ServicePoint
	for rows.Next() {
		var sp models.ServicePoint
		if err := rows.Scan(&sp.ID, &sp.Name, &sp.Description); err != nil {
			return nil, err
		}
		points = append(points, sp)
	}

	return points, rows.Err()
}
