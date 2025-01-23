package routes

import (
	"github.com/gofiber/fiber/v2"
	"backend/database"
	"backend/models"
)

func RegisterWalletRoutes(app *fiber.App) {
	app.Post("/wallets", createWallet)
	app.Get("/wallets", getWallets)
	app.Get("/wallets/username/:username", getWalletsByUsername)
	app.Put("/wallets/:id", updateWallet)
	app.Delete("/wallets/:id", deleteWallet)

	app.Get("/api/chat-history", getChatHistory)
	app.Post("/announcements", createAnnouncement)
	app.Get("/announcements", getAnnouncements)
}

func createWallet(c *fiber.Ctx) error {
	wallet := new(models.Wallet)
	if err := c.BodyParser(wallet); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	database.DB.Create(&wallet)
	return c.Status(201).JSON(wallet)
}

func getWallets(c *fiber.Ctx) error {
	var wallets []models.Wallet
	database.DB.Find(&wallets)
	return c.JSON(wallets)
}


func getWalletsByUsername(c *fiber.Ctx) error {
	username := c.Params("username")
	var wallets []models.Wallet

	if err := database.DB.Where("username = ?", username).Find(&wallets).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "No wallets found for this username"})
	}

	return c.JSON(wallets)
}

func updateWallet(c *fiber.Ctx) error {
	id := c.Params("id")
	var wallet models.Wallet

	if err := database.DB.First(&wallet, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Wallet not found"})
	}

	var updatedData models.Wallet
	if err := c.BodyParser(&updatedData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	wallet.Address = updatedData.Address
	wallet.Balance = updatedData.Balance
	wallet.Currency = updatedData.Currency

	if err := database.DB.Save(&wallet).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update wallet"})
	}

	return c.JSON(wallet)
}

func deleteWallet(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := database.DB.Delete(&models.Wallet{}, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Wallet not found"})
	}
	return c.SendStatus(204)
}

func getChatHistory(c *fiber.Ctx) error {
	var chatHistory []models.ChatMessage

	// Retrieve the chat history from the chat_messages table, ordered by timestamp
	if err := database.DB.Order("timestamp ASC").Find(&chatHistory).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve chat history"})
	}

	return c.JSON(chatHistory)
}
func createAnnouncement(c *fiber.Ctx) error {
    announcement := new(models.Announcement)

    // Parse JSON body
    if err := c.BodyParser(announcement); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Failed to parse request body. Ensure the body is in valid JSON format.",
        })
    }

    // Save to database
    if err := database.DB.Create(&announcement).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to save announcement to the database.",
        })
    }

    // Return the created announcement
    return c.Status(fiber.StatusCreated).JSON(announcement)
}

func getAnnouncements(c *fiber.Ctx) error {
	var announcements []models.Announcement

	// Fetch all announcements from the database
	if err := database.DB.Find(&announcements).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve announcements from the database.",
		})
	}

	if len(announcements) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No announcements found.",
		})
	}

	return c.JSON(announcements)
}