package models

type CreateRoomRequest struct {
	Name           string `json:"name"`
	EstimationType string `json:"estimation_type"`
}
