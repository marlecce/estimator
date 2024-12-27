package services

import (
	"estimator-be/internal/models"
	"estimator-be/internal/repositories"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
)

type RoomService struct {
	repo *repositories.RoomRepository
}

func NewRoomService(repo *repositories.RoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) CreateRoom(name string, estimationType string) string {
	var validType models.EstimationType
	switch estimationType {
	case string(models.EstimationHours), string(models.EstimationDays), string(models.EstimationStoryPoints):
		validType = models.EstimationType(estimationType)
	default:
		validType = models.EstimationStoryPoints
	}

	roomID := uuid.New().String()

	room := &models.Room{
		ID:             roomID,
		Name:           name,
		Participants:   []*models.Participant{},
		Estimates:      []*models.Estimate{},
		EstimationType: validType,
	}
	s.repo.Save(room)

	return roomID
}

func (s *RoomService) AddParticipant(roomID, name string) (string, error) {
	participantID := fmt.Sprintf("p-%06d", rand.Intn(1000000))
	participant := &models.Participant{ID: participantID, Name: name}

	err := s.repo.AddParticipant(roomID, participant)
	if err != nil {
		return "", err
	}

	return participantID, nil
}

func (s *RoomService) AddEstimate(roomID, participantID string, value float64) error {
	estimate := &models.Estimate{ParticipantID: participantID, Value: value}
	return s.repo.AddEstimate(roomID, estimate)
}

func (s *RoomService) RevealEstimates(roomID string) ([]*models.Estimate, error) {
	err := s.repo.RevealEstimates(roomID)
	if err != nil {
		return nil, err
	}

	room, exists := s.repo.FindByID(roomID)
	if !exists {
		return nil, fmt.Errorf("room with ID %s not found", roomID)
	}

	return room.Estimates, nil
}

func (s *RoomService) GetRoomDetails(roomID string) (*models.Room, error) {
	room, exists := s.repo.FindByID(roomID)
	if !exists {
		return nil, fmt.Errorf("room with ID %s not found", roomID)
	}
	return room, nil
}
