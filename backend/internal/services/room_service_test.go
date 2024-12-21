package services

import (
	"estimator-be/internal/repositories"
	"testing"
)

func TestRoomService_CreateRoom(t *testing.T) {
	repo := repositories.NewRoomRepository()
	service := NewRoomService(repo)

	// Create a new room
	roomName := "Test Room"
	roomID := service.CreateRoom(roomName)

	// Verify that the room was created and stored
	savedRoom, exists := repo.FindByID(roomID)
	if !exists {
		t.Fatalf("Expected room with ID %s to exist, but it was not found", roomID)
	}

	if savedRoom.Name != roomName {
		t.Errorf("Expected room name: %s, but got: %s", roomName, savedRoom.Name)
	}
}
