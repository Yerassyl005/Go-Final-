package handler

import (
	"net/http"

	"smartqueue/internal/service"

	"github.com/gorilla/mux"
)

type QRHandler struct {
	service *service.QRService
}

func NewQRHandler(s *service.QRService) *QRHandler {
	return &QRHandler{service: s}
}

// Старый метод для произвольной строки
func (h *QRHandler) GenerateQR(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query().Get("data")

	if data == "" {
		http.Error(w, "data is required", http.StatusBadRequest)
		return
	}

	qr, err := h.service.GenerateQR(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(qr)
}

// Новый метод для конкретного тикета
func (h *QRHandler) GetTicketQR(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	data := "http://localhost:8080/tickets/" + id

	qr, err := h.service.GenerateQR(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(qr)
}
