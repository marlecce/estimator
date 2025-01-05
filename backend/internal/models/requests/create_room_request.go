package models

type CreateRoomRequest struct {
	Name           string `json:"name"`
	HostName       string `json:"host_name"`
	EstimationType string `json:"estimation_type"`
}
