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
	"github.com/golang-jwt/jwt/v4"
	_ "gorm.io/gorm"
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

var secretKey = []byte("your-secret-key")

// Claims represents the structure of the JWT claim
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func generateJWT(c *fiber.Ctx) error {
	username := c.FormValue("username")

	// Create the JWT claims, including the username and expiration time
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Token expires in 24 hours
		},
	}

	// Create a new JWT token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate token")
	}

	// Return the generated token
	return c.JSON(fiber.Map{
		"token": tokenString,
	})
}

// JWT Middleware to validate token
func jwtMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).SendString("Missing or invalid token")
		}

		// Extract the token part
		tokenString := authHeader[len("Bearer "):]

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
		}

		// Store the claims in the context
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			c.Locals("user", claims)
		}

		return c.Next()
	}
}

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

	// JWT route for login (just username, no password)
	app.Post("/login", func(c *fiber.Ctx) error {
		username := c.FormValue("username")

		// Create the JWT claims, which includes the username and expiry time
		claims := Claims{
			Username: username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Token expires in 24 hours
			},
		}

		// Create the token using your secret key
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Sign the token with your secret key
		tokenString, err := token.SignedString(secretKey)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate token")
		}

		// Return the generated token
		return c.JSON(fiber.Map{
			"token": tokenString,
		})
	})

	// WebSocket route for chat with JWT verification
	app.Get("/ws", jwtMiddleware(), websocket.New(func(c *websocket.Conn) {
		// Add the client to the map
		mutex.Lock()
		clients[c] = true
		mutex.Unlock()
		log.Println("WebSocket connection established")

		// Fetch previous chat history when a new client connects
		messages, err := fetchChatHistory()
		if err != nil {
			log.Println("Error fetching chat history:", err)
		} else {
			// Send chat history to the client
			for _, message := range messages {
				err := c.WriteJSON(message)
				if err != nil {
					log.Println("Error sending chat history:", err)
					break
				}
			}
		}

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
			// Save the message to the database
			err = saveMessageToDatabase(msg)
			if err != nil {
				log.Println("Error saving message to database:", err)
			}
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

// fetchChatHistory fetches the previous chat messages from the database
func fetchChatHistory() ([]ChatMessage, error) {
	var messages []ChatMessage
	// Use the global database connection from the database package
	db := database.DB // GORM DB instance

	// Fetch chat history from the database, ordered by creation date
	err := db.Order("created_at asc").Find(&messages).Error
	return messages, err
}

// saveMessageToDatabase saves the chat message to the database
func saveMessageToDatabase(msg ChatMessage) error {
	// Use the global database connection from the database package
	db := database.DB // GORM DB instance

	// Save the chat message to the database
	err := db.Create(&msg).Error
	return err
}
