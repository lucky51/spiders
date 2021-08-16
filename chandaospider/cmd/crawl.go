package cmd

import (
	"chandaospider/model"
	"chandaospider/spider"
	"chandaospider/util"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var crawlway string
var crawlBegin string
var crawlEnd string
var account string
var password string
var outType string
var isNowWeek bool
var crawwayDesc = strings.Join([]string{
	"该子命令支持如下类型数据:",
	"bugs: 爬取我完成的bugs数据",
	"tasks: 爬取我完成的任务数据",
	"all:爬取完成的bugs和tasks数据",
}, "\n")
var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "爬取数据",
	Long:  crawwayDesc,
	Run: func(cmd *cobra.Command, args []string) {
		if len(password) != 32 {
			fmt.Println("密码格式不正确，请输入 32 位 密码的 MD5串")
			return
		}
		if password == "" {
			fmt.Println("请输入用户名")
			return
		}
		if outType != "json" && outType != "txt" && outType != "docx" {
			fmt.Println("请输入正确的 导出文件类型 ，支持  'docx','txt','json'")
			return
		}
		nonelimited := crawlBegin == "" && crawlEnd == "" && !isNowWeek
		var beginTime, endTime time.Time
		var exportModel *util.DocxExportModel
		if !isNowWeek {
			beginTime, err := util.StrToDateTime(timeLayout, crawlBegin)
			if err != nil {
				fmt.Printf("begin:%s 转换失败，请核对时间格式, yyyy-MM-dd HH:mm:ss \n", crawlBegin)
				return
			}
			endTime, err := util.StrToDateTime(timeLayout, crawlEnd)
			if err != nil {
				fmt.Printf("end:%s 转换失败，请核对时间格式, yyyy-MM-dd HH:mm:ss \n", crawlEnd)
				return
			}
			if beginTime.Local().After(endTime) {
				fmt.Println("开始时间 不能 大于结束时间")
				return
			}
		} else {
			beginTime, endTime = util.GetWeekDay()
		}
		sessionSpider, err := spider.NewSessionSpider(account, password)
		if err != nil {
			zap.L().Fatal("initialize session spider error", zap.Error(err))
		}
		switch crawlway {
		case "bugs":
			s := spider.NewBugSpider(sessionSpider, func(item *model.BugItem) bool {
				if nonelimited {
					return true
				}
				time, err := time.ParseInLocation(timeLayout, item.ResolveTime, time.Local)
				if err != nil {
					return false
				}
				if time.After(beginTime) && time.Before(endTime) {
					return true
				} else {
					return false
				}
			})
			s.Crawl()
			exportModel = util.NewDocxExporterModel(s.GetUserName(), beginTime.Local().Year(), endTime.Local().Year(), beginTime.Local().Day(), endTime.Local().Day(), beginTime.Local().Month(), endTime.Local().Month(), s.BugItems, []*model.TaskItem{})

		case "tasks":
			s := spider.NewTaskSpider(sessionSpider, func(item *model.TaskItem) bool {
				if nonelimited {
					return true
				}
				time, err := time.ParseInLocation(timeLayout, item.FinishedTime, time.Local)
				if err != nil {
					return false
				}
				if time.After(beginTime) && time.Before(endTime) {
					return true
				} else {
					return false
				}
			})
			s.Crawl()
			exportModel = util.NewDocxExporterModel(s.GetUserName(), beginTime.Local().Year(), endTime.Local().Year(), beginTime.Local().Day(), endTime.Local().Day(), beginTime.Local().Month(), endTime.Local().Month(), []*model.BugItem{}, s.TaskItems)
		case "all":
			bs := spider.NewBugSpider(sessionSpider, func(item *model.BugItem) bool {
				if nonelimited {
					return true
				}
				time, err := time.ParseInLocation(timeLayout, item.ResolveTime, time.Local)
				if err != nil {
					return false
				}
				if time.After(beginTime) && time.Before(endTime) {
					return true
				} else {
					return false
				}
			})
			bs.Crawl()
			ts := spider.NewTaskSpider(sessionSpider, func(item *model.TaskItem) bool {
				if nonelimited {
					return true
				}
				time, err := time.ParseInLocation(timeLayout, item.FinishedTime, time.Local)
				if err != nil {
					return false
				}
				if time.After(beginTime) && time.Before(endTime) {
					return true
				} else {
					return false
				}
			})
			ts.Crawl()
			exportModel = util.NewDocxExporterModel(bs.GetUserName(), beginTime.Local().Year(), endTime.Local().Year(), beginTime.Local().Day(), endTime.Local().Day(), beginTime.Local().Month(), endTime.Local().Month(), bs.BugItems, ts.TaskItems)
		default:
			fmt.Println("not supported type")
			return
		}
		if exportModel != nil {
			fileName := fmt.Sprintf("./%s-周报-%s.docx", endTime.Local().Format("20060102"), sessionSpider.GetUserName())
			file, err := os.OpenFile(fileName, os.O_TRUNC|os.O_CREATE, 0600)
			if err != nil {
				zap.L().Fatal("open file error:", zap.Error(err))
			}
			defer file.Close()
			switch outType {
			case "docx":
				err = util.InvokeDocxWrite(exportModel, file)
				if err != nil {
					zap.L().Fatal("invoke docx write error", zap.Error(err))
				}
			case "txt", "json":
				jdata, err := json.Marshal(exportModel)
				if err != nil {
					zap.L().Fatal("marshal error", zap.Error(err))
				}
				_, err = file.Write(jdata)
				if err != nil {
					zap.L().Fatal("file write error", zap.Error(err))
				}
			}

		}
	},
}

func init() {
	crawlCmd.Flags().StringVarP(&crawlway, "type", "t", "bugs", "选取所要爬取的数据类型")
	crawlCmd.Flags().StringVarP(&crawlBegin, "begin", "b", "", "输入爬取的开始时间，格式:yyyy-MM-dd HH:mm:ss")
	crawlCmd.Flags().StringVarP(&crawlEnd, "end", "e", "", "输入爬取的结束时间，格式:yyyy-MM-dd HH:mm:ss")
	crawlCmd.Flags().StringVarP(&account, "account", "a", "", "禅道用户名")
	crawlCmd.Flags().StringVarP(&password, "password", "p", "", "禅道密码的MD5密文模式")
	crawlCmd.Flags().StringVarP(&outType, "output", "o", "docx", "导出文件类型  txt ,json , docx")
	crawlCmd.Flags().BoolVarP(&isNowWeek, "now", "n", true, "导出当前周 优先级 大于时间范围 begin ~ end")
}
