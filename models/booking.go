package models

import (
	"time"
)

type Booking struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
}
