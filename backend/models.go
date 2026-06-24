package main

import "time"

type Card struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UID       string    `gorm:"uniqueIndex;not null" json:"uid"`
	Label     string    `gorm:"not null" json:"label"`
	CreatedAt time.Time `json:"created_at"`
}

type AccessLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UID       string    `gorm:"index;not null" json:"uid"`
	Label     string    `json:"label"`
	Granted   bool      `json:"granted"`
	CreatedAt time.Time `json:"created_at"`
}
