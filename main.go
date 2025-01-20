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

	// WebSocket route for notifications
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		log.Println("WebSocket connection established")
		defer c.Close()

		for {
			// Read message from the client
			messageType, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}
			log.Printf("Received: %s\n", msg)

			// Echo message back to the client
			if err := c.WriteMessage(messageType, msg); err != nil {
				log.Println("Write error:", err)
				break
			}
		}
	}))

	// Get the PORT environment variable for Railway deployment
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Fallback if PORT isn't set
	}

	// Start the unified server on the same port
	log.Fatal(app.Listen(":" + port))
}
