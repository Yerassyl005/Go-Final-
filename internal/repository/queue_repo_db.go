package repository

import (
	"database/sql"
	"fmt"
	"smartqueue/internal/models"
)

type QueuePostgresRepository struct {
	db *sql.DB
}

func NewQueuePostgresRepository(db *sql.DB) *QueuePostgresRepository {
	return &QueuePostgresRepository{db: db}
}

func formatTicketNumber(number int) string {
	return fmt.Sprintf("A-%03d", number)
}

func (r *QueuePostgresRepository) Create(q models.Queue) (models.Queue, error) {
	query := `INSERT INTO queues (name, service_point_id) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(query, q.Name, q.ServicePointID).Scan(&q.ID)
	if err != nil {
		return q, err
	}
	return q, nil
}

func (r *QueuePostgresRepository) GetAll() ([]models.Queue, error) {
	rows, err := r.db.Query(`SELECT id, name, service_point_id, is_open FROM queues ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var queues []models.Queue
	for rows.Next() {
		var q models.Queue
		if err := rows.Scan(&q.ID, &q.Name, &q.ServicePointID, &q.IsOpen); err != nil {
			return nil, err
		}
		queues = append(queues, q)
	}

	return queues, rows.Err()
}

func (r *QueuePostgresRepository) GetByServicePoint(servicePointID int) ([]models.Queue, error) {
	rows, err := r.db.Query(`
		SELECT id, name, service_point_id, is_open
		FROM queues
		WHERE service_point_id = $1
		ORDER BY id
	`, servicePointID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var queues []models.Queue
	for rows.Next() {
		var q models.Queue
		if err := rows.Scan(&q.ID, &q.Name, &q.ServicePointID, &q.IsOpen); err != nil {
			return nil, err
		}
		queues = append(queues, q)
	}

	return queues, rows.Err()
}

func (r *QueuePostgresRepository) GetDisplay(queueID int) (models.QueueDisplay, error) {
	display := models.QueueDisplay{
		QueueID:        queueID,
		WaitingTickets: []models.QueueDisplayTicket{},
	}

	if err := r.db.QueryRow(`SELECT id, name FROM queues WHERE id = $1`, queueID).
		Scan(&display.QueueID, &display.QueueName); err != nil {
		return display, err
	}

	var currentID, currentNumber int
	var currentStatus string
	var currentRecallCount int
	var currentIsPriority bool

	err := r.db.QueryRow(`
		SELECT id, number, status, recall_count, is_priority
		FROM tickets
		WHERE queue_id = $1 AND status = $2
		ORDER BY called_at DESC NULLS LAST, id DESC
		LIMIT 1
	`, queueID, models.TicketStatusCalled).Scan(&currentID, &currentNumber, &currentStatus, &currentRecallCount, &currentIsPriority)

	if err != nil && err != sql.ErrNoRows {
		return display, err
	}

	if err == nil {
		display.CurrentTicket = &models.QueueDisplayTicket{
			ID:           currentID,
			TicketNumber: formatTicketNumber(currentNumber),
			Status:       currentStatus,
			RecallCount:  currentRecallCount,
			IsPriority:   currentIsPriority,
		}
	}

	rows, err := r.db.Query(`
		SELECT id, number, status, is_priority
		FROM tickets
		WHERE queue_id = $1 AND status = 'waiting'
		ORDER BY is_priority DESC, number ASC
	`, queueID)
	if err != nil {
		return display, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, number int
		var status string
		var isPriority bool

		if err := rows.Scan(&id, &number, &status, &isPriority); err != nil {
			return display, err
		}

		display.WaitingTickets = append(display.WaitingTickets, models.QueueDisplayTicket{
			ID:           id,
			TicketNumber: formatTicketNumber(number),
			Status:       status,
			IsPriority:   isPriority,
		})
	}

	if err := rows.Err(); err != nil {
		return display, err
	}

	if err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM tickets
		WHERE queue_id = $1 AND status = 'completed'
	`, queueID).Scan(&display.CompletedCount); err != nil {
		return display, err
	}

	return display, nil
}

func (r *QueuePostgresRepository) GetStats(queueID int) (models.QueueStats, error) {
	stats := models.QueueStats{QueueID: queueID}

	query := `
		SELECT
			COUNT(*) AS total_tickets,
			COUNT(*) FILTER (WHERE status = $2) AS waiting_tickets,
			COUNT(*) FILTER (WHERE status = $3) AS called_tickets,
			COUNT(*) FILTER (WHERE status = $4) AS completed_tickets,
			COUNT(*) FILTER (WHERE status = $5) AS skipped_tickets
		FROM tickets
		WHERE queue_id = $1
	`

	err := r.db.QueryRow(
		query,
		queueID,
		models.TicketStatusWaiting,
		models.TicketStatusCalled,
		models.TicketStatusCompleted,
		models.TicketStatusSkipped,
	).Scan(
		&stats.TotalTickets,
		&stats.WaitingTickets,
		&stats.CalledTickets,
		&stats.CompletedTickets,
		&stats.SkippedTickets,
	)
	if err != nil {
		return stats, err
	}

	return stats, nil
}

func (r *QueuePostgresRepository) GetByID(id int) (*models.Queue, error) {
	var q models.Queue

	err := r.db.QueryRow(`
		SELECT id, name, service_point_id, is_open
		FROM queues
		WHERE id = $1
	`, id).Scan(&q.ID, &q.Name, &q.ServicePointID, &q.IsOpen)

	if err != nil {
		return nil, err
	}

	return &q, nil
}
