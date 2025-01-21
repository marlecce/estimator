package api

import (
	"bytes"
	"encoding/json"
	"estimator-be/internal/models"
	requests "estimator-be/internal/models/requests"
	responses "estimator-be/internal/models/responses"
	"estimator-be/internal/repositories"
	"estimator-be/internal/services"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *mux.Router {
	roomRepo := repositories.NewRoomRepository()
	roomService := services.NewRoomService(roomRepo)
	allowedOrigins := []string{"http://localhost:5173"}
	wsService := services.NewWebSocketServer(allowedOrigins)
	router := mux.NewRouter()
	RegisterRoomRoutes(router, roomService, wsService)
	return router
}

func TestCreateRoom(t *testing.T) {
	t.Run("should create room successfully", func(t *testing.T) {
		// Arrange
		router := setupRouter()
		payload := requests.CreateRoomRequest{
			Name:           "Room Test",
			HostName:       "Alice",
			EstimationType: models.EstimationHours,
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/rooms", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		// Act
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Assert
		assert.Equal(t, http.StatusCreated, rr.Code)
		var response responses.CreatedRoomResponse
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.RoomID)
		assert.NotEmpty(t, response.Host.ID)
		assert.Equal(t, "Alice", response.Host.Name)
	})

	t.Run("should return 400 for invalid payload", func(t *testing.T) {
		// Arrange
		router := setupRouter()
		req, _ := http.NewRequest("POST", "/rooms", bytes.NewReader([]byte("invalid_payload")))
		req.Header.Set("Content-Type", "application/json")

		// Act
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid request")
	})
}

func TestJoinRoomIntegration(t *testing.T) {
	t.Run("should join room successfully", func(t *testing.T) {
		// Arrange
		router := setupRouter()

		// Create a room
		roomPayload := requests.CreateRoomRequest{
			Name:           "Room Test",
			HostName:       "Alice",
			EstimationType: models.EstimationStoryPoints,
		}
		roomBody, _ := json.Marshal(roomPayload)
		roomReq, _ := http.NewRequest("POST", "/rooms", bytes.NewReader(roomBody))
		roomReq.Header.Set("Content-Type", "application/json")
		roomResp := httptest.NewRecorder()
		router.ServeHTTP(roomResp, roomReq)
		assert.Equal(t, http.StatusCreated, roomResp.Code)

		var roomResponse responses.CreatedRoomResponse
		err := json.Unmarshal(roomResp.Body.Bytes(), &roomResponse)
		assert.NoError(t, err)
		roomID := roomResponse.RoomID

		// Join the room
		participantPayload := requests.JoinRoomRequest{Name: "Bob"}
		participantBody, _ := json.Marshal(participantPayload)
		participantReq, _ := http.NewRequest("POST", fmt.Sprintf("/rooms/%s/join", roomID), bytes.NewReader(participantBody))
		participantReq.Header.Set("Content-Type", "application/json")
		participantResp := httptest.NewRecorder()
		router.ServeHTTP(participantResp, participantReq)

		// Assert
		assert.Equal(t, http.StatusOK, participantResp.Code)
		var participantResponse models.Participant
		err = json.Unmarshal(participantResp.Body.Bytes(), &participantResponse)
		assert.NoError(t, err)
		assert.NotEmpty(t, participantResponse.ID)
		assert.Equal(t, "Bob", participantResponse.Name)
	})

	t.Run("should return 400 for non-existent room", func(t *testing.T) {
		// Arrange
		router := setupRouter()
		participantPayload := requests.JoinRoomRequest{Name: "Bob"}
		participantBody, _ := json.Marshal(participantPayload)
		participantReq, _ := http.NewRequest("POST", "/rooms/nonexistent/join", bytes.NewReader(participantBody))
		participantReq.Header.Set("Content-Type", "application/json")

		// Act
		participantResp := httptest.NewRecorder()
		router.ServeHTTP(participantResp, participantReq)

		// Assert
		assert.Equal(t, http.StatusBadRequest, participantResp.Code)
		assert.Contains(t, participantResp.Body.String(), "room with ID nonexistent not found")
	})
}

func TestGetRoomDetails(t *testing.T) {
	t.Run("should return room details for valid room ID", func(t *testing.T) {
		// Arrange
		roomRepo := repositories.NewRoomRepository()
		roomService := services.NewRoomService(roomRepo)

		allowedOrigins := []string{"http://localhost:5173"}
		wsService := services.NewWebSocketServer(allowedOrigins)

		router := mux.NewRouter()
		RegisterRoomRoutes(router, roomService, wsService)

		roomName := "Test"
		hostName := "John Doe"
		estimationType := models.EstimationHours
		roomId, host, _ := roomService.CreateRoom(roomName, hostName, estimationType)

		req, _ := http.NewRequest("GET", fmt.Sprintf("/rooms/%s", roomId), nil)

		// Act
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Assert
		assert.Equal(t, http.StatusOK, rr.Code)

		var response models.Room
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, roomId, response.ID)
		assert.Equal(t, roomName, response.Name)
		assert.Equal(t, host, response.Participants[0])
	})

	t.Run("should return 400 if room does not exist", func(t *testing.T) {
		// Arrange
		roomRepo := repositories.NewRoomRepository()
		roomService := services.NewRoomService(roomRepo)

		allowedOrigins := []string{"http://localhost:5173"}
		wsService := services.NewWebSocketServer(allowedOrigins)

		router := mux.NewRouter()
		RegisterRoomRoutes(router, roomService, wsService)

		roomID := "nonexistent_room"

		req, _ := http.NewRequest("GET", fmt.Sprintf("/rooms/%s", roomID), nil)

		// Act
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "room not found")
	})
}
