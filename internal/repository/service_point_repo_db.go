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

func (r *ServicePointPostgresRepository) Create(sp models.ServicePoint) models.ServicePoint {

	query := `INSERT INTO service_points (name) VALUES ($1) RETURNING id`

	err := r.db.QueryRow(query, sp.Name).Scan(&sp.ID)
	if err != nil {
		return sp
	}

	return sp
}

func (r *ServicePointPostgresRepository) GetAll() []models.ServicePoint {

	rows, err := r.db.Query(`SELECT id, name FROM service_points`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var points []models.ServicePoint

	for rows.Next() {
		var sp models.ServicePoint
		rows.Scan(&sp.ID, &sp.Name)
		points = append(points, sp)
	}

	return points
}