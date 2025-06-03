package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// 用于记录用户触发的事件及相关元数据
// metadata 字段为 JSON 类型，建议用 map[string]interface{} 或自定义结构体

type EventRecord struct {
	ID        string    `json:"id" gorm:"primaryKey;column:id;type:varchar(64);"`                     // 事件记录ID，唯一
	EventID   string    `json:"event_id" gorm:"column:event_id;type:varchar(64);index:idx_event_id;"` // 事件ID
	UserID    string    `json:"user_id" gorm:"column:user_id;type:varchar(64);index:idx_user_id;"`    // 触发事件的用户ID
	Metadata  JSONMap   `json:"metadata" gorm:"column:metadata;type:json;"`                           // 事件具体数据
	EventTime time.Time `json:"event_time" gorm:"column:event_time;autoCreateTime;"`                  // 事件发生时间
	CreatedBy string    `json:"created_by" gorm:"column:created_by;type:varchar(64);"`                // 创建人ID
	UpdatedBy string    `json:"updated_by" gorm:"column:updated_by;type:varchar(64);"`                // 更新人ID
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime;"`                  // 创建时间
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime;"`                  // 更新时间
	IsDeleted bool      `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);default:0;"`       // 是否删除（0否，1是）
}

// TableName 设置表名为 event_record
func (*EventRecord) TableName() string {
	return "event_record"
}

// JSONMap 用于处理 JSON 类型字段
// 你可以根据实际业务自定义结构体替换

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
