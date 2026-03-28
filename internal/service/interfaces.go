package service

import "smartqueue/internal/models"

type ServicePointRepository interface {
	Create(models.ServicePoint) models.ServicePoint
	GetAll() []models.ServicePoint
}

type QueueRepository interface {
	Create(models.Queue) models.Queue
	GetAll() []models.Queue
	GetByServicePoint(int) []models.Queue
}

type TicketRepository interface {
	Create(int, int) models.Ticket
	GetAll() []models.Ticket
	CallNext() *models.Ticket
	Complete(int) *models.Ticket
	GetPosition(int) int
	Skip(int) *models.Ticket
}
