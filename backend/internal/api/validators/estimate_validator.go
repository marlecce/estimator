package validators

import (
	"estimator-be/internal/models"
	requests "estimator-be/internal/models/requests"
	"fmt"
)

func ValidateEstimateRequest(req *requests.EstimateRequest) error {
	if req.ParticipantID == "" {
		return fmt.Errorf("ParticipantID is required")
	}

	if req.Value <= 0 {
		return fmt.Errorf("value must be greater than 0")
	}

	if !IsValidEstimationType(req.EstimationType) {
		return fmt.Errorf("invalid estimation type: %s", req.EstimationType)
	}

	return nil
}

func IsValidEstimationType(estimationType models.EstimationType) bool {
	switch estimationType {
	case models.EstimationHours, models.EstimationDays, models.EstimationStoryPoints:
		return true
	default:
		return false
	}
}
