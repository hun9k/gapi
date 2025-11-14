package gapi

import (
	"log/slog"

	"github.com/spf13/viper"
)

var conf *Config

func Conf() *Config {
	if conf != nil {
		return conf
	}

	return newConf()
}

func newConf() *Config {
	configInit()

	c := Config{}
	if err := viper.Unmarshal(&c); err != nil {
		slog.Warn(err.Error())
	}

	return &c
}

type Config struct {
	App struct {
		Name string
		Mode string
	}
	HttpService struct {
		Enabled bool
		Addr    string
		Tls     bool
	}
	Log struct {
		Format string
		Output string
		File   struct {
			Filename   string
			MaxSize    int
			MaxBackups int
			MaxAge     int
			Compress   bool
		}
	}
}

// 整理Config, 以保证正确的配置
func (c *Config) Clean() *Config {
	// app.mode
	modes := map[string]struct{}{
		CONF_APP_MODE_DEV: {}, CONF_APP_MODE_PROD: {}, CONF_APP_MODE_TEST: {},
	}
	if _, exists := modes[c.App.Mode]; !exists {
		c.App.Mode = CONF_APP_MODE_DEFAULT
		viper.Set("app.mode", CONF_APP_MODE_DEFAULT)
	}

	return c
}

func configInit() {
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

const (
	CONF_APP_MODE_TEST = "test"
	CONF_APP_MODE_PROD = "prod"
	CONF_APP_MODE_DEV  = "dev"
	CONF_APP_MODE_DEFAULT
)

// 设置默认值
func configSetDefault() {
	viper.SetDefault("app.name", "GAPI APP")
	viper.SetDefault("app.mode", CONF_APP_MODE_DEFAULT)
	viper.SetDefault("httpService.enabled", true)
	viper.SetDefault("httpService.addr", ":8080")
	viper.SetDefault("httpService.tls", false)
	viper.SetDefault("log.format", "text")
	viper.SetDefault("log.file.filename", "./app.log")
	viper.SetDefault("log.file.maxSize", 100)
	viper.SetDefault("log.file.maxBackups", 7)
	viper.SetDefault("log.file.maxAge", 30)
	viper.SetDefault("log.file.compress", true)
	// viper.SetDefault("mysql.enabled", true)
	// viper.SetDefault("mysql.dsn", "root:@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local")
}
