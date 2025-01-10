package services

import (
	"estimator-be/internal/models"
	"estimator-be/internal/repositories"
	"testing"
)

func TestRoomService_CreateRoom(t *testing.T) {
	repo := repositories.NewRoomRepository()
	service := NewRoomService(repo)

	// Create a new room
	roomName := "Test Room"
	hostName := "John"
	estimationType := "hours"
	roomID, _, _ := service.CreateRoom(roomName, hostName, estimationType)

	// Verify that the room was created and stored
	savedRoom, exists := repo.FindByID(roomID)
	if !exists {
		t.Fatalf("Expected room with ID %s to exist, but it was not found", roomID)
	}

	if savedRoom.Name != roomName {
		t.Errorf("Expected room name: %s, but got: %s", roomName, savedRoom.Name)
	}

	if savedRoom.EstimationType != models.EstimationHours {
		t.Errorf("expected estimation type 'hours', got '%s'", savedRoom.EstimationType)
	}
}

func TestRoomService_CreateRoomWithInvalidEstimationType(t *testing.T) {
	repo := repositories.NewRoomRepository()
	service := NewRoomService(repo)

	// Create a new room
	roomName := "Test Room"
	hostName := "John"
	estimationType := "hours"
	roomID, _, _ := service.CreateRoom(roomName, hostName, estimationType)

	// Verify that the room was created and stored
	savedRoom, exists := repo.FindByID(roomID)
	if !exists {
		t.Fatalf("Expected room with ID %s to exist, but it was not found", roomID)
	}

	if savedRoom.Name != roomName {
		t.Errorf("Expected room name: %s, but got: %s", roomName, savedRoom.Name)
	}

	if savedRoom.EstimationType != models.EstimationHours {
		t.Errorf("expected estimation type 'hours', got '%s'", savedRoom.EstimationType)
	}
}
