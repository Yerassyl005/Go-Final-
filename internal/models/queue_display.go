package models

type QueueDisplayTicket struct {
	ID           int    `json:"id"`
	TicketNumber string `json:"ticket_number"`
	Status       string `json:"status"`
	RecallCount  int    `json:"recall_count,omitempty"`
}

type QueueDisplay struct {
	QueueID        int                  `json:"queue_id"`
	QueueName      string               `json:"queue_name"`
	CurrentTicket  *QueueDisplayTicket  `json:"current_ticket,omitempty"`
	WaitingTickets []QueueDisplayTicket `json:"waiting_tickets"`
	CompletedCount int                  `json:"completed_count"`
}

type QueueStats struct {
	QueueID          int `json:"queue_id"`
	TotalTickets     int `json:"total_tickets"`
	WaitingTickets   int `json:"waiting_tickets"`
	CalledTickets    int `json:"called_tickets"`
	CompletedTickets int `json:"completed_tickets"`
	SkippedTickets   int `json:"skipped_tickets"`
}
