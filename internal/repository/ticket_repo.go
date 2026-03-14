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

func (r *TicketRepository) Create(queueID int) models.Ticket {

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