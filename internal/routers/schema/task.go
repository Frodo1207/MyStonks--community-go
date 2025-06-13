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

type RankItem struct {
	Addr  string `json:"addr"`
	Rank  int    `json:"rank"`
	Score int    `json:"score"`
}

type RankBoards struct {
	RankItems []RankItem `json:"rank_items"`
}
