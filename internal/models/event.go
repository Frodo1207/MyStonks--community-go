package models

import "time"

// 用于记录系统支持的事件类型及描述

type Event struct {
	ID          string    `json:"id" gorm:"primaryKey;column:id;type:varchar(64);"`               // 事件ID，如 TASK_COMPLETE、ON_CHAIN_TRANSFER 等
	Name        string    `json:"name" gorm:"column:name;type:varchar(128);"`                     // 事件名称
	Description string    `json:"description" gorm:"column:description;type:text;"`               // 事件描述
	CreatedBy   string    `json:"created_by" gorm:"column:created_by;type:varchar(64);"`          // 创建人ID
	UpdatedBy   string    `json:"updated_by" gorm:"column:updated_by;type:varchar(64);"`          // 更新人ID
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime;"`            // 创建时间
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime;"`            // 更新时间
	IsDeleted   bool      `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);default:0;"` // 是否删除（0否，1是）
}

func (*Event) TableName() string {
	return "event"
}
