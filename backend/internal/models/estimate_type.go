package models

type EstimationType string

const (
	EstimationHours       EstimationType = "hours"
	EstimationDays        EstimationType = "days"
	EstimationStoryPoints EstimationType = "story_points"
)
