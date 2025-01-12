package responses

import "estimator-be/internal/models"

type CreatedRoomResponse struct {
	RoomID string             `json:"room_id"`
	Host   models.Participant `json:"host"`
}
