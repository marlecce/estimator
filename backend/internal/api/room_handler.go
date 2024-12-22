package api

import (
	"encoding/json"
	"net/http"

	requests "estimator-be/internal/models/requests"
	"estimator-be/internal/services"

	"github.com/gorilla/mux"
)

type RoomHandler struct {
	roomService *services.RoomService
}

func RegisterRoomRoutes(r *mux.Router, roomService *services.RoomService) {
	handler := &RoomHandler{roomService: roomService}
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

	roomID := h.roomService.CreateRoom(req.Name, req.EstimationType)

	resp := map[string]string{"room_id": roomID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *RoomHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	var req *requests.JoinRoomRequest
	vars := mux.Vars(r)
	roomID := vars["room_id"]

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	participantID, err := h.roomService.AddParticipant(roomID, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := map[string]string{"participant_id": participantID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *RoomHandler) Estimate(w http.ResponseWriter, r *http.Request) {
	var req *requests.EstimateRequest
	vars := mux.Vars(r)
	roomID := vars["room_id"]

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := h.roomService.AddEstimate(roomID, req.ParticipantID, req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
		"name":            roomDetails.Name,
		"participants":    roomDetails.Participants,
		"revealed":        roomDetails.Revealed,
		"estimation_type": roomDetails.EstimationType,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
