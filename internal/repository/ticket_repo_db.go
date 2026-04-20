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

type ticketScanner interface {
	Scan(dest ...any) error
}

func scanTicket(row ticketScanner, ticket *models.Ticket) error {
	return row.Scan(
		&ticket.ID,
		&ticket.QueueID,
		&ticket.Number,
		&ticket.Status,
		&ticket.RecallCount,
	)
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
	ticket.Status = models.TicketStatusWaiting

	return ticket, nil
}

func (r *TicketPostgresRepository) GetAll() ([]models.Ticket, error) {
	rows, err := r.db.Query(`
		SELECT id, queue_id, number, status, recall_count
		FROM tickets
		ORDER BY queue_id, number
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []models.Ticket
	for rows.Next() {
		var t models.Ticket
		if err := rows.Scan(&t.ID, &t.QueueID, &t.Number, &t.Status, &t.RecallCount); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}

	return tickets, rows.Err()
}

func (r *TicketPostgresRepository) GetByID(ticketID int) (*models.Ticket, error) {
	var ticket models.Ticket

	err := scanTicket(r.db.QueryRow(`
		SELECT id, queue_id, number, status, recall_count
		FROM tickets
		WHERE id = $1
	`, ticketID), &ticket)
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *TicketPostgresRepository) GetCurrent(queueID int) (*models.Ticket, error) {
	var ticket models.Ticket

	row := r.db.QueryRow(`
		SELECT id, queue_id, number, status, recall_count
		FROM tickets
		WHERE queue_id = $1 AND status = $2
		ORDER BY called_at DESC NULLS LAST, id DESC
		LIMIT 1
	`, queueID, models.TicketStatusCalled)

	err := scanTicket(row, &ticket)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *TicketPostgresRepository) CallSkipped(ticketID int) (*models.Ticket, error) {
	var ticket models.Ticket

	err := scanTicket(r.db.QueryRow(`
		UPDATE tickets
		SET status = $2, called_at = NOW()
		WHERE id = $1 AND status = $3
		RETURNING id, queue_id, number, status, recall_count
	`, ticketID, models.TicketStatusCalled, models.TicketStatusSkipped), &ticket)
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *TicketPostgresRepository) CallNext(queueID int) (*models.Ticket, error) {
	var ticket models.Ticket

	err := scanTicket(r.db.QueryRow(`
		UPDATE tickets
		SET status = $2, called_at = NOW()
		WHERE id = (
			SELECT id
			FROM tickets
			WHERE queue_id = $1 AND status = $3
			ORDER BY number ASC
			LIMIT 1
		)
		RETURNING id, queue_id, number, status, recall_count
	`, queueID, models.TicketStatusCalled, models.TicketStatusWaiting), &ticket)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *TicketPostgresRepository) RecallCurrent(queueID int) (*models.Ticket, error) {
	var ticket models.Ticket

	err := scanTicket(r.db.QueryRow(`
		UPDATE tickets
		SET recall_count = recall_count + 1, called_at = NOW()
		WHERE id = (
			SELECT id
			FROM tickets
			WHERE queue_id = $1 AND status = $2
			ORDER BY called_at DESC NULLS LAST, id DESC
			LIMIT 1
		)
		RETURNING id, queue_id, number, status, recall_count
	`, queueID, models.TicketStatusCalled), &ticket)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *TicketPostgresRepository) SkipCurrent(queueID int) (*models.Ticket, error) {
	var ticket models.Ticket

	err := scanTicket(r.db.QueryRow(`
		UPDATE tickets
		SET status = $2
		WHERE id = (
			SELECT id
			FROM tickets
			WHERE queue_id = $1 AND status = $3
			ORDER BY called_at DESC NULLS LAST, id DESC
			LIMIT 1
		)
		RETURNING id, queue_id, number, status, recall_count
	`, queueID, models.TicketStatusSkipped, models.TicketStatusCalled), &ticket)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *TicketPostgresRepository) CompleteCurrent(queueID int) (*models.Ticket, error) {
	var ticket models.Ticket

	err := scanTicket(r.db.QueryRow(`
		UPDATE tickets
		SET status = $2, completed_at = NOW()
		WHERE id = (
			SELECT id
			FROM tickets
			WHERE queue_id = $1 AND status = $3
			ORDER BY called_at DESC NULLS LAST, id DESC
			LIMIT 1
		)
		RETURNING id, queue_id, number, status, recall_count
	`, queueID, models.TicketStatusCompleted, models.TicketStatusCalled), &ticket)

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

	if status != models.TicketStatusWaiting && status != models.TicketStatusCalled {
		return 0, nil
	}

	var ahead int
	err = r.db.QueryRow(`
		SELECT COUNT(*)
		FROM tickets
		WHERE queue_id = $1
		AND number < $2
		AND status IN ($3, $4)
	`, queueID, number, models.TicketStatusWaiting, models.TicketStatusCalled).Scan(&ahead)
	if err != nil {
		return -1, err
	}

	return ahead + 1, nil
}
