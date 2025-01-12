package models

import "estimator-be/internal/models"

type CreateRoomRequest struct {
	Name           string                `json:"name"`
	HostName       string                `json:"host_name"`
	EstimationType models.EstimationType `json:"estimation_type"`
}
