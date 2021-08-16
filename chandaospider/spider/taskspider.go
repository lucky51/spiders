package spider

import (
	"bytes"
	"chandaospider/model"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"go.uber.org/zap"
)

const (
	mytaskFinishedByUrl  = "http://192.168.0.178:8011/my-task-finishedBy.html"
	taskViewPageTemplate = "http://192.168.0.178:8011/task-view-%s.html"
)

type TaskSpider struct {
	*SessionSpider
	TaskItemWriter chan *model.TaskItem
	TaskItems      []*model.TaskItem
	Filter         func(*model.TaskItem) bool
}

// UrlCollector 收集完成task的URL地址
func (s *TaskSpider) UrlCollector(cancel context.CancelFunc) {
	defer cancel()
	waited := make(chan struct{})
	defer close(waited)
	taskViewCollector := s.Collector.Clone()
	next := true
	var taskviewItem *model.TaskItem

	taskViewCollector.OnResponse(func(resp *colly.Response) {
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body))
		if err != nil {

		} else {
			findStr := "<strong>" + s.username + "</strong> 完成"
			doc.Find("ol.histories-list>li").Each(func(i int, sec *goquery.Selection) {
				raw, err := sec.Html()
				if err != nil {
					zap.L().Error("task get histories list error", zap.Error(err))
					return
				}
				if strings.Contains(raw, findStr) {
					taskviewItem.FinishedTime = finishedTimeReg.FindString(raw)
				}
			})
			if s.Filter != nil && s.Filter(taskviewItem) {
				s.TaskItems = append(s.TaskItems, taskviewItem)
			} else {
				// 不能跳过 ，列表没有完成时间排序
				//next = false
				zap.L().Info("skip task", zap.String("ID", taskviewItem.ID))
			}
		}
	})
	go func() {
	end:
		for {
			select {
			case taskviewItem = <-s.TaskItemWriter:
			case <-time.After(time.Second * 2):
				next = false
			}
			if next {
				zap.L().Info("visiting task", zap.String("ID", taskviewItem.ID), zap.String("TaskUrl", taskviewItem.TaskViewUrl))
				taskViewCollector.Visit(taskviewItem.TaskViewUrl)
			} else {
				break end
			}
		}
		waited <- struct{}{}
	}()
	<-waited
}

// Close 关闭开启的channel
func (s *TaskSpider) Close() {
	close(s.TaskItemWriter)
}

// Crawl 开始爬取
func (s *TaskSpider) Crawl() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	sessionCollector, err := s.login()
	if err != nil {
		zap.L().Error("login error:", zap.Error(err))
		cancel()
		return err
	}
	go func(c *colly.Collector) {
		collector := c.Clone()
		// collector.Async = true
		// collector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})
		collector.OnResponse(func(r *colly.Response) {
			doc, err := goquery.NewDocumentFromReader(bytes.NewReader(r.Body))
			if err != nil {
				zap.L().Info("create new goquery document error", zap.Error(err))
				return
			}
			doc.Find("table#tasktable>tbody>tr").EachWithBreak(func(_ int, selector *goquery.Selection) bool {
				i := 0
				for _, n := range selector.Nodes {

					h := colly.NewHTMLElementFromSelectionNode(r, selector, n, i)
					i++
					id := h.ChildText("td.c-id")
					cpri := h.ChildText("td.c-pri>span")
					project := h.ChildText("td.c-project>a")
					taskName := h.ChildText("td.c-name>a")
					createdBy := h.ChildText("td:nth-child(5)")
					resolvedBy := h.ChildText("td:nth-child(7)")
					taskviewpage := fmt.Sprintf(taskViewPageTemplate, id)
				end:
					for {
						select {
						case s.TaskItemWriter <- &model.TaskItem{
							ID:          id,
							PRI:         cpri,
							Project:     project,
							TaskName:    taskName,
							CreatedBy:   createdBy,
							FinishedBy:  resolvedBy,
							TaskViewUrl: taskviewpage,
						}:
							zap.L().Info("join task", zap.String("ID", id), zap.String("TaskPage", taskviewpage))
							break end
						case <-ctx.Done():
							zap.L().Info("task crawl is done.")
							return false
						}
					}
				}
				return true
			})
		})
		collector.Visit(mytaskFinishedByUrl)
	}(sessionCollector)
	zap.L().Info("Starting Task URL collect..")
	s.UrlCollector(cancel)
	return nil
}

// NewTaskSpider 构造爬取任务
func NewTaskSpider(session *SessionSpider, f func(item *model.TaskItem) bool) *TaskSpider {
	return &TaskSpider{
		SessionSpider:  session,
		TaskItemWriter: make(chan *model.TaskItem),
		TaskItems:      []*model.TaskItem{},
		Filter:         f,
	}
}
