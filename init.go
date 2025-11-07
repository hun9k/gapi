package gapi

import (
	"io"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

func init() {
	// 初始化配置
	initConfig()

	// 初始化日志
	initLog()
}

func initLog() {
	// 声明witter, options, logger
	// 基于配置设置Writter
	var writer io.Writer
	switch viper.GetString("log.output") {
	case "file":
	case "std":
		fallthrough
	default:
		writer = os.Stdout
	}

	// 声明选项
	options := &slog.HandlerOptions{}
	// 设置最低级别日志
	switch viper.GetString("mode") {
	case "prod":
		options.Level = slog.LevelInfo
	case "dev":
		fallthrough
	default:
		options.Level = slog.LevelDebug
	}

	logger := &slog.Logger{}
	// 基于类型设置handler
	switch viper.GetString("log.format") {
	case "json":
		logger = slog.New(slog.NewJSONHandler(writer, options))
	case "text":
		fallthrough
	default:
		logger = slog.New(slog.NewTextHandler(writer, options))
	}
	slog.SetDefault(logger)
}

func initConfig() {
	// 设置默认值
	configSetDefault()

	// 解析配置文件
	// 配置文件名称(无扩展名)
	viper.SetConfigName("configs")
	// 查找配置文件所在的路径
	viper.AddConfigPath(".")

	// 查找并读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		slog.Warn(err.Error())
	}
}

// 设置默认值
func configSetDefault() {
	viper.SetDefault("name", "BACKEND API")
	viper.SetDefault("mode", "dev")
	viper.SetDefault("httpService.enabled", true)
	viper.SetDefault("httpService.addr", ":8080")
	viper.SetDefault("log.format", "text")
	// viper.SetDefault("mysql.enabled", true)
	// viper.SetDefault("mysql.dsn", "root:@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local")
}
