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

func TestRoomService_GetRoomByID(t *testing.T) {
	repo := repositories.NewRoomRepository()
	service := NewRoomService(repo)

	// Create and save a test room
	room := &models.Room{ID: "123456", Name: "Existing Room"}
	repo.Save(room)

	// Retrieve the room using the service
	retrievedRoom, exists := service.GetRoomByID("123456")
	if !exists {
		t.Fatalf("Expected room with ID 123456 to exist, but it was not found")
	}

	if retrievedRoom.Name != "Existing Room" {
		t.Errorf("Expected room name: %s, but got: %s", "Existing Room", retrievedRoom.Name)
	}

	// Try to retrieve a non-existent room
	_, exists = service.GetRoomByID("654321")
	if exists {
		t.Errorf("Expected room with ID 654321 to not exist, but it was found")
	}
}
