package repository

import (
	"database/sql"
	"smartqueue/internal/models"
)

type TicketPostgresRepository struct {
	db *sql.DB
}

func NewTicketPostgresRepository(db *sql.DB) *TicketPostgresRepository {
	return &TicketPostgresRepository{db: db}
}

func (r *TicketPostgresRepository) Create(queueID int, userID int) models.Ticket {

	var number int

	err := r.db.QueryRow(`
		SELECT COALESCE(MAX(number), 0) + 1 
		FROM tickets 
		WHERE queue_id = $1
	`, queueID).Scan(&number)

	if err != nil {
		return models.Ticket{}
	}

	var ticket models.Ticket

	err = r.db.QueryRow(`
		INSERT INTO tickets (queue_id, user_id, number, status)
		VALUES ($1, $2, $3, 'waiting')
		RETURNING id
	`, queueID, userID, number).Scan(&ticket.ID)

	if err != nil {
		return ticket
	}

	ticket.QueueID = queueID
	ticket.Number = number
	ticket.Status = "waiting"

	return ticket
}

func (r *TicketPostgresRepository) GetPosition(ticketID int) int {

	var queueID, number int

	err := r.db.QueryRow(`
		SELECT queue_id, number 
		FROM tickets 
		WHERE id = $1
	`, ticketID).Scan(&queueID, &number)

	if err != nil {
		return -1
	}

	var position int

	err = r.db.QueryRow(`
		SELECT COUNT(*) 
		FROM tickets
		WHERE queue_id = $1
		AND number < $2
		AND status IN ('waiting', 'called', 'in_progress')
	`, queueID, number).Scan(&position)

	if err != nil {
		return -1
	}

	return position
}