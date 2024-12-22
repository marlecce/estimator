package models

type Room struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Participants   []*Participant `json:"participants"`
	Revealed       bool           `json:"revealed"`
	Estimates      []*Estimate    `json:"estimates"`
	EstimationType EstimationType `json:"estimation_type"`
}
