package services

import (
	"encoding/json"
	"estimator-be/internal/models"
	"fmt"
	"log"
	"sync"

	socketio "github.com/googollee/go-socket.io"
)

type HubService struct {
	hubs map[string]*models.Hub
	mu   sync.Mutex
}

func NewHubService() *HubService {
	return &HubService{
		hubs: make(map[string]*models.Hub),
	}
}

func (s *HubService) NewHub(roomID string) *models.Hub {
	s.mu.Lock()
	defer s.mu.Unlock()

	hub := &models.Hub{
		RoomID:    roomID,
		Clients:   make(map[string]socketio.Conn),
		Broadcast: make(chan []byte),
	}
	s.hubs[roomID] = hub
	return hub
}

func (s *HubService) GetHub(roomID string) (*models.Hub, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	hub, exists := s.hubs[roomID]
	return hub, exists
}

func (s *HubService) GetOrCreateHub(roomID string) *models.Hub {
	s.mu.Lock()
	defer s.mu.Unlock()

	hub, exists := s.hubs[roomID]
	if !exists {
		hub = &models.Hub{
			RoomID:    roomID,
			Clients:   make(map[string]socketio.Conn),
			Broadcast: make(chan []byte),
		}
		s.hubs[roomID] = hub
	}
	return hub
}

func (s *HubService) AddClient(hub *models.Hub, conn socketio.Conn, participantID string) {
	hub.Clients[participantID] = conn
	conn.Join(hub.RoomID)

	log.Printf("Participant %s joined room %s", participantID, hub.RoomID)
}

func (s *HubService) RemoveClient(hub *models.Hub, participantID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	conn, exists := hub.Clients[participantID]
	if exists {
		conn.Leave(hub.RoomID)
		delete(hub.Clients, participantID)
		log.Printf("Participant %s left room %s", participantID, hub.RoomID)
	}
}

func (s *HubService) ListenForMessages(hub *models.Hub, server *socketio.Server, participantID string) {
	server.OnEvent("/", "message", func(s socketio.Conn, msg string) {
		log.Printf("Message received from %s in room %s: %s", participantID, hub.RoomID, msg)
		hub.Broadcast <- []byte(msg)
	})

	server.OnDisconnect("/", func(conn socketio.Conn, reason string) {
		log.Printf("Participant %s disconnected: %s", participantID, reason)
		s.RemoveClient(hub, participantID)
	})
}

func (s *HubService) Run(hub *models.Hub) {
	go func() {
		for message := range hub.Broadcast {
			for participantID, conn := range hub.Clients {
				conn.Emit("message", string(message))
				log.Printf("Message sent to participant %s: %s", participantID, string(message))
			}
		}
	}()
}

func (s *HubService) SendBroadcastMessage(hub *models.Hub, message models.WebSocketEvent) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to serialize message: %v", err)
		return fmt.Errorf("failed to serialize message: %w", err)
	}
	hub.Broadcast <- messageBytes
	return nil
}
