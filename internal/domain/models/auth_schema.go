package models

import "time"

// User represents the users table in the database
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Password  string    `gorm:"not null" json:"-"` // json:"-" ensures the password hash is never serialized to JSON
	Accounts  []Account `json:"accounts,omitempty"`
	Sessions  []Session `json:"sessions,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Account represents an OAuth provider account (e.g., Google, GitHub) linked to a user.
type Account struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	UserID            uint      `gorm:"not null;index" json:"user_id"`
	User              User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Provider          string    `gorm:"not null;index" json:"provider"`
	ProviderAccountID string    `gorm:"not null;index" json:"provider_account_id"`
	RefreshToken      *string   `json:"-"`
	AccessToken       *string   `json:"-"`
	ExpiresAt         *int64    `json:"expires_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Session represents an active user session in the database.
type Session struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	SessionToken string    `gorm:"uniqueIndex;not null" json:"-"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	User         User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	ExpiresAt    time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
