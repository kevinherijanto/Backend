package models

import "gorm.io/gorm"

type Announcement struct {
    ID      uint   `json:"id" gorm:"primaryKey"`  // Primary key
    Title   string `json:"title"`                // Judul pengumuman
    Content string `json:"content"`              // Isi pengumuman
    gorm.Model                                   // Timestamp fields (CreatedAt, UpdatedAt, DeletedAt)
}
