package service

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Reward      int    `json:"reward"`
	Completed   bool   `json:"completed"`
	Category    string `json:"category,omitempty"`
	Step        int    `json:"step,omitempty"`     // 仅新人任务用
	Icon        string `json:"icon,omitempty"`     // 仅日常任务用
	Progress    int    `json:"progress,omitempty"` // 部分任务可能用
}

type UserInfoTask struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Tasks    []Task `json:"tasks"`
	Point    int    `json:"point"`
	Rank     int    `json:"rank"`
}

func GetTasksByCategory(category string) ([]Task, error) {
	// 所有 mock 任务
	allTasks := []Task{
		// 新手任务
		{ID: 1, Step: 1, Title: "创建钱包", Description: "创建你的第一个Web3钱包", Reward: 200, Completed: true, Category: "newbie"},
		{ID: 2, Step: 2, Title: "首次登录", Description: "使用钱包连接登录平台", Reward: 100, Completed: true, Category: "newbie"},
		{ID: 3, Step: 3, Title: "完善资料", Description: "设置你的用户名和头像", Reward: 150, Completed: false, Category: "newbie"},
		{ID: 4, Step: 4, Title: "首次交互", Description: "在社区发表你的第一条动态", Reward: 200, Completed: false, Category: "newbie"},
		{ID: 5, Step: 5, Title: "邀请好友", Description: "邀请至少1位好友加入社区", Reward: 300, Completed: false, Category: "newbie"},

		// 日常任务
		{ID: 101, Title: "每日签到", Description: "访问社区网站并签到", Reward: 50, Completed: true, Category: "daily", Icon: "check-circle"},
		{ID: 102, Title: "社区互动", Description: "在论坛发表1条评论", Reward: 100, Completed: false, Category: "daily", Icon: "message-square"},
		{ID: 103, Title: "学习Web3", Description: "完成一篇教程学习", Reward: 150, Completed: false, Category: "daily", Icon: "book-open"},

		// 其他任务
		{ID: 201, Title: "优质内容创作", Description: "发布一篇被管理员标记为优质的内容", Reward: 500, Completed: false, Category: "other"},
		{ID: 202, Title: "完成Web3课程", Description: "完成平台提供的Web3基础课程学习", Reward: 800, Completed: false, Category: "other", Progress: 35},
		{ID: 203, Title: "参与AMA活动", Description: "参加一次社区AMA并提问", Reward: 300, Completed: false, Category: "other"},
		{ID: 204, Title: "邀请3位好友", Description: "成功邀请3位好友加入并活跃", Reward: 1000, Completed: false, Category: "other"},
	}

	// 分类过滤
	var filtered []Task
	for _, task := range allTasks {
		if task.Category == category {
			filtered = append(filtered, task)
		}
	}

	return filtered, nil
}

func IsTaskComplete(userid string, taskid string) (bool, error) {
	// 这里仅为示例，实际应用中应从数据库查询用户完成任务情况
	return true, nil
}

func GetUserInfoTask(userid string) (UserInfoTask, error) {
	return UserInfoTask{
		UserID:   userid,
		UserName: "test_user",
		Tasks:    []Task{},
		Point:    100,
		Rank:     1,
	}, nil
}
