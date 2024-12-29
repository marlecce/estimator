package models

import socketio "github.com/googollee/go-socket.io"

type Hub struct {
	RoomID    string
	Clients   map[string]socketio.Conn
	Broadcast chan []byte
}
