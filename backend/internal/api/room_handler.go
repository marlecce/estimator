package api

import (
	"encoding/json"
	"log"
	"net/http"

	"estimator-be/internal/models"
	requests "estimator-be/internal/models/requests"
	"estimator-be/internal/services"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type RoomHandler struct {
	roomService *services.RoomService
	hubService  *services.HubService
	upgrader    websocket.Upgrader
}

func NewRoomHandler(roomService *services.RoomService, hubService *services.HubService) *RoomHandler {
	return &RoomHandler{
		roomService: roomService,
		hubService:  hubService,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// TODO Permettiamo tutte le origini, eventualmente si può aggiungere logica di validazione.
				return true
			},
		},
	}
}

func RegisterRoomRoutes(r *mux.Router, roomService *services.RoomService, hubService *services.HubService) {
	handler := NewRoomHandler(roomService, hubService)
	r.HandleFunc("/rooms", handler.CreateRoom).Methods("POST")
	r.HandleFunc("/rooms/{room_id}/join", handler.JoinRoom).Methods("POST")
	r.HandleFunc("/rooms/{room_id}/estimate", handler.Estimate).Methods("POST")
	r.HandleFunc("/rooms/{room_id}/reveal", handler.Reveal).Methods("POST")
	r.HandleFunc("/rooms/{room_id}", handler.GetRoomDetails).Methods("GET")
	r.HandleFunc("/rooms/{room_id}/ws", handler.HandleWebSocket).Methods("GET")
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
	w.WriteHeader(http.StatusCreated)
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

func (h *RoomHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["room_id"]
	log.Printf("Starting websocket for room %s", roomID)

	participantName := r.URL.Query().Get("name")
	if participantName == "" {
		http.Error(w, "Missing 'name' query parameter", http.StatusBadRequest)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	participantID, err := h.roomService.AddParticipant(roomID, participantName)
	if err != nil {
		log.Println("Error adding participant:", err)
		http.Error(w, "Failed to add participant to room", http.StatusInternalServerError)
		return
	}
	log.Printf("Participant added: ID=%s, Name=%s, Room=%s\n", participantID, participantName, roomID)

	hub := h.hubService.GetOrCreateHub(roomID)

	h.hubService.AddClient(hub, conn)

	h.listenForMessages(hub, conn)
}

func (h *RoomHandler) listenForMessages(hub *models.Hub, conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			h.hubService.RemoveClient(hub, conn)
			return
		}

		hub.Broadcast <- message
	}
}
