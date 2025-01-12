package api

import (
	"bytes"
	"encoding/json"
	"estimator-be/internal/repositories"
	"estimator-be/internal/services"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestCreateRoomIntegration(t *testing.T) {
	// Arrange
	roomRepo := repositories.NewRoomRepository()
	roomService := services.NewRoomService(roomRepo)

	allowedOrigins := []string{"http://localhost:5173"}
	wsService := services.NewWebSocketServer(allowedOrigins)

	router := mux.NewRouter()
	RegisterRoomRoutes(router, roomService, wsService)

	payload := map[string]string{"name": "Test Room"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/rooms", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Act
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["room_id"])
}

func TestAddParticipantIntegration(t *testing.T) {
	// Arrange
	roomRepo := repositories.NewRoomRepository()
	roomService := services.NewRoomService(roomRepo)

	allowedOrigins := []string{"http://localhost:5173"}
	wsService := services.NewWebSocketServer(allowedOrigins)

	router := mux.NewRouter()
	RegisterRoomRoutes(router, roomService, wsService)

	// Create a room
	roomPayload := map[string]string{"name": "Test Room"}
	roomBody, _ := json.Marshal(roomPayload)
	roomReq, _ := http.NewRequest("POST", "/rooms", bytes.NewReader(roomBody))
	roomReq.Header.Set("Content-Type", "application/json")
	roomResp := httptest.NewRecorder()
	router.ServeHTTP(roomResp, roomReq)

	assert.Equal(t, http.StatusCreated, roomResp.Code)

	var roomResponse map[string]string
	err := json.Unmarshal(roomResp.Body.Bytes(), &roomResponse)
	assert.NoError(t, err)
	roomID := roomResponse["room_id"]

	// Add a partecipant
	participantPayload := map[string]string{"name": "John Doe"}
	participantBody, _ := json.Marshal(participantPayload)
	participantReq, _ := http.NewRequest("POST", fmt.Sprintf("/rooms/%s/join", roomID), bytes.NewReader(participantBody))
	participantReq.Header.Set("Content-Type", "application/json")
	participantResp := httptest.NewRecorder()
	router.ServeHTTP(participantResp, participantReq)

	// Assert
	assert.Equal(t, http.StatusOK, participantResp.Code)

	var participantResponse map[string]string
	err = json.Unmarshal(participantResp.Body.Bytes(), &participantResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, participantResponse["participant_id"])
}
