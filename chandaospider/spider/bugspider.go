package spider

import (
	"bytes"
	"chandaospider/model"
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"go.uber.org/zap"
)

const (
	mybugResolveByUrl   string = "http://192.168.0.178:8011/my-bug-resolvedBy.html"
	bugViewPageTemplate string = "http://192.168.0.178:8011/bug-view-%s.html"
)
const (
	finishedTimeRegString string = "\\d{4}-\\d{2}-\\d{2}\\s+\\d{2}:\\d{2}:\\d{2}"
)

var finishedTimeReg = regexp.MustCompile(finishedTimeRegString)

type BugSpider struct {
	*SessionSpider
	BugItemWriter chan *model.BugItem
	BugItems      []*model.BugItem
	Filter        func(*model.BugItem) bool
	PagerMyBug    int
}

var reg = regexp.MustCompile(finishedTimeRegString)

func ValidateTimeStr(s string) bool {
	isvalid := reg.Match([]byte(s))
	return isvalid
}

func (s *BugSpider) UrlCollector(cancel context.CancelFunc) {
	waited := make(chan struct{})
	defer cancel()
	defer close(waited)
	bugViewCollector := s.Collector.Clone()
	// bugViewCollector.Async = true
	// bugViewCollector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})
	var bugviewItem *model.BugItem
	bugViewCollector.OnResponse(func(resp *colly.Response) {
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body))
		if err != nil {
			zap.L().Error("create goquery document reader error:", zap.Error(err))
			return
		} else {
			resolveTimeStr := doc.Find("div#legendLife>table>tbody>tr:nth-child(3)>td").Text()
			resolveTimeStr = finishedTimeReg.FindString(resolveTimeStr)
			bugviewItem.ResolveTime = resolveTimeStr
			if s.Filter != nil && s.Filter(bugviewItem) {
				s.BugItems = append(s.BugItems, bugviewItem)
			} else {
				zap.L().Info("skip bug", zap.String("ID", bugviewItem.ID))
				// 不能跳过 ，列表没有完成时间排序
				//	next = false
			}
		}
	})
	go func() {
	end:
		for len(s.BugItemWriter) > 0 {
			select {
			case bugviewItem = <-s.BugItemWriter:
			case <-time.After(time.Second * 2):
				break end
			}
			zap.L().Info("visiting bug", zap.String("ID", bugviewItem.ID), zap.String("BugUrl", bugviewItem.BugViewUrl))
			bugViewCollector.Visit(bugviewItem.BugViewUrl)
		}
		waited <- struct{}{}
	}()
	<-waited
	zap.L().Info("bug url crawl finished")
}
func (s *BugSpider) Close() {
	close(s.BugItemWriter)
}
func (s *BugSpider) Crawl() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	sessionCollector, err := s.login()
	if err != nil {
		zap.L().Error("bug spider is  err:", zap.Error(err))
		cancel()
		return err
	}
	go func(c *colly.Collector) {
		collector := c.Clone()
		collector.OnRequest(func(r *colly.Request) {
			collector.SetCookies(r.URL.String(), []*http.Cookie{{Name: "pagerMyBug", Value: strconv.Itoa(s.PagerMyBug)}})
		})
		collector.OnResponse(func(r *colly.Response) {
			doc, err := goquery.NewDocumentFromReader(bytes.NewReader(r.Body))
			if err != nil {
				zap.L().Error("resolve to goquery error", zap.Error(err))
				return
			}
			doc.Find("table#bugList>tbody>tr").EachWithBreak(func(_ int, selector *goquery.Selection) bool {
				i := 0
				for _, n := range selector.Nodes {
					h := colly.NewHTMLElementFromSelectionNode(r, selector, n, i)
					id := h.ChildText("td.c-id")
					level := h.ChildAttr("td>span.label-severity", "title")
					bugType := h.ChildText("td:nth-child(4)")
					bugTitle := h.ChildText("td:nth-child(5)>a")
					createdBy := h.ChildText("td:nth-child(6)")
					resolvedBy := h.ChildText("td:nth-child(8)")
					resolution := h.ChildText("td:nth-child(9)")
					bugviewpage := fmt.Sprintf(bugViewPageTemplate, id)
				end:
					for {
						select {
						case s.BugItemWriter <- &model.BugItem{
							ID:         id,
							Level:      level,
							BugType:    bugType,
							CreatedBy:  createdBy,
							BugTitle:   bugTitle,
							BugViewUrl: bugviewpage,
							ResolveBy:  resolvedBy,
							Resolution: resolution,
						}:
							zap.L().Info("join bug", zap.String("ID", id), zap.String("BugUrl", bugviewpage))
							break end
						case <-ctx.Done():
							zap.L().Info("bugs crawl is done.")
							return false
						}
					}
				}
				return true
			})

		})
		collector.Visit(mybugResolveByUrl)
	}(sessionCollector)
	zap.L().Info("start bug URL collect...")
	s.UrlCollector(cancel)
	return nil
}
func NewBugSpider(sessionSpider *SessionSpider, f func(item *model.BugItem) bool) *BugSpider {
	return &BugSpider{
		SessionSpider: sessionSpider,
		BugItemWriter: make(chan *model.BugItem),
		BugItems:      []*model.BugItem{},
		Filter:        f,
		PagerMyBug:    50,
	}
}
