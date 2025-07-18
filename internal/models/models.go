package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"unique;not null" json:"username"`
	Email     string         `gorm:"unique;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Role      string         `gorm:"default:'user'" json:"role"` // 'user' or 'admin'
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Movie struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	Genre       string         `json:"genre"`
	Director    string         `json:"director"`
	ReleaseYear int            `json:"release_year"`
	Duration    int            `json:"duration"` // in minutes
	Rating      float64        `json:"rating"`
	FilePath    string         `gorm:"not null" json:"file_path"`
	ThumbnailPath string       `json:"thumbnail_path"`
	FileSize    int64          `json:"file_size"`
	UploadedBy  uint           `json:"uploaded_by"`
	User        User           `gorm:"foreignKey:UploadedBy" json:"user"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserSession struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Token     string    `gorm:"not null" json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
}

type ViewHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	MovieID   uint      `gorm:"not null" json:"movie_id"`
	WatchedAt time.Time `json:"watched_at"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	Movie     Movie     `gorm:"foreignKey:MovieID" json:"movie"`
}
