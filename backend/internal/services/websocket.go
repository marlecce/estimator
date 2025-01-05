package services

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/coder/websocket"
)

type WebSocketServer struct {
	clients        map[*websocket.Conn]string
	broadcast      chan []byte
	allowedOrigins []string
	mu             sync.Mutex
}

func NewWebSocketServer(allowedOrigins []string) *WebSocketServer {
	return &WebSocketServer{
		clients:        make(map[*websocket.Conn]string),
		broadcast:      make(chan []byte),
		allowedOrigins: allowedOrigins,
	}
}

func (server *WebSocketServer) isValidOrigin(origin string) bool {
	for _, allowedOrigin := range server.allowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}
	return false
}

func (server *WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	if !server.isValidOrigin(origin) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		log.Printf("WebSocket connection rejected: invalid Origin %s", origin)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"localhost:5173"},
	})
	if err != nil {
		log.Printf("WebSocket connection failed: %v", err)
		return
	}

	defer conn.Close(websocket.StatusInternalError, "internal server error")

	var setupMessage map[string]interface{}
	_, msg, err := conn.Read(r.Context())
	if err != nil {
		log.Printf("Error reading setup message: %v", err)
		return
	}

	err = json.Unmarshal(msg, &setupMessage)
	if err != nil {
		log.Printf("Error unmarshalling setup message: %v", err)
		return
	}

	participantId, ok := setupMessage["participantId"].(string)
	if !ok {
		log.Printf("Missing or invalid participantId")
		return
	}

	server.mu.Lock()
	server.clients[conn] = participantId
	server.mu.Unlock()

	log.Printf("User connected with participantId: %s", participantId)

	defer func() {
		server.mu.Lock()
		delete(server.clients, conn)
		server.mu.Unlock()
		log.Printf("User disconnected: %s", participantId)
	}()

	response := map[string]interface{}{
		"type":   "user_id_assigned",
		"userId": participantId,
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
	}

	err = conn.Write(r.Context(), websocket.MessageText, responseBytes)
	if err != nil {
		log.Printf("Error sending user ID: %v", err)
	}

	for {
		_, msg, err := conn.Read(r.Context())
		if websocket.CloseStatus(err) != -1 {
			delete(server.clients, conn)
			return
		}

		var message map[string]interface{}
		err = json.Unmarshal(msg, &message)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		server.SendBroadcast(msg)
	}
}

func (server *WebSocketServer) HandleMessages() {
	for {
		msg := <-server.broadcast
		server.mu.Lock()
		for client := range server.clients {
			ctx := context.Background()
			err := client.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				log.Printf("Error writing message to client: %v", err)
				client.Close(websocket.StatusInternalError, "error sending message")
				delete(server.clients, client)
			}
		}
		server.mu.Unlock()
	}
}

func (server *WebSocketServer) GetCurrentUserId(conn *websocket.Conn) (string, bool) {
	server.mu.Lock()
	defer server.mu.Unlock()
	userID, exists := server.clients[conn]
	return userID, exists
}

func (server *WebSocketServer) SendBroadcast(msg []byte) {
	server.broadcast <- msg
}
