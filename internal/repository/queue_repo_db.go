package repository

import (
	"database/sql"
	"smartqueue/internal/models"
)

type QueuePostgresRepository struct {
	db *sql.DB
}

func NewQueuePostgresRepository(db *sql.DB) *QueuePostgresRepository {
	return &QueuePostgresRepository{db: db}
}

func (r *QueuePostgresRepository) GetByServicePoint(servicePointID int) []models.Queue {

	rows, err := r.db.Query(`
		SELECT id, name, service_point_id 
		FROM queues 
		WHERE service_point_id = $1
	`, servicePointID)

	if err != nil {
		return nil
	}
	defer rows.Close()

	var queues []models.Queue

	for rows.Next() {
		var q models.Queue
		rows.Scan(&q.ID, &q.Name, &q.ServicePointID)
		queues = append(queues, q)
	}

	return queues
}