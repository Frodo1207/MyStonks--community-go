package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type EventType string

const (
	EventTypeTaskComplete    EventType = "TASK_COMPLETE"
	EventTypeOnChainTransfer EventType = "ON_CHAIN_TRANSFER"
)

type Event struct {
	ID        uint      `json:"id" gorm:"primaryKey;column:id;"`
	UserID    uint      `json:"user_id" gorm:"column:user_id;"`
	EventType EventType `json:"event_type" gorm:"column:event_type;"`
	Metadata  JSONMap   `json:"metadata" gorm:"column:metadata;"`     // 事件具体数据
	EventTime time.Time `json:"event_time" gorm:"column:event_time;"` // 事件发生时间
	CreatedBy uint      `json:"created_by" gorm:"column:created_by;"` // 创建人ID
	UpdatedBy uint      `json:"updated_by" gorm:"column:updated_by"`  // 更新人ID
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;"` // 创建时间
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;"` // 更新时间
	IsDeleted bool      `json:"is_deleted" gorm:"column:is_deleted;"` // 是否删除（0否，1是）
}

func (*Event) TableName() string {
	return "event"
}

type JSONMap map[string]interface{}

func (m JSONMap) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*m = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, m)
}
