package models

// type Queue struct {
// 	ID             int    `json:"id"`
// 	Name           string `json:"name"`
// 	ServicePointID int    `json:"service_point_id"`
// }

type Queue struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	ServicePointID int    `json:"service_point_id"`
	IsOpen         bool   `json:"is_open"` // 👈 ДОБАВЬ
}
