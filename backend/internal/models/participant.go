package models

type Participant struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	HasEstimated bool   `json:"hasEstimated"`
}
