package services

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
)

var server *socketio.Server

func InitSocketServer() *socketio.Server {
	server = socketio.NewServer(nil)

	// Evento di connessione
	server.OnConnect("/", func(s socketio.Conn) error {
		log.Printf("Client connected: %s", s.ID())
		return nil
	})

	// Gestisci l'ingresso in una stanza
	server.OnEvent("/", "join_room", func(s socketio.Conn, roomID string) {
		s.Join(roomID) // Aggiungi il client alla stanza specificata
		log.Printf("Client %s joined room %s", s.ID(), roomID)
	})

	// Gestisci l'invio di un messaggio a una stanza
	server.OnEvent("/", "send_message", func(s socketio.Conn, roomID string, message string) {
		server.BroadcastToRoom("/", roomID, "new_message", message)
		log.Printf("Message sent to room %s: %s", roomID, message)
	})

	// Evento di disconnessione
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Printf("Client disconnected: %s, Reason: %s", s.ID(), reason)
	})

	return server
}
