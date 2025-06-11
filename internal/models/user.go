package models

import (
	"time"
)

type User struct {
	ID          uint      `json:"id" gorm:"primaryKey;column:id;"`
	SolAddress  string    `json:"sol_address" gorm:"column:sol_address;type:varchar(128);uniqueIndex;not null"`
	Username    string    `json:"username" gorm:"column:username;type:varchar(64)"`
	TotalPoints int       `json:"total_points" gorm:"column:total_points;default:0;not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	IsDeleted   bool      `json:"is_deleted" gorm:"column:is_deleted;default:false;not null"`
}
