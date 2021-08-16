package queue

const MAX_SIZE int = 20

type CrawlTask interface {
	Crawl()
}
type QueueItem struct {
	f       CrawlTask
	account string
}
