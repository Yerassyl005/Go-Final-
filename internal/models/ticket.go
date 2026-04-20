package models

type Ticket struct {
	ID          int    `json:"id"`
	QueueID     int    `json:"queue_id"`
	Number      int    `json:"number"`
	Status      string `json:"status"`
	RecallCount int    `json:"recall_count"`
}

const (
	TicketStatusWaiting   = "waiting"
	TicketStatusCalled    = "called"
	TicketStatusSkipped   = "skipped"
	TicketStatusCompleted = "completed"
)
