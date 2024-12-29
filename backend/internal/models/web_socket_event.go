package models

type WebSocketEvent struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}
