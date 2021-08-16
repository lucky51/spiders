package config

import (
	"errors"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	APP                   string = "app"
	PORT                  string = "app.port"
	ENV                   string = "app.env"
	DOCX_EXPORT_URL       string = "app.docxExporterUrl"
	ACCOUNT               string = "app.account"
	PASSWORD              string = "app.password"
	REDIS_SERVER_AND_PORT string = "storage.redisServerAndPort"
)

func InitSetting() {
	viper.SetConfigName("setting")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		zap.L().Info("config settings change", zap.String("Event", in.Name))
	})

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s \n", err))
	}
}

func MustGetAccountFromViper() (string, string, error) {
	account, password := viper.GetString(ACCOUNT), viper.GetString(PASSWORD)
	if account == "" {
		return "", "", errors.New("account is required")
	}
	if password == "" {
		return "", "", errors.New("password is required")
	}
	return account, password, nil
}
