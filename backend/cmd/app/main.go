package main

import (
	"estimator-be/internal/api"
	"estimator-be/internal/repositories"
	"estimator-be/internal/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()

	roomRepo := repositories.NewRoomRepository()
	roomService := services.NewRoomService(roomRepo)

	allowedOrigins := []string{"http://localhost:5173"}
	wsServer := services.NewWebSocketServer(allowedOrigins)
	router.HandleFunc("/ws", wsServer.HandleConnections)
	go wsServer.HandleMessages()

	apiRouter := router.PathPrefix("/api").Subrouter()
	api.RegisterRoomRoutes(apiRouter, roomService, wsServer)

	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./frontend/dist/"))))

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
	}).Handler(router)

	port := "8181"

	log.Printf("Server running on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Error running the server: %v", err)
	}
}
