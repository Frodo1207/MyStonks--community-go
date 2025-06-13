package store

import (
	"MyStonks-go/internal/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type UserStore interface {
	CreateUserIfNotExists(solAddress string) error
	GetUserBySolAddress(solAddress string) (*models.User, error)
	AddPointsByAddr(addr string, points int) error
	DeductPoints(user *models.User, points int, tx *gorm.DB) error
	GetUserRank(user *models.User) (int64, error)
	UpdateUsername(user *models.User, newUsername string) error
	GetLeaderboard(limit int) ([]models.User, error)
	AddCompletedTask(addr string, taskID int) error
	BindTgToUser(tg models.TelegramBinding) error
	GetUserTgInfoByAddr(addr string) (*models.TelegramBinding, error)
	GetTopsUsers() ([]models.User, error)
}

type userStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) UserStore {
	return &userStore{db: db}
}

func (s *userStore) CreateUserIfNotExists(solAddress string) error {
	var user models.User
	if err := s.db.Where(&models.User{SolAddress: solAddress}).FirstOrCreate(&user, models.User{
		SolAddress: solAddress,
		Username:   "", // 默认空用户名
	}).Error; err != nil {
		return err
	}
	return nil
}

func (s *userStore) GetUserBySolAddress(solAddress string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("sol_address = ? AND is_deleted = ?", solAddress, false).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// AddPointsByAddr 给指定SolAddress的用户增加积分
func (s *userStore) AddPointsByAddr(addr string, points int) error {
	if points <= 0 {
		return errors.New("points must be greater than 0")
	}

	// 使用事务确保一致性
	return s.db.Transaction(func(tx *gorm.DB) error {
		var user models.User

		// 查询用户
		if err := tx.Where("sol_address = ? AND is_deleted = false", addr).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("user not found")
			}
			return err
		}

		// 更新积分
		user.TotalPoints += points
		if err := tx.Model(&models.User{}).Where("id = ?", user.ID).Update("total_points", user.TotalPoints).Error; err != nil {
			return err
		}

		return nil
	})
}
func (s *userStore) DeductPoints(user *models.User, points int, tx *gorm.DB) error {
	if tx == nil {
		tx = s.db
	}

	return tx.Model(user).
		Where("total_points >= ?", points).
		Update("total_points", gorm.Expr("total_points - ?", points)).Error
}

func (s *userStore) GetUserRank(user *models.User) (int64, error) {
	var rank int64
	err := s.db.Model(&models.User{}).
		Where("total_points > ? AND is_deleted = ?", user.TotalPoints, false).
		Count(&rank).Error
	if err != nil {
		return 0, err
	}
	return rank + 1, nil // 排名从1开始
}

func (s *userStore) UpdateUsername(user *models.User, newUsername string) error {
	return s.db.Model(user).Update("username", newUsername).Error
}

func (s *userStore) GetLeaderboard(limit int) ([]models.User, error) {
	var users []models.User
	query := s.db.Where("is_deleted = ?", false).
		Order("total_points DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (s *userStore) BindTgToUser(tg models.TelegramBinding) error {
	var existing models.TelegramBinding

	// 检查该 Telegram ID 是否已绑定
	err := s.db.Where("telegram_id = ?", tg.TelegramID).First(&existing).Error
	if err == nil {
		// 已存在记录，更新绑定信息
		tg.ID = existing.ID
		return s.db.Model(&existing).Updates(tg).Error
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 不存在，创建新的绑定记录
		return s.db.Create(&tg).Error
	}

	// 其他数据库错误
	return err
}

func (s *userStore) AddCompletedTask(addr string, taskID int) error {
	// 首先查询用户ID
	var user models.User
	if err := s.db.Where("sol_address = ?", addr).First(&user).Error; err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// 创建用户任务记录
	userTask := models.UserTask{
		UserID:     uint64(user.ID),
		SolAddress: addr,
		TaskID:     taskID,
		Verified:   true,
	}

	// 使用Create或者FirstOrCreate避免重复
	return s.db.Where(models.UserTask{UserID: uint64(user.ID), TaskID: taskID}).
		FirstOrCreate(&userTask).Error
}

func (s *userStore) GetUserTgInfoByAddr(addr string) (*models.TelegramBinding, error) {
	var binding models.TelegramBinding
	err := s.db.Where("addr = ?", addr).First(&binding).Error
	if err != nil {
		return nil, err
	}
	return &binding, nil
}

func (s *userStore) GetTopsUsers() ([]models.User, error) {
	var users []models.User
	// 查询未删除的用户，按积分降序排列，限制10条记录
	err := s.db.Where("is_deleted = ?", false).
		Order("total_points DESC").
		Limit(10).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
