package models

import "gorm.io/gorm"

type Announcement struct {
    gorm.Model
    Title   string `json:"title"`
    Content string `json:"content"`
}
// TableName overrides the default table name used by GORM
func (Announcement) TableName() string {
    return "announcements"  // Ensure the table name is "announcements" (plural)
}