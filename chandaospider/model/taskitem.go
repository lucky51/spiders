package model

// TaskItem 爬取任务的模型
type TaskItem struct {
	ID string
	// 优先级
	PRI          string
	Project      string
	TaskName     string
	CreatedBy    string
	FinishedBy   string
	FinishedTime string
	TaskViewUrl  string
}
