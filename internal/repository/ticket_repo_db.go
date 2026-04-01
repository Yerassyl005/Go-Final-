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

func (r *TicketPostgresRepository) Create(queueID int, userID int) (models.Ticket, error) {
	var number int

	err := r.db.QueryRow(`
		SELECT COALESCE(MAX(number), 0) + 1
		FROM tickets
		WHERE queue_id = $1
	`, queueID).Scan(&number)
	if err != nil {
		return models.Ticket{}, err
	}

	var ticket models.Ticket
	err = r.db.QueryRow(`
		INSERT INTO tickets (queue_id, user_id, number, status)
		VALUES ($1, $2, $3, 'waiting')
		RETURNING id
	`, queueID, userID, number).Scan(&ticket.ID)
	if err != nil {
		return models.Ticket{}, err
	}

	ticket.QueueID = queueID
	ticket.Number = number
	ticket.Status = "waiting"

	return ticket, nil
}

func (r *TicketPostgresRepository) GetAll() ([]models.Ticket, error) {
	rows, err := r.db.Query(`SELECT id, queue_id, number, status FROM tickets ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []models.Ticket
	for rows.Next() {
		var t models.Ticket
		if err := rows.Scan(&t.ID, &t.QueueID, &t.Number, &t.Status); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}

	return tickets, rows.Err()
}

func (r *TicketPostgresRepository) CallNext() (*models.Ticket, error) {
	var ticket models.Ticket

	err := r.db.QueryRow(`
		UPDATE tickets
		SET status = 'called', called_at = NOW()
		WHERE id = (
			SELECT id
			FROM tickets
			WHERE status = 'waiting'
			ORDER BY created_at ASC, id ASC
			LIMIT 1
		)
		RETURNING id, queue_id, number, status
	`).Scan(&ticket.ID, &ticket.QueueID, &ticket.Number, &ticket.Status)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *TicketPostgresRepository) Complete(id int) (*models.Ticket, error) {
	var ticket models.Ticket

	err := r.db.QueryRow(`
		UPDATE tickets
		SET status = 'completed', completed_at = NOW()
		WHERE id = $1
		RETURNING id, queue_id, number, status
	`, id).Scan(&ticket.ID, &ticket.QueueID, &ticket.Number, &ticket.Status)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *TicketPostgresRepository) GetPosition(ticketID int) (int, error) {
	var queueID, number int
	var status string

	err := r.db.QueryRow(`
		SELECT queue_id, number, status
		FROM tickets
		WHERE id = $1
	`, ticketID).Scan(&queueID, &number, &status)
	if err != nil {
		return -1, err
	}

	if status != "waiting" && status != "called" && status != "in_progress" {
		return 0, nil
	}

	var ahead int
	err = r.db.QueryRow(`
		SELECT COUNT(*)
		FROM tickets
		WHERE queue_id = $1
		AND number < $2
		AND status IN ('waiting', 'called', 'in_progress')
	`, queueID, number).Scan(&ahead)
	if err != nil {
		return -1, err
	}

	return ahead + 1, nil
}

func (r *TicketPostgresRepository) Skip(id int) (*models.Ticket, error) {
	var ticket models.Ticket

	err := r.db.QueryRow(`
		UPDATE tickets
		SET status = 'skipped'
		WHERE id = $1
		RETURNING id, queue_id, number, status
	`, id).Scan(&ticket.ID, &ticket.QueueID, &ticket.Number, &ticket.Status)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}
