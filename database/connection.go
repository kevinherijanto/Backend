package database

import (
	"fmt"
	"log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"Wallet-Crypto-Crud/backend/models"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "root:@tcp(127.0.0.1:3306)/crypto_wallet?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate the schema
	if err := DB.AutoMigrate(&models.Wallet{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Database connected successfully")
}
