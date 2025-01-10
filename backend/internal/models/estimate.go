package models

type Estimate struct {
	ParticipantID  string         `json:"participant_id"`
	Value          float64        `json:"value"`
	EstimationType EstimationType `json:"estimation_type"`
}
