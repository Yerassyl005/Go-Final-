package repository

import "smartqueue/internal/models"

type QueueRepository struct {
	queues []models.Queue
	nextID int
}

func NewQueueRepository() *QueueRepository {
	return &QueueRepository{
		queues: []models.Queue{},
		nextID: 1,
	}
}

func (r *QueueRepository) Create(q models.Queue) models.Queue {
	q.ID = r.nextID
	r.nextID++

	r.queues = append(r.queues, q)

	return q
}

func (r *QueueRepository) GetAll() []models.Queue {
	return r.queues
}

func (r *QueueRepository) GetByServicePoint(servicePointID int) []models.Queue {
	var result []models.Queue

	for _, q := range r.queues {
		if q.ServicePointID == servicePointID {
			result = append(result, q)
		}
	}

	return result
}