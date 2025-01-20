package models

import (
	"time"
	"gorm.io/gorm"
)

type ChatMessage struct {
	gorm.Model
	Username  string    `json:"username"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
