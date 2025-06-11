package models

import (
	"time"
)

type Task struct {
	ID            int       `gorm:"primaryKey;column:id" json:"id"`
	Step          int       `gorm:"column:step" json:"step"`
	Title         string    `gorm:"column:title;type:varchar(128);not null" json:"title"`
	Description   string    `gorm:"column:description;type:text" json:"description"`
	Reward        int       `gorm:"column:reward;not null;default:0" json:"reward"`
	Completed     bool      `gorm:"-" json:"completed"` // 不存储到数据库，仅用于业务逻辑
	Category      string    `gorm:"column:category;type:varchar(64);index" json:"category"`
	Icon          string    `gorm:"column:icon;type:varchar(64)" json:"icon"`
	SpecialAction string    `gorm:"column:special_action;type:varchar(64)" json:"special_action"`
	CreatedBy     string    `gorm:"column:created_by;type:varchar(64);not null" json:"created_by"`
	UpdatedBy     string    `gorm:"column:updated_by;type:varchar(64);not null" json:"updated_by"`
	CreatedAt     time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	IsDeleted     bool      `gorm:"column:is_deleted;not null;default:0" json:"is_deleted"`
}

// UserTask 用户任务完成记录表模型
type UserTask struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`                           // 自增主键
	UserID      uint64    `gorm:"not null;column:user_id;index:uk_user_task,unique" json:"user_id"`       // 用户ID
	SolAddress  string    `gorm:"not null;column:sol_address;type:varchar(128);index" json:"sol_address"` // Solana钱包地址
	TaskID      int       `gorm:"not null;column:task_id;index:uk_user_task,unique" json:"task_id"`       // 任务ID
	CompletedAt time.Time `gorm:"column:completed_at;default:CURRENT_TIMESTAMP" json:"completed_at"`      // 完成时间
	Verified    bool      `gorm:"column:verified;default:false" json:"verified"`                          // 是否已验证

	CreatedBy string    `gorm:"not null;column:created_by;type:varchar(64)" json:"created_by"` // 创建人ID
	UpdatedBy string    `gorm:"not null;column:updated_by;type:varchar(64)" json:"updated_by"` // 更新人ID
	CreatedAt time.Time `gorm:"not null;column:created_at;autoCreateTime" json:"created_at"`   // 创建时间
	UpdatedAt time.Time `gorm:"not null;column:updated_at;autoUpdateTime" json:"updated_at"`   // 更新时间
	IsDeleted bool      `gorm:"not null;column:is_deleted;default:false" json:"is_deleted"`    // 是否删除（0否，1是）
}

// TableName 指定表名
func (UserTask) TableName() string {
	return "user_tasks"
}

// UserInfoTask 用户任务信息DTO
type UserInfoTask struct {
	UserID   string     `json:"user_id"`
	UserName string     `json:"user_name"`
	Tasks    []UserTask `json:"tasks"`
	Point    int        `json:"point"`
	Rank     int        `json:"rank"`
}
