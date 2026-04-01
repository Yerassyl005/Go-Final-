package service

import "smartqueue/internal/models"

type QueueService struct {
	repo QueueRepository
}

func NewQueueService(r QueueRepository) *QueueService {
	return &QueueService{repo: r}
}

func (s *QueueService) Create(q models.Queue) (models.Queue, error) {
	return s.repo.Create(q)
}

func (s *QueueService) GetAll() ([]models.Queue, error) {
	return s.repo.GetAll()
}

func (s *QueueService) GetByServicePoint(servicePointID int) ([]models.Queue, error) {
	return s.repo.GetByServicePoint(servicePointID)
}

func (s *QueueService) GetDisplay(queueID int) (models.QueueDisplay, error) {
	return s.repo.GetDisplay(queueID)
}

func (s *QueueService) GetStats(queueID int) (models.QueueStats, error) {
	return s.repo.GetStats(queueID)
}
