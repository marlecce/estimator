package models

type Room struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Participants []*Participant `json:"participants"`
	Estimates    []*Estimate    `json:"estimates"`
	Revealed     bool           `json:"revealed"`
}
