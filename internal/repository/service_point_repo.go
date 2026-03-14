package repository

import "smartqueue/internal/models"

type ServicePointRepository struct {
	servicePoints []models.ServicePoint
	nextID        int
}

func NewServicePointRepository() *ServicePointRepository {
	return &ServicePointRepository{
		servicePoints: []models.ServicePoint{},
		nextID:        1,
	}
}

func (r *ServicePointRepository) Create(sp models.ServicePoint) models.ServicePoint {

	sp.ID = r.nextID
	r.nextID++

	r.servicePoints = append(r.servicePoints, sp)

	return sp
}

func (r *ServicePointRepository) GetAll() []models.ServicePoint {
	return r.servicePoints
}