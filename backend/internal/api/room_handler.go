package api

import (
	"encoding/json"
	"log"
	"net/http"

	requests "estimator-be/internal/models/requests"
	"estimator-be/internal/services"

	socketio "github.com/googollee/go-socket.io"

	"github.com/gorilla/mux"
)

var roomSocketServers = make(map[string]*socketio.Server)

type RoomHandler struct {
	roomService *services.RoomService
	hubService  *services.HubService
}

func NewRoomHandler(roomService *services.RoomService, hubService *services.HubService) *RoomHandler {
	return &RoomHandler{
		roomService: roomService,
		hubService:  hubService,
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

	// TODO check if room exists

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	participantID, err := h.roomService.AddParticipant(roomID, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	server, exists := roomSocketServers[roomID]
	if !exists {
		server = socketio.NewServer(nil)

		// Configura gli eventi WebSocket
		server.OnConnect("/", func(s socketio.Conn) error {
			log.Printf("New connection in room %s", roomID)
			s.Join(roomID)
			return nil
		})

		server.OnEvent("/", "join_room", func(s socketio.Conn, data map[string]interface{}) {
			participantID := data["participantId"].(string)
			log.Printf("Participant %s joined room %s", participantID, roomID)
			s.Emit("participant_joined", data)
		})

		server.OnDisconnect("/", func(s socketio.Conn, reason string) {
			log.Printf("Participant disconnected from room %s: %s", roomID, reason)
			s.Leave(roomID)
		})

		go func() {
			if err := server.Serve(); err != nil {
				log.Fatalf("SocketIO server error: %s", err)
			}
		}()
		roomSocketServers[roomID] = server
	} else {
		log.Printf("WebSocket server already exists for room %s", roomID)
	}

	// Invia un evento di aggiornamento della stanza
	server.BroadcastToRoom("/", roomID, "participant_joined", map[string]interface{}{
		"id":   participantID,
		"name": req.Name,
	})

	// Risposta HTTP
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

	hub, exists := h.hubService.GetHub(roomID)
	if exists {
		message := map[string]interface{}{
			"type":        "estimate_notification",
			"participant": req.ParticipantID,
			"message":     "A participant has made an estimate.",
		}
		msg, _ := json.Marshal(message)
		hub.Broadcast <- msg
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

	hub, exists := h.hubService.GetHub(roomID)
	if exists {
		message := map[string]interface{}{
			"type":      "estimates_revealed",
			"estimates": estimates,
		}
		msg, _ := json.Marshal(message)
		hub.Broadcast <- msg
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

	server, exists := roomSocketServers[roomID]
	if !exists {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	server.ServeHTTP(w, r)
}
