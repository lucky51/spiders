package test

import (
	"bytes"
	"chandaospider/util"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const timeLayout = "2006-01-02 15:04:05"

func TestParse(t *testing.T) {
	time, err := time.Parse("2006-01-02 15:04:05", "2021-05-08 08:50:56")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(time.Format("2006-01-02 15:04:05"))
	//t.Log(time)
}

func TestGoQuery(t *testing.T) {
	//regStr1 := "(\\d{4}-\\d{2}-\\d{2}\\s+\\d{2}:\\d{2}:\\d{2}),.*?完成"
	regStr := "\\d{4}-\\d{2}-\\d{2}\\s+\\d{2}:\\d{2}:\\d{2}"
	reg := regexp.MustCompile(regStr)
	//	reg1 := regexp.MustCompile(regStr1)
	htmlRaw := `
	<ul id="userNav" class="nav nav-default">
            <li><a class="dropdown-toggle" data-toggle="dropdown"><span class="user-name">王某人</span><span class="caret"></span></a><ul class="dropdown-menu pull-right"><li class="user-profile-item"><a href="/my-profile.html?onlybody=yes" class="iframe" data-width="600"><div class="avatar avatar bg-secondary avatar-circle">W</div>
<div class="user-profile-name">王某人</div><div class="user-profile-role">研发</div></a></li><li class="divider"></li><li><a href="/my-profile.html?onlybody=yes" class="iframe" data-width="600">个人档案</a>
</li><li><a href="/my-changepassword.html?onlybody=yes" class="iframe" data-width="500">更改密码</a>
</li><li class="divider"></li><li class="dropdown-submenu"><a href="javascript:;">主题</a><ul class="dropdown-menu pull-left"><li class="selected"><a href="javascript:selectTheme(&quot;default&quot;);" data-value="default">禅道蓝（默认）</a></li><li><a href="javascript:selectTheme(&quot;green&quot;);" data-value="green">叶兰绿</a></li><li><a href="javascript:selectTheme(&quot;red&quot;);" data-value="red">赤诚红</a></li><li><a href="javascript:selectTheme(&quot;purple&quot;);" data-value="purple">玉烟紫</a></li><li><a href="javascript:selectTheme(&quot;pink&quot;);" data-value="pink">芙蕖粉</a></li><li><a href="javascript:selectTheme(&quot;blackberry&quot;);" data-value="blackberry">露莓黑</a></li><li><a href="javascript:selectTheme(&quot;classic&quot;);" data-value="classic">经典蓝</a></li></ul></li><li class="dropdown-submenu"><a href="javascript:;">Language</a><ul class="dropdown-menu pull-left"><li class="selected"><a href="javascript:selectLang(&quot;zh-cn&quot;);">简体</a></li><li><a href="javascript:selectLang(&quot;zh-tw&quot;);">繁體</a></li><li><a href="javascript:selectLang(&quot;en&quot;);">English</a></li></ul></li><li class="custom-item"><a href="/custom-ajaxMenu-task-view.html?onlybody=yes" data-toggle="modal" data-type="iframe" data-icon="cog" data-width="80%">自定义导航</a></li><li class="divider"></li><li class="dropdown-submenu"><a data-toggle="dropdown">帮助</a><ul class="dropdown-menu pull-left"><li><a href="/tutorial-start.html" class="iframe" data-class-name="modal-inverse" data-width="800" data-headerless="true" data-backdrop="true" data-keyboard="true">新手教程</a>
</li><li><a href="https://www.zentao.net/book/zentaopmshelp.html?fullScreen=zentao" class="open-help-tab" target="_blank">手册</a>
</li><li><a href="/misc-changeLog.html" class="iframe" data-width="800" data-headerless="true" data-backdrop="true" data-keyboard="true">修改日志</a>
</li></ul></li>
<li><a href="/misc-about.html" class="about iframe" data-width="900" data-headerless="true" data-backdrop="true" data-keyboard="true" data-class="modal-about">关于禅道</a>
</li><li class="divider"></li><li><a href="/user-logout.html">退出</a>
</li></ul></li>
          </ul>
	<div class="detail-content">
	<ol class="histories-list">
					  <li value="1">
			  2021-04-21 09:27:35, 由 <strong>李某人</strong> 创建。
					</li>
				<li value="2">
			  2021-04-21 09:36:30, 由 <strong>王某人</strong> 完成。
			  <button type="button" class="btn btn-mini switch-btn btn-icon btn-expand" title="切换显示"><i class="change-show icon icon-plus icon-sm"></i></button>
	  <div class="history-changes" id="changeBox3">
		修改了 <strong><i>总消耗   </i></strong>，旧值为 "0"，新值为 "8"。<br>
修改了 <strong><i>指派给   </i></strong>，旧值为 "aaa"，新值为 "xxx"。<br>
修改了 <strong><i>完成时间</i></strong>，旧值为 ""，新值为 "2021-04-21 09:36:30"。<br>
修改了 <strong><i>预计剩余</i></strong>，旧值为 "3"，新值为 "0"。<br>
修改了 <strong><i>任务状态</i></strong>，旧值为 "wait"，新值为 "done"。<br>
修改了 <strong><i>由谁完成</i></strong>，旧值为 ""，新值为 "aaa"。<br>
	  </div>
					</li>
				<li value="3">
			  2021-04-21 09:37:07, 由 <strong>王某人</strong> 编辑。
			  <button type="button" class="btn btn-mini switch-btn btn-icon btn-expand" title="切换显示"><i class="change-show icon icon-plus icon-sm"></i></button>
	  <div class="history-changes" id="changeBox4">
		修改了 <strong><i>最初预计</i></strong>，旧值为 "3"，新值为 "8"。<br>
	  </div>
					</li>
		</ol>
</div>`

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(htmlRaw)))
	if err != nil {
		t.Log(err)
	}
	userName := doc.Find("ul#userNav span.user-name").Text()
	userName = strings.TrimSpace(userName)
	findStr := "<strong>" + userName + "</strong> 完成"
	t.Log("username", userName)
	doc.Find("ol.histories-list>li").Each(func(i int, sec *goquery.Selection) {
		raw, err := sec.Html()
		t.Log("raw html", raw)
		if err != nil {
			t.Log(err)
		}
		t.Log("find str", findStr)
		if strings.Contains(raw, findStr) {
			t.Log("find string ****:", reg.FindString(raw))
		}
	})
}

func TestDefaultTime(t *testing.T) {
	var t1 time.Time
	var t2 time.Time
	t.Log("time.Time default value : ", t1, t1 == t2, t1.After(t2), t1.Before(t2))
}

func TestBefore(t *testing.T) {
	t1 := time.Now()
	t2 := t1.Add(time.Hour * 24)
	t.Log("t1 < t2", t1.After(t2))
}
func TestParseBefore(t *testing.T) {
	t1, _ := time.ParseInLocation(timeLayout, "2021-04-04 00:00:01", time.Local)
	t2, _ := time.ParseInLocation(timeLayout, "2021-06-04 00:00:01", time.Local)
	n := time.Now()
	t.Log(n.After(t1) && n.Before(t2))
}

func TestNowWeekRange(t *testing.T) {
	t1, t2 := util.GetWeekDay()
	t.Log(t1.Format(timeLayout))
	t.Log(t2.Format(timeLayout))
}
