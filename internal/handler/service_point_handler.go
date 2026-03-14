package handler

import (
	"encoding/json"
	"net/http"

	"smartqueue/internal/models"
	"smartqueue/internal/service"
)

type ServicePointHandler struct {
	service *service.ServicePointService
}

func NewServicePointHandler(s *service.ServicePointService) *ServicePointHandler {
	return &ServicePointHandler{service: s}
}

func (h *ServicePointHandler) Create(w http.ResponseWriter, r *http.Request) {

	var sp models.ServicePoint

	err := json.NewDecoder(r.Body).Decode(&sp)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	result := h.service.Create(sp)

	json.NewEncoder(w).Encode(result)
}

func (h *ServicePointHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	points := h.service.GetAll()

	json.NewEncoder(w).Encode(points)
}