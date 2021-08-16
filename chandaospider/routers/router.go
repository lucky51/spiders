package routers

import (
	"chandaospider/model"
	"chandaospider/spider"
	"chandaospider/util"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const timeLayout = "2006-01-02 15:04:05"

type ServeWeeklyModel struct {
	Account   string    `form:"account" query:"account" binding:"required"`
	Password  string    `form:"password" query:"password" binding:"required"`
	BeginTime time.Time `form:"beginTime" query:"beginTime" binding:"required" time_format:"2006-01-02 15:04"`
	EndTime   time.Time `form:"endTime" query:"endTime" binding:"required"  time_format:"2006-01-02 15:04"`
}

func NewRouter() *gin.Engine {
	r := gin.Default()
	zap.L().Info("开启http服务...")
	r.Static("/www", "./static")
	r.LoadHTMLGlob("templates/*")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.GET("/weekly", func(c *gin.Context) {
		input := &ServeWeeklyModel{}
		err := c.MustBindWith(input, binding.Query)

		zap.L().Info("request weekly query string", zap.String("password", input.Password), zap.String("account", input.Account), zap.Time("beginTime", input.BeginTime), zap.Time("endTime", input.EndTime))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if len(input.Password) != 32 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "密码MD5长度等于32位",
			})
			return
		}
		err = AttarchSpider(input, c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	})
	return r
}
func AttarchSpider(m *ServeWeeklyModel, c *gin.Context) error {

	//beginTime, endTime := util.GetWeekDay()

	sessionSpider, err := spider.NewSessionSpider(m.Account, m.Password)
	if err != nil {
		return err
	}

	bs := spider.NewBugSpider(sessionSpider, func(item *model.BugItem) bool {
		time, err := time.ParseInLocation(timeLayout, item.ResolveTime, time.Local)
		if err != nil {
			return false
		}
		if time.After(m.BeginTime) && time.Before(m.EndTime) {
			return true
		} else {
			return false
		}
	})
	bugpage := viper.GetInt("app.pagermybug")
	if bugpage > 0 {
		bs.PagerMyBug = bugpage
	}
	defer bs.Close()
	err = bs.Crawl()
	if err != nil {
		return err
	}

	ts := spider.NewTaskSpider(sessionSpider, func(item *model.TaskItem) bool {
		time, err := time.ParseInLocation(timeLayout, item.FinishedTime, time.Local)
		if err != nil {
			return false
		}
		if time.After(m.BeginTime) && time.Before(m.EndTime) {
			return true
		} else {
			return false
		}
	})
	defer ts.Close()
	err = ts.Crawl()
	if err != nil {
		return err
	}
	exportModel := util.NewDocxExporterModel(sessionSpider.GetUserName(), m.BeginTime.Local().Year(), m.EndTime.Local().Year(), m.BeginTime.Local().Day(), m.EndTime.Local().Day(), m.BeginTime.Local().Month(), m.EndTime.Local().Month(), bs.BugItems, ts.TaskItems)
	zap.L().Info("docx exporter", zap.String("input", exportModel.ToJson()))
	fileName := fmt.Sprintf("%s-周报-%s.docx", m.EndTime.Local().Format("20060102"), sessionSpider.GetUserName())
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.Writer.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	return util.InvokeDocxWrite(exportModel, c.Writer)
}
