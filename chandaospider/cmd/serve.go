package cmd

import (
	"chandaospider/config"
	"chandaospider/routers"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var serverCmd = &cobra.Command{
	Use:   "serve",
	Short: "http服务",
	Long:  "开启周报导出的http服务",
	Run: func(cmd *cobra.Command, args []string) {
		config.InitSetting()
		port := viper.GetString("app.port")
		if port == "" {
			zap.L().Fatal("请在setting.yaml中配置服务启动的端口号 app.port")
		}
		StartHttpServe(port)
	},
}

func StartHttpServe(port string) {
	r := routers.NewRouter()
	r.Run(fmt.Sprintf(":%s", port))
}
