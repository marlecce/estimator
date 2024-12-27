package services

import (
	"estimator-be/internal/models"
	"log"
	"sync"

	"github.com/gorilla/websocket"
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
		RoomID:     roomID,
		Clients:    make(map[*websocket.Conn]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
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
			RoomID:     roomID,
			Clients:    make(map[*websocket.Conn]bool),
			Broadcast:  make(chan []byte),
			Register:   make(chan *websocket.Conn),
			Unregister: make(chan *websocket.Conn),
		}
		s.hubs[roomID] = hub
	}
	return hub
}

func (s *HubService) AddClient(hub *models.Hub, conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	hub.Clients[conn] = true
}

func (s *HubService) RemoveClient(hub *models.Hub, conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(hub.Clients, conn)
}

func (s *HubService) ListenForMessages(hub *models.Hub, conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			s.RemoveClient(hub, conn)
			conn.Close()
			break
		}
		hub.Broadcast <- msg
	}
}

func (h *HubService) Run(hub *models.Hub) {
	for {
		select {
		case conn := <-hub.Register:
			hub.Clients[conn] = true
			log.Printf("Client registered in room %s", hub.RoomID)
		case conn := <-hub.Unregister:
			if _, ok := hub.Clients[conn]; ok {
				delete(hub.Clients, conn)
				conn.Close()
				log.Printf("Client unregistered from room %s", hub.RoomID)
			}
		case message := <-hub.Broadcast:
			for conn := range hub.Clients {
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Println("Error sending message:", err)
					conn.Close()
					hub.Unregister <- conn
				}
			}
		}
	}
}
