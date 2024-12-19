package handlers

import (
	"chat-app/config"
	"chat-app/models"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var clients = make(map[string]*websocket.Conn) // Mapping of user IDs to their WebSocket connection

// WebSocket message struct
type WSMessage struct {
	FromUser string `json:"fromUser"`
	ToUser   string `json:"toUser"`
	Message  string `json:"message"`
}

// HandleWebSocket manages the WebSocket connection for users
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	// Get user ID from URL path
	userID := r.URL.Path[len("/ws/"):]

	// Add user to active clients map
	clients[userID] = conn
	log.Printf("User %s connected", userID)

	// Ensure to remove user from map when connection closes
	defer func() {
		delete(clients, userID)
		log.Printf("User %s disconnected", userID)
	}()

	// Listen for incoming messages
	for {
		var msg WSMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		// Send message to the target user if they are connected
		sendMessageToUser(msg)

		// Save message in MongoDB
		saveMessageInDB(msg)
	}
}

// sendMessageToUser sends a message to the intended recipient via WebSocket
func sendMessageToUser(msg WSMessage) {
	toUserConn, ok := clients[msg.ToUser]
	if ok {
		// Send message to the recipient
		err := toUserConn.WriteJSON(msg)
		if err != nil {
			log.Printf("Error sending message to %s: %v", msg.ToUser, err)
		} else {
			log.Printf("Message sent from %s to %s: %s", msg.FromUser, msg.ToUser, msg.Message)
		}
	} else {
		log.Printf("User %s not connected", msg.ToUser)
	}
}

// saveMessageInDB saves the message to MongoDB for persistence
func saveMessageInDB(msg WSMessage) {
	message := models.Message{
		FromUser:  msg.FromUser,
		ToUser:    msg.ToUser,
		Content:   msg.Message,
		Timestamp: time.Now(),
	}

	// Get MongoDB collection
	collection := config.GetDB().Collection("messages")

	// Insert the message into MongoDB
	_, err := collection.InsertOne(nil, message)
	if err != nil {
		log.Printf("Error saving message to DB: %v", err)
	} else {
		log.Printf("Message from %s to %s saved in DB", msg.FromUser, msg.ToUser)
	}
}
