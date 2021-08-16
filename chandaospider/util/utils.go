package util

import (
	"chandaospider/config"
	"chandaospider/model"
	"encoding/json"
	"os"

	"crypto/md5"
	"fmt"
	"io"
	"time"

	"github.com/levigross/grequests"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// type DocxExportModel map[string]interface{}

// Md4String 生成md5加密串
func Md5String(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}
func GenPasswordString(pwd string, salt string) string {
	fmt.Println(pwd, salt)
	return Md5String(Md5String(pwd) + salt)
}

// GenPasswodStringByMd5Pwd generate password by md5 password string
func GenPasswodStringByMd5Pwd(md5pwd string, salt string) string {
	return Md5String(md5pwd + salt)
}
func InvokeDocxWrite(data interface{}, writer io.Writer) error {
	docxexportUrl := viper.GetString(config.DOCX_EXPORT_URL)
	postOptions := &grequests.RequestOptions{
		JSON: data,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	resp, err := grequests.Post(docxexportUrl, postOptions)
	if err != nil {
		zap.L().Error("InvokeDocxRender POST Error", zap.Error(err))
		return err
	}
	_, err = writer.Write(resp.Bytes())
	if err != nil {
		return err
	}
	return nil
}
func FixDateNumber(d int) string {
	return fmt.Sprintf("%02d", d)
}

// StrToDateTime 字符串转时间
func StrToDateTime(str string, layOut string) (time.Time, error) {
	return time.ParseInLocation(layOut, str, time.Local)
}

type DocxExportModel struct {
	UserName, Day1, Day2, Month1, Month2 string
	Year1, Year2                         int
	MyBugs                               []*model.BugItem
	MyTasks                              []*model.TaskItem
	HasBug, HasTask                      bool
}

func NewDocxExporterModel(
	userName string,
	year1, year2, day1, day2 int,
	month1, month2 time.Month,
	mybugs []*model.BugItem,
	mytasks []*model.TaskItem) *DocxExportModel {
	return &DocxExportModel{
		UserName: userName,
		MyBugs:   mybugs,
		MyTasks:  mytasks,
		Year1:    year1,
		Year2:    year2,
		Month1:   FixDateNumber(int(month1)),
		Month2:   FixDateNumber(int(month2)),
		Day1:     FixDateNumber(day1),
		Day2:     FixDateNumber(day2),
		HasBug:   len(mybugs) > 0,
		HasTask:  len(mytasks) > 0,
	}
}
func (m *DocxExportModel) SetMyBugs(mybugs []*model.BugItem) {
	m.MyBugs = mybugs
	m.HasBug = len(m.MyBugs) > 0
}

func (m *DocxExportModel) SetTasks(mytasks []*model.TaskItem) {
	m.MyTasks = mytasks
	m.HasTask = len(mytasks) > 0
}
func (m *DocxExportModel) ToJson() string {
	jbytes, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(jbytes)
}

// GetWeekDay 获得当前周的初始和结束日期
func GetWeekDay() (time.Time, time.Time) {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	//周日做特殊判断 因为time.Monday = 0
	if offset > 0 {
		offset = -6
	}

	lastoffset := int(time.Saturday - now.Weekday())
	//周日做特殊判断 因为time.Monday = 0
	if lastoffset == 6 {
		lastoffset = -1
	}

	firstOfWeek := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	lastOfWeeK := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local).AddDate(0, 0, lastoffset+1)
	// f := firstOfWeek.Unix()
	// l := lastOfWeeK.Unix()
	return firstOfWeek, lastOfWeeK //time.Unix(f, 0).Format("2006-01-02") + " 00:00:00", time.Unix(l, 0).Format("2006-01-02") + " 23:59:59"
}

// FileExists 判断路径是否为文件
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
