package service

import (
	"MyStonks-go/internal/models"
	"MyStonks-go/internal/store"
	"MyStonks-go/internal/taskVerifier"
	"time"
)

type TaskService struct {
	tStore store.TaskStore
	uStore store.UserStore
}

func NewTaskService(tStore store.TaskStore, uStore store.UserStore) *TaskService {
	return &TaskService{tStore: tStore, uStore: uStore}
}

func (t *TaskService) GetTasksByCategory(category, addr string) ([]models.Task, error) {
	tasks, err := t.tStore.GetTasksByCategory(category)
	if err != nil {
		return nil, err
	}
	if len(addr) < 15 {
		return tasks, nil
	}

	userTasks, err := t.tStore.GetUserCompletedTasks(addr)
	if err != nil {
		return nil, err
	}
	userTaskMap := make(map[int]models.UserTask)
	for _, ut := range userTasks {
		userTaskMap[ut.TaskID] = ut
	}
	now := time.Now()
	for i := range tasks {
		if userTask, exists := userTaskMap[tasks[i].ID]; exists {
			// 如果是每日任务，检查完成时间是否在今天
			if tasks[i].Category == "daily" {
				// 获取完成时间的年月日部分
				completedDate := userTask.CompletedAt.Truncate(24 * time.Hour)
				currentDate := now.Truncate(24 * time.Hour)

				// 如果是在今天完成的，标记为已完成
				tasks[i].Completed = completedDate.Equal(currentDate)
			} else {
				// 非每日任务，直接标记为已完成
				tasks[i].Completed = true
			}
		} else {
			// 用户未完成该任务
			tasks[i].Completed = false
		}
	}

	return tasks, nil
}

func (t *TaskService) GetUserCompleteTasks(addr string) ([]models.UserTask, error) {
	tasks, err := t.tStore.GetUserCompletedTasks(addr)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t *TaskService) CompleteTask(addr string, taskID int) error {

	task, err := t.GetTaskById(taskID)
	if err != nil {
		return err
	}

	err = t.uStore.AddPointsByAddr(addr, task.Reward)
	if err != nil {
		return err
	}

	err = t.uStore.AddCompletedTask(addr, taskID)
	if err != nil {
		return err
	}

	return nil
}

func (t *TaskService) GetTaskById(taskid int) (*models.Task, error) {
	task, err := t.tStore.GetTaskByID(taskid)
	if err != nil {
		return &models.Task{}, err
	}
	return task, nil
}

func (t *TaskService) GetUserInfoTask(userid string) (models.UserInfoTask, error) {
	user, err := t.uStore.GetUserBySolAddress(userid)
	if err != nil {
		return models.UserInfoTask{}, err
	}

	rank, err := t.uStore.GetUserRank(user)
	if err != nil {
		return models.UserInfoTask{}, err
	}

	completedTasks, err := t.tStore.GetUserCompletedTasks(userid)
	if err != nil {
		return models.UserInfoTask{}, err
	}

	return models.UserInfoTask{
		UserID:   userid,
		UserName: user.Username,
		Tasks:    completedTasks,
		Point:    user.TotalPoints,
		Rank:     int(rank),
	}, nil
}

func (t *TaskService) CheckStonksTrade(addr string) (bool, error) {
	return taskVerifier.GetInstance().VerifyStonksTradeTask(addr)
}

func (t *TaskService) RefreshDailyTask() error {
	return t.tStore.RefreshDailyTasks()
}
