// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"time"
)

type User struct {
	ID            string    `json:"id"`
	Role          string    `json:"role"`
	Phone         string    `json:"phone"`
	Name          string    `json:"name"`
	ImageUrl      string    `json:"image_url"`
	StatusText    string    `json:"status_text"`
	IsBlocked     bool      `json:"is_blocked"`
	BlockedReason string    `json:"blocked_reason"`
	CreatedAt     time.Time `json:"created_at"`
}

type Vehicle struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	Brand     string    `json:"brand"`
	Number    string    `json:"number"`
	Region    int32     `json:"region"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
}
