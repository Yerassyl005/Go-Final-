package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"smartqueue/internal/service"
)

type TicketHandler struct {
	service *service.TicketService
}

func NewTicketHandler(s *service.TicketService) *TicketHandler {
	return &TicketHandler{service: s}
}

func (h *TicketHandler) TakeTicket(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		QueueID int `json:"queue_id"`
	}

	var req Request

	json.NewDecoder(r.Body).Decode(&req)

	ticket := h.service.Create(req.QueueID)

	json.NewEncoder(w).Encode(ticket)
}

func (h *TicketHandler) GetTickets(w http.ResponseWriter, r *http.Request) {

	tickets := h.service.GetAll()

	json.NewEncoder(w).Encode(tickets)
}

func (h *TicketHandler) CallNext(w http.ResponseWriter, r *http.Request) {

	ticket := h.service.CallNext()

	json.NewEncoder(w).Encode(ticket)
}

func (h *TicketHandler) CompleteTicket(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	ticket := h.service.Complete(id)

	json.NewEncoder(w).Encode(ticket)
}