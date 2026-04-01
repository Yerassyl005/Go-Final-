package service

import "smartqueue/internal/models"

type TicketService struct {
	repo TicketRepository
}

func NewTicketService(r TicketRepository) *TicketService {
	return &TicketService{repo: r}
}

func (s *TicketService) Create(queueID int, userID int) (models.Ticket, error) {
	return s.repo.Create(queueID, userID)
}

func (s *TicketService) GetAll() ([]models.Ticket, error) {
	return s.repo.GetAll()
}

func (s *TicketService) CallNext() (*models.Ticket, error) {
	return s.repo.CallNext()
}

func (s *TicketService) Complete(id int) (*models.Ticket, error) {
	return s.repo.Complete(id)
}

func (s *TicketService) GetPosition(id int) (int, error) {
	return s.repo.GetPosition(id)
}

func (s *TicketService) Skip(id int) (*models.Ticket, error) {
	return s.repo.Skip(id)
}
