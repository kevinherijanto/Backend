package main

import (
	"backend/database"
	"backend/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
)

// Global variable to store connected WebSocket clients
var clients = make(map[*websocket.Conn]bool)

func main() {
	// Main app for HTTP routes (app1)
	httpApp := fiber.New()

	// Enable CORS middleware for all origins
	httpApp.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Initialize Database
	database.ConnectDB()

	// Enable logger middleware
	httpApp.Use(logger.New())

	// Register routes
	routes.RegisterWalletRoutes(httpApp)

	// WebSocket server (app2)
	wsApp := fiber.New()

	// WebSocket route for notifications
	wsApp.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// Add the client to the clients map
		clients[c] = true
		log.Println("New WebSocket connection established")

		// Keep the connection alive and listen for any errors or disconnections
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				log.Println("WebSocket error:", err)
				delete(clients, c) // Remove client from the list on disconnect
				break
			}
		}
	}))

	// Function to send a notification to all connected clients
	sendWalletNotification := func(message string) {
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println("Error sending message to WebSocket client:", err)
				client.Close()
				delete(clients, client) // Remove client if there's an error
			}
		}
	}

	// Simulate wallet creation (for demonstration purposes)
	httpApp.Post("/create-wallet", func(c *fiber.Ctx) error {
		// Simulate wallet creation process here

		// After wallet is created, send a notification to all WebSocket clients
		sendWalletNotification("New wallet created!")

		// Respond to HTTP request
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Wallet created and notification sent",
		})
	})

	// Get the PORT environment variable for Railway deployment
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // fallback if PORT isn't set
	}

	// Start HTTP server
	go func() {
		log.Fatal(httpApp.Listen(":" + port)) // Port for HTTP routes
	}()

	// Start WebSocket server on a different port (e.g., 8081)
	go func() {
		log.Fatal(wsApp.Listen(":" + "8081")) // Separate port for WebSocket
	}()

	// Block to keep both servers running
	select {}
}
