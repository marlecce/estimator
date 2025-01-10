package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"estimator-be/internal/api/validators"
	requests "estimator-be/internal/models/requests"
	"estimator-be/internal/services"

	"github.com/gorilla/mux"
)

type RoomHandler struct {
	roomService *services.RoomService
	wsServer    *services.WebSocketServer
}

func NewRoomHandler(roomService *services.RoomService, wsServer *services.WebSocketServer) *RoomHandler {
	return &RoomHandler{
		roomService: roomService,
		wsServer:    wsServer,
	}
}

func RegisterRoomRoutes(r *mux.Router, roomService *services.RoomService, wsServer *services.WebSocketServer) {
	handler := NewRoomHandler(roomService, wsServer)

	r.HandleFunc("/rooms", handler.CreateRoom).Methods("POST")
	r.HandleFunc("/rooms/{room_id}/join", handler.JoinRoom).Methods("POST")
	r.HandleFunc("/rooms/{room_id}/estimate", handler.Estimate).Methods("POST")
	r.HandleFunc("/rooms/{room_id}/reveal", handler.Reveal).Methods("POST")
	r.HandleFunc("/rooms/{room_id}", handler.GetRoomDetails).Methods("GET")
}

func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var req *requests.CreateRoomRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	roomID, host, _ := h.roomService.CreateRoom(req.Name, req.HostName, req.EstimationType)

	resp := map[string]interface{}{
		"room_id": roomID,
		"host":    host,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *RoomHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	var req *requests.JoinRoomRequest
	vars := mux.Vars(r)
	roomID := vars["room_id"]

	// TODO check if room exists
	/*
		if !h.roomService.RoomExists(roomID) {
			http.Error(w, "Room does not exist", http.StatusNotFound)
			return
		}
	*/

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	isHost := false
	participant, err := h.roomService.AddParticipant(roomID, req.Name, isHost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := map[string]interface{}{
		"participant": participant,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *RoomHandler) Estimate(w http.ResponseWriter, r *http.Request) {
	var req *requests.EstimateRequest
	vars := mux.Vars(r)
	roomID := vars["room_id"]

	fmt.Print(r.Body)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	fmt.Print(req)

	if err := validators.ValidateEstimateRequest(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.roomService.AddEstimate(roomID, req.ParticipantID, req.Value, req.EstimationType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg, _ := json.Marshal(map[string]interface{}{
		"type":        "estimate_submitted",
		"participant": req.ParticipantID,
		"message":     "A participant has made an estimate.",
	})
	h.wsServer.SendBroadcast(msg)

	resp := map[string]string{"status": "success"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *RoomHandler) Reveal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["room_id"]

	estimates, err := h.roomService.RevealEstimates(roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg, _ := json.Marshal(map[string]interface{}{
		"type":      "estimates_revealed",
		"estimates": estimates,
	})
	h.wsServer.SendBroadcast(msg)

	resp := map[string]interface{}{"estimates": estimates}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *RoomHandler) GetRoomDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["room_id"]

	roomDetails, err := h.roomService.GetRoomDetails(roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := map[string]interface{}{
		"name":           roomDetails.Name,
		"participants":   roomDetails.Participants,
		"hostId":         roomDetails.HostID,
		"revealed":       roomDetails.Revealed,
		"estimationType": roomDetails.EstimationType,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
