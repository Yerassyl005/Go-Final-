package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"smartqueue/internal/models"
	"smartqueue/internal/service"

	"github.com/gorilla/mux"
)

type QueueHandler struct {
	service *service.QueueService
}

func NewQueueHandler(s *service.QueueService) *QueueHandler {
	return &QueueHandler{service: s}
}

func (h *QueueHandler) Create(w http.ResponseWriter, r *http.Request) {
	var q models.Queue

	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.Create(q)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func (h *QueueHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	queues, err := h.service.GetAll()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, queues)
}

func (h *QueueHandler) GetByServicePoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid service point id")
		return
	}

	queues, err := h.service.GetByServicePoint(id)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, queues)
}

func (h *QueueHandler) GetDisplay(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid queue id")
		return
	}

	display, err := h.service.GetDisplay(id)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, display)
}

func (h *QueueHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid queue id")
		return
	}

	stats, err := h.service.GetStats(id)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, stats)
}
