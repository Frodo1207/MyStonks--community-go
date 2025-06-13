package models

import "time"

type Events struct {
	ID               string    `json:"id" gorm:"primaryKey;column:id;type:varchar(64);"`                      // 事件ID
	Title            string    `json:"title" gorm:"column:title;type:varchar(128);"`                          // 活动标题
	Name             string    `json:"name" gorm:"column:name;type:varchar(128);"`                            // 活动名称(内部标识，英文或拼音)
	Description      string    `json:"description" gorm:"column:description;type:text;"`                      // 活动详细描述
	Content          string    `json:"content" gorm:"column:content;type:longtext;"`                          // 活动富文本内容
	CoverImageURL    string    `json:"cover_image_url" gorm:"column:cover_image_url;type:varchar(255);"`      // 封面图片链接
	Location         string    `json:"location" gorm:"column:location;type:varchar(128);"`                    // 活动地点(如"Twitter Space"、"Telegram群"、"北京798艺术区")
	StartTime        time.Time `json:"start_time" gorm:"column:start_time;"`                                  // 活动开始时间
	EndTime          time.Time `json:"end_time" gorm:"column:end_time;"`                                      // 活动结束时间
	RegisterURL      string    `json:"register_url" gorm:"column:register_url;type:varchar(255);"`            // 参与/报名链接
	MaxAttendees     int       `json:"max_attendees" gorm:"column:max_attendees;type:int;"`                   // 最大参与人数(0表示不限)
	CurrentAttendees int       `json:"current_attendees" gorm:"column:current_attendees;type:int;default:0;"` // 当前参与人数
	Popularity       int       `json:"popularity" gorm:"column:popularity;type:int;default:0;"`               // 热度(可根据浏览量、参与数等计算)
	ViewCount        int       `json:"view_count" gorm:"column:view_count;type:int;default:0;"`               // 浏览量
	CategoryID       string    `json:"category_id" gorm:"column:category_id;type:varchar(64);"`               // 分类ID
	LocationType     string    `json:"location_type" gorm:"column:location_type;type:varchar(32);"`           // 地点类型(online/offline)
	Status           int       `json:"status" gorm:"column:status;type:tinyint;default:0;"`                   // 状态(0:未开始,1:进行中,2:已结束,3:已取消)
	IsExpired        bool      `json:"is_expired" gorm:"column:is_expired;type:tinyint(1);default:0;"`        // 是否过期
	IsFeatured       bool      `json:"is_featured" gorm:"column:is_featured;type:tinyint(1);default:0;"`      // 是否推荐
	IsDeleted        bool      `json:"is_deleted" gorm:"column:is_deleted;type:tinyint(1);default:0;"`        // 是否删除（0否，1是）
	CreatedBy        string    `json:"created_by" gorm:"column:created_by;type:varchar(64);"`                 // 创建人ID
	UpdatedBy        string    `json:"updated_by" gorm:"column:updated_by;type:varchar(64);"`                 // 更新人ID
	CreatedAt        time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime;"`                   // 创建时间
	UpdatedAt        time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime;"`                   // 更新时间
}

func (*Events) TableName() string {
	return "events"
}
