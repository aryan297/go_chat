package main

import (
	"chat-app/config"
	"chat-app/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Initialize MongoDB connection
	config.InitMongoDB()

	// Setup HTTP handlers
	http.HandleFunc("/ws", handlers.HandleWebSocket)

	// Start the WebSocket server
	serverAddr := ":8080"
	server := &http.Server{Addr: serverAddr}

	// Run the server in a goroutine so it doesn't block
	go func() {
		log.Printf("WebSocket server is running on %s...", serverAddr)
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	// Graceful shutdown: wait for current connections to close
	log.Println("Shutting down server...")

	// Set a deadline for shutting down
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown failed: %v", err)
	}

	// Close MongoDB connection
	config.CloseMongoDB()
	log.Println("Server gracefully shut down")
}
