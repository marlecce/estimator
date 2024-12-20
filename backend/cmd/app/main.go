package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/create-room", createRoom)
	http.HandleFunc("/join-room", joinRoom)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createRoom(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Room created")
}

func joinRoom(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Joined room")
}
