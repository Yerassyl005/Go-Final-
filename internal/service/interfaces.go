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
	CallNext() (*models.Ticket, error)
	Complete(int) (*models.Ticket, error)
	GetPosition(int) (int, error)
	Skip(int) (*models.Ticket, error)
}
