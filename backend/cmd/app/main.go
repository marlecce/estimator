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
	hubService := services.NewHubService()

	apiRouter := router.PathPrefix("/api").Subrouter()
	api.RegisterRoomRoutes(apiRouter, roomService, hubService)

	socketServer := services.InitSocketServer()
	router.PathPrefix("/socket.io/").Handler(socketServer)

	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./frontend/dist/"))))

	handler := cors.Default().Handler(router)

	port := "8181"

	log.Printf("Server running on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Error running the server: %v", err)
	}
}
