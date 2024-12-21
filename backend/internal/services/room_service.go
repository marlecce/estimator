package services

import (
	"estimator-be/internal/models"
	"estimator-be/internal/repositories"
	"fmt"
	"math/rand"
	"time"
)

type RoomService struct {
	repo *repositories.RoomRepository
}

func NewRoomService(repo *repositories.RoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) CreateRoom(name string) string {
	rand.Seed(time.Now().UnixNano())
	roomID := fmt.Sprintf("%06d", rand.Intn(1000000)) // Genera un ID casuale

	room := &models.Room{
		ID:   roomID,
		Name: name,
	}
	s.repo.Save(room)

	return roomID
}

func (s *RoomService) GetRoomByID(id string) (*models.Room, bool) {
	return s.repo.FindByID(id)
}
