package models

type Ticket struct {
	ID      int    `json:"id"`
	QueueID int    `json:"queue_id"`
	Number  int    `json:"number"`
	Status  string `json:"status"`
}