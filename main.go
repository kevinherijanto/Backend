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
