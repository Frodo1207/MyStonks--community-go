package store

import (
	"MyStonks-go/internal/models"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type TaskStore interface {
	GetTasksByCategory(category string) ([]models.Task, error)
	GetTaskByID(id int) (*models.Task, error)
	GetUserCompletedTaskIDs(addr string) (map[int]bool, error)
	IsTaskCompleted(addr string, taskID int) (bool, error)
	RecordTaskCompletion(tx *gorm.DB, addr string, taskID int) error
	GetUserCompletedTasks(addr string) ([]models.UserTask, error)
	GetUserCompletedTask(addr string) ([]models.UserTask, error)
	RefreshDailyTasks() error
}

type taskStore struct {
	db *gorm.DB
}

func NewTaskStore(db *gorm.DB) TaskStore {
	return &taskStore{
		db: db,
	}
}

func (s *taskStore) GetTasksByCategory(category string) ([]models.Task, error) {
	var tasks []models.Task
	err := s.db.Where("category = ? AND is_deleted = ?", category, false).Find(&tasks).Error
	return tasks, err
}

func (s *taskStore) GetTaskByID(id int) (*models.Task, error) {
	var task models.Task
	err := s.db.Where("id = ?", id).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *taskStore) GetUserCompletedTaskIDs(addr string) (map[int]bool, error) {
	var userTasks []models.UserTask
	err := s.db.Where("sol_address = ?", addr).Find(&userTasks).Error
	if err != nil {
		return nil, err
	}
	fmt.Println("userTasks:", userTasks)
	result := make(map[int]bool)
	for _, ut := range userTasks {
		result[ut.TaskID] = true
	}

	return result, nil
}

func (s *taskStore) GetUserCompletedTask(addr string) ([]models.UserTask, error) {
	var userTasks []models.UserTask
	err := s.db.Where("sol_address = ?", addr).Find(&userTasks).Error
	if err != nil {
		return nil, err
	}

	return userTasks, nil
}

func (s *taskStore) IsTaskCompleted(addr string, taskID int) (bool, error) {
	var count int64
	err := s.db.Model(&models.UserTask{}).
		Where("user_id = ? AND task_id = ?", addr, taskID).
		Count(&count).Error
	return count > 0, err
}

func (s *taskStore) RecordTaskCompletion(tx *gorm.DB, addr string, taskID int) error {
	return nil
}

func (s *taskStore) GetUserCompletedTasks(addr string) ([]models.UserTask, error) {
	var userTasks []models.UserTask
	if err := s.db.Where("sol_address = ? AND is_deleted = false", addr).Find(&userTasks).Error; err != nil {
		return nil, err
	}
	fmt.Println("userTasks:", userTasks)
	return userTasks, nil
}

func (s *taskStore) RefreshDailyTasks() error {
	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 查询所有未删除的daily任务
	var existingTasks []models.Task
	if err := tx.Where("category = ? AND is_deleted = ?", "daily", false).Find(&existingTasks).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to query existing daily tasks: %v", err)
	}

	if len(existingTasks) == 0 {
		tx.Rollback()
		return nil // 没有需要刷新的任务
	}

	// 2. 生成新的ID并复制任务创建新记录
	newTasks := make([]models.Task, len(existingTasks))
	now := time.Now()
	userID := "system" // 假设系统用户执行此操作

	// 获取当前最大ID
	var maxID int
	if err := tx.Model(&models.Task{}).Select("COALESCE(MAX(id), 0)").Scan(&maxID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get max task ID: %v", err)
	}

	for i, task := range existingTasks {
		maxID++ // 递增ID
		newTasks[i] = models.Task{
			ID:            maxID, // 显式设置新ID
			Step:          task.Step,
			Title:         task.Title,
			Description:   task.Description,
			Reward:        task.Reward,
			Category:      task.Category,
			Icon:          task.Icon,
			SpecialAction: task.SpecialAction,
			CreatedBy:     userID,
			UpdatedBy:     userID,
			CreatedAt:     now,
			UpdatedAt:     now,
			IsDeleted:     false,
		}
	}

	// 批量创建新任务
	if err := tx.Create(&newTasks).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create new daily tasks: %v", err)
	}

	// 3. 标记旧任务为已删除
	oldTaskIDs := make([]int, len(existingTasks))
	for i, task := range existingTasks {
		oldTaskIDs[i] = task.ID
	}

	if err := tx.Model(&models.Task{}).
		Where("id IN (?)", oldTaskIDs).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"updated_by": userID,
			"updated_at": now,
		}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to mark old tasks as deleted: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
