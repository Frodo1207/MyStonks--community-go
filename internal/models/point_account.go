package models

import "time"

// 用于记录用户积分账户信息，包括余额、创建/更新人、时间等
// 如需精确小数建议用 github.com/shopspring/decimal，这里先用 float64

type PointAccount struct {
	ID        string    `json:"id" gorm:"primaryKey;column:id;type:varchar(64);"`                  // 积分账户ID，唯一
	UserID    string    `json:"user_id" gorm:"column:user_id;type:varchar(64);index:idx_user_id;"` // 关联用户ID
	Balance   float64   `json:"balance" gorm:"column:balance;type:decimal(18,6);default:0;"`       // 积分余额
	CreatedBy string    `json:"created_by" gorm:"column:created_by;type:varchar(64);"`             // 创建人ID
	UpdatedBy string    `json:"updated_by" gorm:"column:updated_by;type:varchar(64);"`             // 更新人ID
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime;"`               // 创建时间
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime;"`               // 更新时间
	IsDeleted bool      `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);default:0;"`    // 是否删除（0否，1是）
}

func (*PointAccount) TableName() string {
	return "point_account"
}
