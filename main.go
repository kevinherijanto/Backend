package main

import (
	"backend/database"
	"backend/routes"
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
)

// ChatMessage represents the structure of a chat message
type ChatMessage struct {
	Username string `json:"username"`
	Message  string `json:"message"`
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
