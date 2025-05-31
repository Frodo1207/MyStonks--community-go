package models

import "time"

type User struct {
	ID         uint      `json:"id" gorm:"primaryKey;column:id;"`
	SolAddress string    `json:"sol_address" gorm:"column:sol_address;"`
	Username   string    `json:"username" gorm:"column:username;"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at;"`
	IsDeleted  bool      `json:"is_deleted" gorm:"column:is_deleted;"`
}

func (*User) TableName() string {
	return "users"
}

func CreateUserIfNotExists(solAddress string) error {
	var user User
	if err := db.Where(&User{SolAddress: solAddress}).FirstOrCreate(&user).Error; err != nil {
		return err
	}
	return nil
}

func GetUserBySolAddress(solAddress string) (*User, error) {
	var user User
	if err := db.Where("sol_address = ? and is_deleted = ?", solAddress, false).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
