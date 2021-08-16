package model

// BugItem 爬取 bug 模型
type BugItem struct {
	ID          string
	Level       string
	BugType     string
	BugTitle    string
	CreatedBy   string
	ResolveBy   string
	Resolution  string
	ResolveTime string
	BugViewUrl  string
}
