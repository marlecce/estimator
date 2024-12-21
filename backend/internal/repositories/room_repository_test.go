package repositories

import (
	"estimator-be/internal/models"
	"testing"
)

func TestRoomRepository(t *testing.T) {
	repo := NewRoomRepository()

	room := &models.Room{ID: "123456", Name: "Test Room"}
	repo.Save(room)

	savedRoom, exists := repo.FindByID("123456")
	if !exists || savedRoom.Name != "Test Room" {
		t.Errorf("Expected room: %v, but got: %v", room, savedRoom)
	}

	_, exists = repo.FindByID("654321")
	if exists {
		t.Errorf("A non-existing room ID should not be found.")
	}
}
