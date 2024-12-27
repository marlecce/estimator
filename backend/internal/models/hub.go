package models

import (
	"github.com/gorilla/websocket"
)

type Hub struct {
	RoomID     string
	Clients    map[*websocket.Conn]bool
	Broadcast  chan []byte
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
}
