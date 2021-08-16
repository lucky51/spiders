package spider

import (
	"bytes"
	"chandaospider/config"
	"chandaospider/util"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/redisstorage"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const redirectRegStr = "self.location='/user-login"
const (
	userNameSelector string = "ul#userNav span.user-name"
)
const (
	rootUrl      string = "http://192.168.0.178:8011/"
	loginEntry   string = "http://192.168.0.178:8011/user-login-Lw==.html"
	loginPostUrl string = "http://192.168.0.178:8011/user-login.html "
)

type SessionSpider struct {
	account   string
	password  string
	username  string
	Collector *colly.Collector
}

var redirectReg = regexp.MustCompile(redirectRegStr)

// SetSessionStorage set session storage
func (s *SessionSpider) SetSessionStorage() error {
	address := viper.GetString(config.REDIS_SERVER_AND_PORT)
	if address == "" {
		zap.L().Warn(fmt.Sprintf("RedisStorage is disabled ,not find %s in setting.yaml", config.REDIS_SERVER_AND_PORT))
		return nil
	}
	var storage = &redisstorage.Storage{
		Address:  address,
		Password: "",
		DB:       0,
		Prefix:   fmt.Sprintf("chandao-%s", s.account),
	}
	return s.Collector.SetStorage(storage)
}
func (s *SessionSpider) GetUserName() string {
	return s.username
}

// login 登录禅道
func (s *SessionSpider) login() (coll *colly.Collector, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("login error:%v", err)
		}
	}()
	s.Collector.OnResponse(func(resp *colly.Response) {
		if resp.StatusCode != http.StatusOK {
			panic("request error")
		}
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body))
		if err != nil {
			panic(err)
		}
		rawHtm, err := doc.Html()
		if err != nil {
			panic(err)
		}

		if redirectReg.MatchString(rawHtm) {
			loginColl := s.Collector.Clone()
			loginColl.OnHTML("input[name='verifyRand']", func(h *colly.HTMLElement) {

				rvalue := h.Attr("value")
				pwdmd5 := util.GenPasswodStringByMd5Pwd(s.password, rvalue)
				zap.L().Info("account:",
					zap.String("random value", rvalue),
					zap.String("account", s.account),
					zap.String("password", s.password),
					zap.String("final password", pwdmd5))
				frm := map[string]string{
					"account":     s.account,
					"password":    pwdmd5,
					"keepLogin[]": "on",
					"referer":     "/",
					"verifyRand":  rvalue,
				}
				session := loginColl.Clone()
				session.OnRequest(func(r *colly.Request) {
					r.Headers.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.104 Safari/537.36")
					r.Headers.Add("Referer", "/")
					r.Headers.Add("Host", "192.168.0.178:8011")
					r.Headers.Add("Content-Type", "application/x-www-form-urlencoded")
				})
				session.OnResponse(func(r *colly.Response) {
					if r.StatusCode != http.StatusOK {
						panic(fmt.Errorf("登录提交表单返回错误状态码:%d", r.StatusCode))
					}
					if strings.Contains(string(r.Body), "登录失败") {
						panic(errors.New("登录提交表单失败"))
					}
				})
				err := session.Post(loginPostUrl, frm)
				if err != nil {
					panic(err)
				} else {
					mainpagecollector := session.Clone()
					mainpagecollector.OnResponse(func(mainresp *colly.Response) {
						if mainresp.StatusCode != http.StatusOK {
							panic(fmt.Errorf("禅道跟页面返回错误状态码:%d", mainresp.StatusCode))
						}
						maindoc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body))
						if err != nil {
							panic(err)
						}
						s.username = strings.TrimSpace(maindoc.Find(userNameSelector).Text())
						fmt.Println("user name:", s.username)
						s.Collector = mainpagecollector
					})
					mainpagecollector.Visit(rootUrl)
				}
			})
			loginColl.Visit(loginEntry)
		} else {
			s.username = strings.TrimSpace(doc.Find(userNameSelector).Text())
			fmt.Println("user name:", s.username)
		}
	})
	err = s.Collector.Visit(rootUrl)
	if err != nil {
		panic(err)
	}
	s.Collector.Wait()
	coll = s.Collector
	return
}

func NewSessionSpider(account string, password string) (*SessionSpider, error) {
	if account == "" || password == "" {
		return nil, errors.New("account and password is required")
	}
	s := &SessionSpider{
		account:   account,
		password:  password,
		Collector: colly.NewCollector(
		//colly.Async(true),
		),
	}
	//s.Collector.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 4})
	s.Collector.AllowURLRevisit = true
	err := s.SetSessionStorage()
	if err != nil {
		return nil, err
	}
	return s, nil
}
