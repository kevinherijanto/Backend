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
	// Main app for HTTP routes
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

	// Register routes
	routes.RegisterWalletRoutes(app)

	// Use the PORT environment variable for Railway deployment
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // fallback if PORT isn't set
	}

	// Start the main Fiber app
	go func() {
		log.Fatal(app.Listen(":" + port))
	}()

	// WebSocket server on port 8080
	wsApp := fiber.New()

	// WebSocket route for notifications
	wsApp.Get("/ws", websocket.New(func(c *websocket.Conn) {
		for {
			// Simple echo server for WebSocket
			message := []byte("New wallet created!")
			if err := c.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("Error sending WebSocket message:", err)
				break
			}
		}
	}))

	// Start WebSocket server on port 8080
	log.Fatal(wsApp.Listen(":8080"))
}
