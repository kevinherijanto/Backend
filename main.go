package main

import (
	"backend/database"
	"backend/routes"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
)

// ChatMessage represents the structure of a chat message
type ChatMessage struct {
	Username string    `json:"username"`
	Message  string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// Global variables to manage WebSocket connections
var (
	clients   = make(map[*websocket.Conn]bool) // Active WebSocket connections
	broadcast = make(chan ChatMessage)         // Broadcast channel for messages
	mutex     sync.Mutex                       // Mutex to handle concurrent access to clients map
)

func main() {
	// Main app for both HTTP routes and WebSocket
	app := fiber.New()

	// Enable CORS middleware for all origins
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Initialize Database
	database.ConnectDB()
	defer database.CloseDB() // Ensure the DB connection is closed on app exit

	// Enable logger middleware
	app.Use(logger.New())

	// Register HTTP routes
	routes.RegisterWalletRoutes(app)

	// WebSocket route for chat
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// Add the client to the map
		mutex.Lock()
		clients[c] = true
		mutex.Unlock()
		log.Println("WebSocket connection established")

		// Send previous messages from database to the client
		sendChatHistory(c)

		defer func() {
			// Remove the client from the map when they disconnect
			mutex.Lock()
			delete(clients, c)
			mutex.Unlock()
			c.Close()
			log.Println("WebSocket connection closed")
		}()

		// Read messages from the client
		for {
			var msg ChatMessage
			err := c.ReadJSON(&msg)
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}
			log.Printf("Received message from %s: %s\n", msg.Username, msg.Message)

			// Store the message in the database
			saveMessageToDB(msg)

			// Send the message to the broadcast channel
			broadcast <- msg
		}
	}))

	// Start a goroutine to handle broadcasting messages
	go handleBroadcast()

	// Get the PORT environment variable for Railway deployment
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Fallback if PORT isn't set
	}

	// Start the server
	log.Fatal(app.Listen(":" + port))
}

// handleBroadcast listens to the broadcast channel and sends messages to all connected clients
func handleBroadcast() {
	for {
		// Receive a message from the broadcast channel
		msg := <-broadcast

		// Send the message to all connected clients
		mutex.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("Error writing message to client:", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

// sendChatHistory sends the last 100 messages to a new WebSocket client
func sendChatHistory(c *websocket.Conn) {
	rows, err := database.DB.Query("SELECT username, message, timestamp FROM messages ORDER BY timestamp DESC LIMIT 100")
	if err != nil {
		log.Println("Error fetching chat history:", err)
		return
	}
	defer rows.Close()

	var messages []ChatMessage
	for rows.Next() {
		var msg ChatMessage
		if err := rows.Scan(&msg.Username, &msg.Message, &msg.Timestamp); err != nil {
			log.Println("Error scanning message:", err)
			return
		}
		messages = append(messages, msg)
	}

	// Send the messages to the client
	for _, msg := range messages {
		if err := c.WriteJSON(msg); err != nil {
			log.Println("Error sending message to client:", err)
			return
		}
	}
}

// saveMessageToDB stores a chat message in the database
func saveMessageToDB(msg ChatMessage) {
	_, err := database.DB.Exec("INSERT INTO messages (username, message) VALUES ($1, $2)", msg.Username, msg.Message)
	if err != nil {
		log.Println("Error saving message to database:", err)
	}
}
