package service

import "smartqueue/internal/models"

type ServicePointService struct {
	repo ServicePointRepository
}

func NewServicePointService(r ServicePointRepository) *ServicePointService {
	return &ServicePointService{repo: r}
}

func (s *ServicePointService) Create(sp models.ServicePoint) (models.ServicePoint, error) {
	return s.repo.Create(sp)
}

func (s *ServicePointService) GetAll() ([]models.ServicePoint, error) {
	return s.repo.GetAll()
}
