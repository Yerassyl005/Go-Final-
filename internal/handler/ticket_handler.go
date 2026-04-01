package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"smartqueue/internal/middleware"
	"smartqueue/internal/service"

	"github.com/gorilla/mux"
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json")
		return
	}

	claims, ok := r.Context().Value(middleware.UserContextKey).(*service.AuthClaims)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ticket, err := h.service.Create(req.QueueID, claims.UserID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, ticket)
}

func (h *TicketHandler) GetTickets(w http.ResponseWriter, r *http.Request) {
	tickets, err := h.service.GetAll()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, tickets)
}

func (h *TicketHandler) CallNext(w http.ResponseWriter, r *http.Request) {
	ticket, err := h.service.CallNext()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if ticket == nil {
		writeJSONError(w, http.StatusNotFound, "no waiting tickets")
		return
	}

	writeJSON(w, http.StatusOK, ticket)
}

func (h *TicketHandler) CompleteTicket(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid ticket id")
		return
	}

	ticket, err := h.service.Complete(id)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if ticket == nil {
		writeJSONError(w, http.StatusNotFound, "ticket not found")
		return
	}

	writeJSON(w, http.StatusOK, ticket)
}

func (h *TicketHandler) GetPosition(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid ticket id")
		return
	}

	position, err := h.service.GetPosition(id)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"ticket_id": id,
		"position":  position,
	})
}

func (h *TicketHandler) SkipTicket(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid ticket id")
		return
	}

	ticket, err := h.service.Skip(id)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if ticket == nil {
		writeJSONError(w, http.StatusNotFound, "ticket not found")
		return
	}

	writeJSON(w, http.StatusOK, ticket)
}
