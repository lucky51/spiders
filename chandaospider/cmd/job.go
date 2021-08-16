package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"chandaospider/config"
	"chandaospider/model"
	"chandaospider/spider"
	"chandaospider/util"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var isWithServe bool = false
var exportCmdDesc = strings.Join([]string{}, "\n")
var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "生成周报",
	Long:  exportCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		if !util.FileExists("./setting.yaml") {
			zap.L().Fatal("error: 无法启动任务，请在程序目录创建 setting.yaml 文件")
		}
		config.InitSetting()
		account, password, err := config.MustGetAccountFromViper()
		if err != nil {
			zap.L().Fatal("get setting", zap.Error(err))
		}
		if isWithServe && !isNowWeek {
			port := viper.GetString("app.port")
			if port == "" {
				zap.L().Fatal("请在setting.yaml中配置服务启动的端口号 app.port")
			}
			go StartHttpServe(port)
		}
		job := &JobRunner{
			account:  account,
			password: password,
		}
		if !isNowWeek {
			c := cron.New()
			expression := viper.GetString("job.cron")
			if expression == "" {
				zap.L().Fatal("setting.yaml job.cron is required")
			}
			sche, err := cron.ParseStandard(expression)
			if err != nil {
				zap.L().Fatal("cron expression invalid", zap.Error(err))
			}
			c.Schedule(sche, job)
			zap.L().Info("starting job ....")
			c.Start()
			zap.L().Info("started job")
			defer c.Stop()
			select {}
		} else {
			job.Run()
		}
	},
}

type JobRunner struct {
	password string
	account  string
}

func (j *JobRunner) Run() {
	zap.L().Info("Starting docx render")
	beginTime, endTime := util.GetWeekDay()
	sessionSpider, err := spider.NewSessionSpider(j.account, j.password)
	if err != nil {
		zap.L().Error("initialize session spider error", zap.Error(err))
		return
	}
	bs := spider.NewBugSpider(sessionSpider, func(item *model.BugItem) bool {
		time, err := time.ParseInLocation(timeLayout, item.ResolveTime, time.Local)
		if err != nil {
			zap.L().Error("parse in location error", zap.Error(err))
			return false
		}
		if time.After(beginTime) && time.Before(endTime) {
			return true
		} else {
			return false
		}
	})
	defer bs.Close()
	err = bs.Crawl()
	if err != nil {
		zap.L().Error("crawl bugs error", zap.Error(err))
		return
	}
	ts := spider.NewTaskSpider(sessionSpider, func(item *model.TaskItem) bool {
		time, err := time.ParseInLocation(timeLayout, item.FinishedTime, time.Local)
		if err != nil {
			zap.L().Error("parse in location error", zap.Error(err))
			return false
		}
		if time.After(beginTime) && time.Before(endTime) {
			return true
		} else {
			return false
		}
	})
	defer ts.Close()
	err = ts.Crawl()
	if err != nil {
		zap.L().Error("crawl tasks error", zap.Error(err))
		return
	}
	exportModel := util.NewDocxExporterModel(bs.GetUserName(), beginTime.Local().Year(), endTime.Local().Year(), beginTime.Local().Day(), endTime.Local().Day(), beginTime.Local().Month(), endTime.Local().Month(), bs.BugItems, ts.TaskItems)
	fileName := fmt.Sprintf("./%s-周报-%s.docx", endTime.Local().Format("20060102"), sessionSpider.GetUserName())
	file, err := os.OpenFile(fileName, os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		fmt.Printf("%v \n", err)
		return
	}
	defer file.Close()
	util.InvokeDocxWrite(exportModel, file)
}

func init() {
	jobCmd.Flags().BoolVarP(&isNowWeek, "now", "n", false, "是否立刻输出周报，非后台任务。")
	jobCmd.Flags().BoolVarP(&isWithServe, "serve", "s", false, "是否开启http服务")
}
