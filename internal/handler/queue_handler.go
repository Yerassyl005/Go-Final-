package handler

import (
	"encoding/json"
	"net/http"

	"smartqueue/internal/models"
	"smartqueue/internal/service"
)

type QueueHandler struct {
	service *service.QueueService
}

func NewQueueHandler(s *service.QueueService) *QueueHandler {
	return &QueueHandler{service: s}
}

func (h *QueueHandler) Create(w http.ResponseWriter, r *http.Request) {

	var q models.Queue

	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	result := h.service.Create(q)

	json.NewEncoder(w).Encode(result)
}

func (h *QueueHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	queues := h.service.GetAll()

	json.NewEncoder(w).Encode(queues)
}