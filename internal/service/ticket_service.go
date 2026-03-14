package service

import "smartqueue/internal/models"

type TicketService struct {
	repo TicketRepository
}

func NewTicketService(r TicketRepository) *TicketService {
	return &TicketService{repo: r}
}

func (s *TicketService) Create(queueID int) models.Ticket {
	return s.repo.Create(queueID)
}

func (s *TicketService) GetAll() []models.Ticket {
	return s.repo.GetAll()
}

func (s *TicketService) CallNext() *models.Ticket {
	return s.repo.CallNext()
}

func (s *TicketService) Complete(id int) *models.Ticket {
	return s.repo.Complete(id)
}