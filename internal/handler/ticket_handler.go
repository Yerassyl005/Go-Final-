package handler

import (
	"encoding/json"
	"errors"
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
		if errors.Is(err, service.ErrInvalidQueueID) {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		if isForeignKeyViolation(err) && hasConstraint(err, "tickets_queue_id_fkey") {
			writeJSONError(w, http.StatusNotFound, "queue not found")
			return
		}
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
	queueID, err := getIntParam(r, "id")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid queue id")
		return
	}

	ticket, err := h.service.CallNext(queueID)
	if err != nil {
		writeTicketActionError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "ticket called",
		"ticket":  ticket,
	})
}

func (h *TicketHandler) RecallCurrent(w http.ResponseWriter, r *http.Request) {
	queueID, err := getIntParam(r, "id")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid queue id")
		return
	}

	ticket, err := h.service.RecallCurrent(queueID)
	if err != nil {
		writeTicketActionError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "ticket recalled",
		"ticket":  ticket,
	})
}

func (h *TicketHandler) CallSkipped(w http.ResponseWriter, r *http.Request) {
	ticketID, err := getIntParam(r, "id")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid ticket id")
		return
	}

	ticket, err := h.service.CallSkipped(ticketID)
	if err != nil {
		writeTicketActionError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "skipped ticket called again",
		"ticket":  ticket,
	})
}

func (h *TicketHandler) GetPosition(w http.ResponseWriter, r *http.Request) {
	id, err := getIntParam(r, "id")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid ticket id")
		return
	}

	position, err := h.service.GetPosition(id)
	if err != nil {
		if errors.Is(err, service.ErrTicketNotFound) {
			writeJSONError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"ticket_id": id,
		"position":  position,
	})
}

func (h *TicketHandler) SkipCurrent(w http.ResponseWriter, r *http.Request) {
	queueID, err := getIntParam(r, "id")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid queue id")
		return
	}

	ticket, err := h.service.SkipCurrent(queueID)
	if err != nil {
		writeTicketActionError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "ticket skipped",
		"ticket":  ticket,
	})
}

func (h *TicketHandler) CompleteCurrent(w http.ResponseWriter, r *http.Request) {
	queueID, err := getIntParam(r, "id")
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid queue id")
		return
	}

	ticket, err := h.service.CompleteCurrent(queueID)
	if err != nil {
		writeTicketActionError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "ticket completed",
		"ticket":  ticket,
	})
}

func getIntParam(r *http.Request, name string) (int, error) {
	params := mux.Vars(r)
	return strconv.Atoi(params[name])
}

func writeTicketActionError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidQueueID), errors.Is(err, service.ErrInvalidTicketID):
		writeJSONError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrTicketNotFound):
		writeJSONError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrNoWaitingTickets), errors.Is(err, service.ErrNoActiveTicket), errors.Is(err, service.ErrActiveTicketExist), errors.Is(err, service.ErrTicketNotSkipped):
		writeJSONError(w, http.StatusConflict, err.Error())
	case isDuplicateKey(err) && hasConstraint(err, "tickets_one_active_per_queue_idx"):
		writeJSONError(w, http.StatusConflict, service.ErrActiveTicketExist.Error())
	default:
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
