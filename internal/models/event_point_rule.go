package models

import "time"

// 用于定义事件触发时奖励的积分数量及描述

type EventPointRule struct {
	ID          string    `json:"id" gorm:"primaryKey;column:id;type:varchar(64);"`                              // 规则ID，唯一
	EventID     string    `json:"event_id" gorm:"column:event_id;type:varchar(64);uniqueIndex:uniq_event_rule;"` // 关联事件ID
	Points      float64   `json:"points" gorm:"column:points;type:decimal(18,6);"`                               // 触发该事件时奖励的积分数量，支持正负
	Description string    `json:"description" gorm:"column:description;type:varchar(256);"`                      // 规则描述
	CreatedBy   string    `json:"created_by" gorm:"column:created_by;type:varchar(64);"`                         // 创建人ID
	UpdatedBy   string    `json:"updated_by" gorm:"column:updated_by;type:varchar(64);"`                         // 更新人ID
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime;"`                           // 创建时间
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime;"`                           // 更新时间
	IsDeleted   bool      `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);default:0;"`                // 是否删除（0否，1是）
}

func (*EventPointRule) TableName() string {
	return "event_point_rule"
}
