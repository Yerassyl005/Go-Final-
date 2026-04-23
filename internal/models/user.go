package models

import "time"

const (
	PriorityCategoryNone     = "none"
	PriorityCategoryPregnant = "pregnant"
	PriorityCategoryElderly  = "elderly"
	PriorityCategoryDisabled = "disabled"
)

type User struct {
	ID               int       `json:"id"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Phone            string    `json:"phone"`
	PriorityCategory string    `json:"priority_category"`
	PasswordHash     string    `json:"-"`
	CreatedAt        time.Time `json:"created_at"`
	ServicePointID int `json:"service_point_id"`
}
