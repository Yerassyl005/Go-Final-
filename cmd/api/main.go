package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"smartqueue/internal/handler"
	"smartqueue/internal/repository"
	"smartqueue/internal/service"
)

func main() {

	router := mux.NewRouter()

	servicePointRepo := repository.NewServicePointRepository()
	queueRepo := repository.NewQueueRepository()
	ticketRepo := repository.NewTicketRepository()

	servicePointService := service.NewServicePointService(servicePointRepo)
	queueService := service.NewQueueService(queueRepo)
	ticketService := service.NewTicketService(ticketRepo)

	servicePointHandler := handler.NewServicePointHandler(servicePointService)
	queueHandler := handler.NewQueueHandler(queueService)
	ticketHandler := handler.NewTicketHandler(ticketService)

	router.HandleFunc("/servicepoints", servicePointHandler.Create).Methods("POST")
	router.HandleFunc("/servicepoints", servicePointHandler.GetAll).Methods("GET")

	router.HandleFunc("/queues", queueHandler.Create).Methods("POST")
	router.HandleFunc("/queues", queueHandler.GetAll).Methods("GET")

	router.HandleFunc("/tickets", ticketHandler.TakeTicket).Methods("POST")
	router.HandleFunc("/tickets", ticketHandler.GetTickets).Methods("GET")

	router.HandleFunc("/tickets/call", ticketHandler.CallNext).Methods("POST")

	router.HandleFunc("/tickets/{id}/complete", ticketHandler.CompleteTicket).Methods("POST")

	log.Println("Server running on port 8080")

	http.ListenAndServe(":8080", router)
}