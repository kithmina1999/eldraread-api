package models

import "time"

const (
	// Roles
	RoleUser  = "user"
	RoleAdmin = "admin"

	// Statuses
	StatusActive   = "active"
	StatusDisabled = "disabled"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	UID       string    `gorm:"unique;not null"` // Firebase UID
	Email     string    `gorm:"unique;not null"`
	Username  string    `gorm:"not null"`
	Role      string    `gorm:"default:user"` // user or admin
	Status    string    `gorm:"default:active"` // active or disabled
	CreatedAt time.Time
	UpdatedAt time.Time
}
