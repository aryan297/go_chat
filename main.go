package main

import (
	"chat-app/config"
	"chat-app/handlers"
	"log"
	"net/http"
)

func main() {
	// Initialize MongoDB connection
	config.InitMongoDB()

	// Setup WebSocket handler
	http.HandleFunc("/ws/", handlers.HandleWebSocket)

	// Start HTTP server
	serverAddr := ":8080"
	log.Printf("WebSocket server is running on %s...", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
