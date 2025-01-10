package models

import "estimator-be/internal/models"

type EstimateRequest struct {
	ParticipantID  string                `json:"participant_id"`
	Value          float64               `json:"value"`
	EstimationType models.EstimationType `json:"estimation_type"`
}
