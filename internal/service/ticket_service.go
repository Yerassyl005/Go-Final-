package service

import (
	"database/sql"
	"errors"
	"smartqueue/internal/models"
)

var (
	ErrInvalidQueueID    = errors.New("invalid queue id")
	ErrInvalidTicketID   = errors.New("invalid ticket id")
	ErrTicketNotFound    = errors.New("ticket not found")
	ErrTicketNotSkipped  = errors.New("ticket is not in skipped status")
	ErrNoWaitingTickets  = errors.New("no waiting tickets in this queue")
	ErrNoActiveTicket    = errors.New("no active ticket in this queue")
	ErrActiveTicketExist = errors.New("another ticket is already active in this queue")

	ErrQueueClosed     = errors.New("queue is closed")
	ErrTicketNotCalled = errors.New("ticket is not called")
)

type TicketService struct {
	repo      TicketRepository
	userRepo  UserReader
	queueRepo QueueReader // 👈 ВАЖНО: QueueReader, не QueueRepository
}

func NewTicketService(r TicketRepository, userRepo UserReader, queueRepo QueueReader) *TicketService {
	return &TicketService{
		repo:      r,
		userRepo:  userRepo,
		queueRepo: queueRepo,
	}
}

func (s *TicketService) Create(queueID int, userID int) (models.Ticket, error) {
	if queueID <= 0 {
		return models.Ticket{}, ErrInvalidQueueID
	}

	// 🔒 проверка очереди
	queue, err := s.queueRepo.GetByID(queueID)
	if err != nil {
		return models.Ticket{}, err
	}
	if queue == nil {
		return models.Ticket{}, ErrInvalidQueueID
	}
	if !queue.IsOpen {
		return models.Ticket{}, ErrQueueClosed
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return models.Ticket{}, err
	}

	isPriority := user.PriorityCategory != models.PriorityCategoryNone
	return s.repo.Create(queueID, userID, isPriority)
}

func (s *TicketService) GetAll() ([]models.Ticket, error) {
	return s.repo.GetAll()
}

func (s *TicketService) CallNext(queueID int) (*models.Ticket, error) {
	if queueID <= 0 {
		return nil, ErrInvalidQueueID
	}

	current, err := s.repo.GetCurrent(queueID)
	if err != nil {
		return nil, err
	}
	if current != nil {
		return nil, ErrActiveTicketExist
	}

	ticket, err := s.repo.CallNext(queueID)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, ErrNoWaitingTickets
	}

	return ticket, nil
}

func (s *TicketService) RecallCurrent(queueID int) (*models.Ticket, error) {
	if queueID <= 0 {
		return nil, ErrInvalidQueueID
	}

	ticket, err := s.repo.RecallCurrent(queueID)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, ErrNoActiveTicket
	}

	return ticket, nil
}

func (s *TicketService) CallSkipped(ticketID int) (*models.Ticket, error) {
	if ticketID <= 0 {
		return nil, ErrInvalidTicketID
	}

	ticket, err := s.repo.GetByID(ticketID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTicketNotFound
	}
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, ErrTicketNotFound
	}
	if ticket.Status != models.TicketStatusSkipped {
		return nil, ErrTicketNotSkipped
	}

	current, err := s.repo.GetCurrent(ticket.QueueID)
	if err != nil {
		return nil, err
	}
	if current != nil {
		return nil, ErrActiveTicketExist
	}

	ticket, err = s.repo.CallSkipped(ticketID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTicketNotFound
	}
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, ErrTicketNotFound
	}

	return ticket, nil
}

func (s *TicketService) SkipCurrent(queueID int) (*models.Ticket, error) {
	if queueID <= 0 {
		return nil, ErrInvalidQueueID
	}

	ticket, err := s.repo.SkipCurrent(queueID)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, ErrNoActiveTicket
	}

	return ticket, nil
}

func (s *TicketService) CompleteCurrent(queueID int) (*models.Ticket, error) {
	if queueID <= 0 {
		return nil, ErrInvalidQueueID
	}

	current, err := s.repo.GetCurrent(queueID)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, ErrNoActiveTicket
	}
	if current.Status != models.TicketStatusCalled {
		return nil, ErrTicketNotCalled
	}

	ticket, err := s.repo.CompleteCurrent(queueID)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, ErrNoActiveTicket
	}

	return ticket, nil
}

func (s *TicketService) GetPosition(id int) (int, error) {
	if id <= 0 {
		return 0, ErrInvalidTicketID
	}

	position, err := s.repo.GetPosition(id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrTicketNotFound
	}
	if err != nil {
		return 0, err
	}

	return position, nil
}
