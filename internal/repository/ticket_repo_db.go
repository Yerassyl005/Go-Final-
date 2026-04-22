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
		&ticket.IsPriority,
	)
}

func (r *TicketPostgresRepository) Create(queueID int, userID int, isPriority bool) (models.Ticket, error) {
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
		INSERT INTO tickets (queue_id, user_id, number, status, is_priority)
		VALUES ($1, $2, $3, 'waiting', $4)
		RETURNING id
	`, queueID, userID, number, isPriority).Scan(&ticket.ID)
	if err != nil {
		return models.Ticket{}, err
	}

	ticket.QueueID = queueID
	ticket.Number = number
	ticket.Status = models.TicketStatusWaiting
	ticket.IsPriority = isPriority

	return ticket, nil
}

func (r *TicketPostgresRepository) GetAll() ([]models.Ticket, error) {
	rows, err := r.db.Query(`
		SELECT id, queue_id, number, status, recall_count, is_priority
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
		if err := rows.Scan(&t.ID, &t.QueueID, &t.Number, &t.Status, &t.RecallCount, &t.IsPriority); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}

	return tickets, rows.Err()
}

func (r *TicketPostgresRepository) GetByID(ticketID int) (*models.Ticket, error) {
	var ticket models.Ticket

	err := scanTicket(r.db.QueryRow(`
		SELECT id, queue_id, number, status, recall_count, is_priority
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
		SELECT id, queue_id, number, status, recall_count, is_priority
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
		RETURNING id, queue_id, number, status, recall_count, is_priority
	`, ticketID, models.TicketStatusCalled, models.TicketStatusSkipped), &ticket)
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *TicketPostgresRepository) CallNext(queueID int) (*models.Ticket, error) {
	var ticket models.Ticket

	// Priority tickets (is_priority DESC) are always served before regular ones.
	// Within the same priority tier, FIFO ordering applies (number ASC).
	err := scanTicket(r.db.QueryRow(`
		UPDATE tickets
		SET status = $2, called_at = NOW()
		WHERE id = (
			SELECT id
			FROM tickets
			WHERE queue_id = $1 AND status = $3
			ORDER BY is_priority DESC, number ASC
			LIMIT 1
		)
		RETURNING id, queue_id, number, status, recall_count, is_priority
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
		RETURNING id, queue_id, number, status, recall_count, is_priority
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
		RETURNING id, queue_id, number, status, recall_count, is_priority
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
		RETURNING id, queue_id, number, status, recall_count, is_priority
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
	var isPriority bool

	err := r.db.QueryRow(`
		SELECT queue_id, number, status, is_priority
		FROM tickets
		WHERE id = $1
	`, ticketID).Scan(&queueID, &number, &status, &isPriority)
	if err != nil {
		return -1, err
	}

	if status != models.TicketStatusWaiting && status != models.TicketStatusCalled {
		return 0, nil
	}

	var ahead int
	if isPriority {
		// For priority ticket: only other priority tickets with a lower number are ahead,
		// plus any currently called ticket.
		err = r.db.QueryRow(`
			SELECT COUNT(*)
			FROM tickets
			WHERE queue_id = $1
			AND (
				status = $2
				OR (is_priority = true AND status = $3 AND number < $4)
			)
		`, queueID, models.TicketStatusCalled, models.TicketStatusWaiting, number).Scan(&ahead)
	} else {
		// For regular ticket: all priority waiting tickets are ahead,
		// plus regular tickets with a lower number, plus any currently called ticket.
		err = r.db.QueryRow(`
			SELECT COUNT(*)
			FROM tickets
			WHERE queue_id = $1
			AND (
				status = $2
				OR (is_priority = true AND status = $3)
				OR (is_priority = false AND status = $3 AND number < $4)
			)
		`, queueID, models.TicketStatusCalled, models.TicketStatusWaiting, number).Scan(&ahead)
	}
	if err != nil {
		return -1, err
	}

	return ahead + 1, nil
}
