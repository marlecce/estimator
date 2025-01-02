package services

import (
	"context"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

type WebSocketServer struct {
	clients        map[*websocket.Conn]bool
	broadcast      chan []byte
	allowedOrigins []string
}

func NewWebSocketServer(allowedOrigins []string) *WebSocketServer {
	return &WebSocketServer{
		clients:        make(map[*websocket.Conn]bool),
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
	server.clients[conn] = true

	for {
		_, msg, err := conn.Read(r.Context())
		if websocket.CloseStatus(err) != -1 {
			delete(server.clients, conn)
			return
		}

		server.broadcast <- msg
	}
}

func (server *WebSocketServer) HandleMessages() {
	for {
		msg := <-server.broadcast
		for client := range server.clients {
			ctx := context.Background()
			err := client.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				log.Printf("Error writing message to client: %v", err)
				client.Close(websocket.StatusInternalError, "error sending message")
				delete(server.clients, client)
			}
		}
	}
}

func (server *WebSocketServer) SendBroadcast(msg []byte) {
	server.broadcast <- msg
}
