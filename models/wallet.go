package models

import "gorm.io/gorm"

type Wallet struct {
    ID       uint   `json:"id" gorm:"primaryKey"`  // Add ID as the primary key
    Username string `json:"username"`             // User associated with the wallet
    Address  string `json:"address"`              // Wallet address
    Balance  float64 `json:"balance"`             // Wallet balance
    Currency string `json:"currency"`             // Wallet currency
    gorm.Model
}
