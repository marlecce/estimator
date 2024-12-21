package api

import (
	"encoding/json"
	"net/http"

	"estimator-be/internal/services"

	"github.com/gorilla/mux"
)

type CreateRoomRequest struct {
	Name string `json:"name"`
}

type RoomHandler struct {
	roomService *services.RoomService
}

func RegisterRoomRoutes(r *mux.Router, roomService *services.RoomService) {
	handler := &RoomHandler{roomService: roomService}
	r.HandleFunc("/rooms", handler.CreateRoom).Methods("POST")
}

func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var req CreateRoomRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	roomID := h.roomService.CreateRoom(req.Name)

	resp := map[string]string{"room_id": roomID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
