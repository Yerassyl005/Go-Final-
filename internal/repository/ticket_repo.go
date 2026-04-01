package repository

import "smartqueue/internal/models"

type TicketRepository struct {
	tickets []models.Ticket
	nextID  int
	nextNum int
}

func NewTicketRepository() *TicketRepository {
	return &TicketRepository{
		tickets: []models.Ticket{},
		nextID:  1,
		nextNum: 1,
	}
}

func (r *TicketRepository) Create(queueID int, userID int) models.Ticket {
	ticket := models.Ticket{
		ID:      r.nextID,
		QueueID: queueID,
		Number:  r.nextNum,
		Status:  "waiting",
	}

	r.nextID++
	r.nextNum++

	r.tickets = append(r.tickets, ticket)

	return ticket
}

func (r *TicketRepository) GetAll() []models.Ticket {
	return r.tickets
}

func (r *TicketRepository) CallNext() *models.Ticket {
	for i := range r.tickets {
		if r.tickets[i].Status == "waiting" {
			r.tickets[i].Status = "called"
			return &r.tickets[i]
		}
	}
	return nil
}

func (r *TicketRepository) Complete(id int) *models.Ticket {
	for i := range r.tickets {
		if r.tickets[i].ID == id {
			r.tickets[i].Status = "completed"
			return &r.tickets[i]
		}
	}
	return nil
}

func (r *TicketRepository) GetPosition(id int) int {
	position := 0

	for _, ticket := range r.tickets {
		if ticket.Status == "waiting" {
			position++
		}
		if ticket.ID == id {
			if ticket.Status == "waiting" {
				return position
			}
			return 0
		}
	}

	return 0
}
func (r *TicketRepository) Skip(id int) *models.Ticket {
	for i := range r.tickets {
		if r.tickets[i].ID == id {
			r.tickets[i].Status = "skipped"
			return &r.tickets[i]
		}
	}
	return nil
}