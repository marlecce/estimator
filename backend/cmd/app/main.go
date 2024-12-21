package main

import (
	"estimator-be/internal/api"
	"estimator-be/internal/repositories"
	"estimator-be/internal/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	roomRepo := repositories.NewRoomRepository()
	roomService := services.NewRoomService(roomRepo)

	api.RegisterRoomRoutes(r, roomService)

	port := "8181"

	log.Printf("Server running on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Error running the server: %v", err)
	}
}
