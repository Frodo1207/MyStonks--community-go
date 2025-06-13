package models

import "time"

type TelegramBinding struct {
	ID         uint64 `gorm:"primaryKey"`
	FirstName  string
	Addr       string `gorm:"not null;uniqueIndex"`
	TelegramID int64  `gorm:"not null;uniqueIndex"`
	Username   string
	PhotoURL   string
	AuthDate   int64
	Hash       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
