package store

import (
	"MyStonks-go/internal/models"
	"fmt"
	"gorm.io/gorm"
)

type TaskStore interface {
	GetTasksByCategory(category string) ([]models.Task, error)
	GetTaskByID(id int) (*models.Task, error)
	GetUserCompletedTaskIDs(addr string) (map[int]bool, error)
	IsTaskCompleted(addr string, taskID int) (bool, error)
	RecordTaskCompletion(tx *gorm.DB, addr string, taskID int) error
	GetUserCompletedTasks(addr string) ([]models.UserTask, error)
	GetUserCompletedTask(addr string) ([]models.UserTask, error)
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
	err := s.db.Where("category = ?", category).Find(&tasks).Error
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

	var user models.User
	var userTasks []models.UserTask

	// 2. 查询该用户的任务完成记录
	if err := s.db.Where("sol_address = ? AND is_deleted = false", user.ID).Find(&userTasks).Error; err != nil {
		return nil, err
	}

	return userTasks, nil
}
