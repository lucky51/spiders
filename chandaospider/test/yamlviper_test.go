package test

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
)

func TestViperYaml(t *testing.T) {
	yamlBytes := []byte(`app:
  port: 3002
  env: debug
  docxExporterUrl: http://localhost:3001/file
  account: wangfucheng
  password: 36e48c5d322c489dd3ee4e51ffda1494
storage:
  redisServerAndPort: 127.0.0.1:6379
job:
  cron: 10 17 * * 5
`)
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewReader(yamlBytes)) // Find and read the config file
	if err != nil {                                     // Handle errors reading the config file
		t.Fatalf("read yaml error:%v", err)
	}
	appstr := viper.GetString("app")
	appmap := viper.GetStringMapString("app")
	for key, val := range appmap {
		t.Logf("key =%s , value=%s \r\n", key, val)
	}
	app_port := viper.GetString("app.port")
	app_cron := viper.GetString("job.cron")
	t.Logf("%s ,%s ,%s ,%s", appstr, appmap, app_port, app_cron)
}
