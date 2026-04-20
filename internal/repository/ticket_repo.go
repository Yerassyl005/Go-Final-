package repository

import "smartqueue/internal/models"

type TicketRepository struct {
	tickets []models.Ticket
	nextID  int
}

func NewTicketRepository() *TicketRepository {
	return &TicketRepository{
		tickets: []models.Ticket{},
		nextID:  1,
	}
}

func (r *TicketRepository) Create(queueID int, userID int) (models.Ticket, error) {
	nextNumber := 1
	for _, ticket := range r.tickets {
		if ticket.QueueID == queueID && ticket.Number >= nextNumber {
			nextNumber = ticket.Number + 1
		}
	}

	ticket := models.Ticket{
		ID:      r.nextID,
		QueueID: queueID,
		Number:  nextNumber,
		Status:  models.TicketStatusWaiting,
	}

	r.nextID++
	r.tickets = append(r.tickets, ticket)

	return ticket, nil
}

func (r *TicketRepository) GetAll() ([]models.Ticket, error) {
	return r.tickets, nil
}

func (r *TicketRepository) GetByID(id int) (*models.Ticket, error) {
	for i := range r.tickets {
		if r.tickets[i].ID == id {
			return &r.tickets[i], nil
		}
	}

	return nil, nil
}

func (r *TicketRepository) GetCurrent(queueID int) (*models.Ticket, error) {
	for i := range r.tickets {
		if r.tickets[i].QueueID == queueID && r.tickets[i].Status == models.TicketStatusCalled {
			return &r.tickets[i], nil
		}
	}

	return nil, nil
}

func (r *TicketRepository) CallSkipped(ticketID int) (*models.Ticket, error) {
	ticket, err := r.GetByID(ticketID)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, nil
	}

	ticket.Status = models.TicketStatusCalled
	return ticket, nil
}

func (r *TicketRepository) CallNext(queueID int) (*models.Ticket, error) {
	current, err := r.GetCurrent(queueID)
	if err != nil {
		return nil, err
	}
	if current != nil {
		return nil, nil
	}

	for i := range r.tickets {
		if r.tickets[i].QueueID == queueID && r.tickets[i].Status == models.TicketStatusWaiting {
			r.tickets[i].Status = models.TicketStatusCalled
			return &r.tickets[i], nil
		}
	}

	return nil, nil
}

func (r *TicketRepository) RecallCurrent(queueID int) (*models.Ticket, error) {
	ticket, err := r.GetCurrent(queueID)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, nil
	}

	ticket.RecallCount++
	return ticket, nil
}

func (r *TicketRepository) SkipCurrent(queueID int) (*models.Ticket, error) {
	ticket, err := r.GetCurrent(queueID)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, nil
	}

	ticket.Status = models.TicketStatusSkipped
	return ticket, nil
}

func (r *TicketRepository) CompleteCurrent(queueID int) (*models.Ticket, error) {
	ticket, err := r.GetCurrent(queueID)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, nil
	}

	ticket.Status = models.TicketStatusCompleted
	return ticket, nil
}

func (r *TicketRepository) GetPosition(id int) (int, error) {
	position := 0
	queueID := 0

	for _, ticket := range r.tickets {
		if ticket.ID == id {
			queueID = ticket.QueueID
			break
		}
	}

	for _, ticket := range r.tickets {
		if ticket.QueueID != queueID {
			continue
		}
		if ticket.Status == models.TicketStatusWaiting || ticket.Status == models.TicketStatusCalled {
			position++
		}
		if ticket.ID == id {
			if ticket.Status == models.TicketStatusWaiting || ticket.Status == models.TicketStatusCalled {
				return position, nil
			}
			return 0, nil
		}
	}

	return 0, nil
}
