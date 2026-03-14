package service

import "smartqueue/internal/models"

type QueueService struct {
	repo QueueRepository
}

func NewQueueService(r QueueRepository) *QueueService {
	return &QueueService{repo: r}
}

func (s *QueueService) Create(q models.Queue) models.Queue {
	return s.repo.Create(q)
}

func (s *QueueService) GetAll() []models.Queue {
	return s.repo.GetAll()
}