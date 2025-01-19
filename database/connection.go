package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Ganti string ini dengan string koneksi yang sudah diformat
	dsn := "root:TKpUAYFvcQXHfmtiApbvAQTXsPzcGvoO@tcp(mysql.railway.internal:3306)/railway?charset=utf8mb4&parseTime=True&loc=Local"

	// Koneksi ke database menggunakan GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	log.Println("Database connected successfully!")
}
