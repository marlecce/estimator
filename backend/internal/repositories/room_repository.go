package repositories

import (
	"estimator-be/internal/models"
	"sync"
)

type RoomRepository struct {
	mu    sync.RWMutex
	rooms map[string]*models.Room
}

func NewRoomRepository() *RoomRepository {
	return &RoomRepository{
		rooms: make(map[string]*models.Room),
	}
}

func (r *RoomRepository) Save(room *models.Room) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rooms[room.ID] = room
}

func (r *RoomRepository) FindByID(id string) (*models.Room, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	room, exists := r.rooms[id]
	return room, exists
}

func (r *RoomRepository) FindAll() []*models.Room {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rooms := make([]*models.Room, 0, len(r.rooms))
	for _, room := range r.rooms {
		rooms = append(rooms, room)
	}
	return rooms
}
