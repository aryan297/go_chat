package handlers

import (
	"chat-app/config"
	"chat-app/models"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	db := config.GetDB()
	collection := db.Collection("chats")

	for {
		// Read message from client
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading WebSocket message: %v", err)
			break
		}

		// Unmarshal the message into a ChatMessage struct
		var chatMsg models.Chat
		err = json.Unmarshal(msg, &chatMsg)
		if err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		// Check if userId is provided
		if chatMsg.UserID == "" {
			log.Println("User ID not provided")
			continue
		}

		// Add timestamp
		chatMsg.Timestamp = time.Now()

		// Save to MongoDB
		_, err = collection.InsertOne(r.Context(), chatMsg)
		if err != nil {
			log.Printf("Failed to save message: %v", err)
			continue
		}

		// Echo message back to client
		err = conn.WriteJSON(chatMsg)
		if err != nil {
			log.Printf("Error sending message back: %v", err)
			break
		}
	}
}
