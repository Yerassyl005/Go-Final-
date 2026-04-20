package service

import "smartqueue/internal/models"

type ServicePointRepository interface {
	Create(models.ServicePoint) (models.ServicePoint, error)
	GetAll() ([]models.ServicePoint, error)
}

type QueueRepository interface {
	Create(models.Queue) (models.Queue, error)
	GetAll() ([]models.Queue, error)
	GetByServicePoint(int) ([]models.Queue, error)
	GetDisplay(int) (models.QueueDisplay, error)
	GetStats(int) (models.QueueStats, error)
}

type TicketRepository interface {
	Create(int, int) (models.Ticket, error)
	GetAll() ([]models.Ticket, error)
	GetByID(int) (*models.Ticket, error)
	GetCurrent(int) (*models.Ticket, error)
	CallNext(int) (*models.Ticket, error)
	CallSkipped(int) (*models.Ticket, error)
	RecallCurrent(int) (*models.Ticket, error)
	SkipCurrent(int) (*models.Ticket, error)
	CompleteCurrent(int) (*models.Ticket, error)
	GetPosition(int) (int, error)
}
