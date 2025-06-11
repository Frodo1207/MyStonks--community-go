package schema

type Task struct {
	ID            int    `json:"id"`
	Step          int    `json:"step,omitempty"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Reward        int    `json:"reward"`
	Completed     bool   `json:"completed"` // 是否完成
	Category      string `json:"category"`  // 任务分类
	Icon          string `json:"icon,omitempty"`
	SpecialAction string `json:"special_action,omitempty"`
}
