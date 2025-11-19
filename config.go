package gapi

import (
	"log/slog"

	"github.com/spf13/viper"
)

var _config *Config

func Conf() *Config {
	if _config == nil {
		c, err := newConf()
		if err != nil {
			slog.Warn("config failed", "error", err.Error())
		}
		_config = c
	}

	return _config
}

func newConf() (*Config, error) {
	// 默认配置
	c := NewDefaultConf()

	// 初始化配置选项
	if err := confRead(); err != nil {
		return c, err
	}

	// 解析为特定类型
	if err := viper.Unmarshal(&c); err != nil {
		return c, err
	}

	return c, nil
}

type Config struct {
	App struct {
		Name string `yaml:"name,omitempty" json:"name,omitempty"`
		Mode string `yaml:"mode,omitempty" json:"mode,omitempty"`
	} `yaml:"app,omitempty" json:"app,omitempty"`
	HttpService struct {
		Enable bool   `yaml:"enable,omitempty" json:"enable,omitempty"`
		Addr   string `yaml:"addr,omitempty" json:"addr,omitempty"`
		Tls    bool   `yaml:"tls,omitempty" json:"tls,omitempty"`
	} `yaml:"http_service,omitempty" json:"http_service,omitempty"`
	Log struct {
		Format string `yaml:"format,omitempty" json:"format,omitempty"`
		Output string `yaml:"output,omitempty" json:"output,omitempty"`
		File   struct {
			Filename   string `yaml:"filename,omitempty" json:"filename,omitempty"`
			MaxSize    int    `yaml:"max_size,omitempty" json:"max_size,omitempty"`
			MaxBackups int    `yaml:"max_backups,omitempty" json:"max_backups,omitempty"`
			MaxAge     int    `yaml:"max_age,omitempty" json:"max_age,omitempty"`
			Compress   bool   `yaml:"compress,omitempty" json:"compress,omitempty"`
		} `yaml:"file,omitempty" json:"file,omitempty"`
	} `yaml:"log,omitempty" json:"log,omitempty"`
	MySQL struct {
		DSN string `yaml:"dsn,omitempty" json:"dsn,omitempty"`
	} `yaml:"my_sql,omitempty" json:"my_sql,omitempty"`
}

// 整理Config, 以保证正确的配置
func (c *Config) Clean() *Config {
	// app.mode
	modes := map[string]struct{}{
		APP_MODE_DEV: {}, APP_MODE_PROD: {}, APP_MODE_TEST: {},
	}
	if _, exists := modes[c.App.Mode]; !exists {
		c.App.Mode = APP_MODE_DFT
		viper.Set("app.mode", APP_MODE_DFT)
	}

	// log.format
	formats := map[string]struct{}{
		LOG_FORMAT_JSON: {}, LOG_FORMAT_TEXT: {},
	}
	if _, exists := formats[c.Log.Format]; !exists {
		c.Log.Format = LOG_FORMAT_DFT
		viper.Set("log.format", LOG_FORMAT_DFT)
	}

	return c
}

func confRead() error {
	// 解析配置文件
	// 配置文件名称(无扩展名)
	viper.SetConfigName("configs")
	// 查找配置文件所在的路径
	viper.AddConfigPath(".")

	// 查找并读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

// config constants
const (
	// app
	APP_NAME_DFT  = "GAPI APP"
	APP_MODE_TEST = "test"
	APP_MODE_PROD = "prod"
	APP_MODE_DEV  = "dev"
	APP_MODE_DFT

	// httpService
	HTTPSERVICE_ADDR_DFT    = ":8080"
	HTTPSERVICE_ENABLED_DFT = true
	HTTPSERVICE_TLS_DFT     = false

	// log
	LOG_FORMAT_JSON = "json"
	LOG_FORMAT_TEXT = "text"
	LOG_FORMAT_DFT
	LOG_OUTPUT_FILE = "file"
	LOG_OUTPUT_STD  = "std"
	LOG_OUTPUT_DFT
	LOG_FILE_FILENAME_DFT   = "logs/app.log"
	LOG_FILE_MAXSIZE_DFT    = 512  // MB
	LOG_FILE_MAXBACKUPS_DFT = 10   //
	LOG_FILE_MAXAGE_DFT     = 30   //
	LOG_FILE_COMPRESS_DFT   = true //

	// mysql
	MYSQL_DSN_DFT = "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

	// redis
	REDIS_DSN_DFT = ""
)

func NewDefaultConf() *Config {
	c := &Config{}
	c.App.Name = APP_NAME_DFT
	c.App.Mode = APP_MODE_DFT

	c.HttpService.Enable = HTTPSERVICE_ENABLED_DFT
	c.HttpService.Addr = HTTPSERVICE_ADDR_DFT
	c.HttpService.Tls = HTTPSERVICE_TLS_DFT

	c.Log.Format = LOG_FORMAT_DFT
	c.Log.Output = LOG_OUTPUT_DFT
	c.Log.File.Filename = LOG_FILE_FILENAME_DFT
	c.Log.File.MaxSize = LOG_FILE_MAXSIZE_DFT
	c.Log.File.MaxBackups = LOG_FILE_MAXBACKUPS_DFT
	c.Log.File.MaxAge = LOG_FILE_MAXAGE_DFT
	c.Log.File.Compress = LOG_FILE_COMPRESS_DFT

	c.MySQL.DSN = MYSQL_DSN_DFT

	return c
}
