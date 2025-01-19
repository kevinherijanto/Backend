package main

import (
	"Wallet-Crypto-Crud/backend/database"
	"Wallet-Crypto-Crud/backend/routes"
	"log"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
)


func main() {
	// Initialize Fiber app
	app := fiber.New()

	// Enable CORS middleware for all origins
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Initialize Database
	database.ConnectDB(); 
	// Enable logger middleware
	app.Use(logger.New())

	// Register routes
	routes.RegisterWalletRoutes(app)

	// WebSocket route for notifications
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		for {
			// Simple echo server for WebSocket
			message := []byte("New wallet created!")
			if err := c.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("Error sending WebSocket message:", err)
				break
			}
		}
	}))
	// Use the PORT environment variable for Railway deployment
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // fallback if PORT isn't set
	}
	log.Fatal(app.Listen(":" + port))

	// Start server
	log.Fatal(app.Listen(":3000"))
}
