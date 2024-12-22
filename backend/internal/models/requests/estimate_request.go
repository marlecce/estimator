package models

type EstimateRequest struct {
	ParticipantID string  `json:"participant_id"`
	Value         float64 `json:"value"`
}
