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

	if err := json.NewDecoder(r.Body).Decode(&sp); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.Create(sp)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func (h *ServicePointHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	points, err := h.service.GetAll()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, points)
}
