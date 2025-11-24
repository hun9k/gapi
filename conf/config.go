package conf

import (
	"log/slog"

	"github.com/spf13/viper"
)

// config constants
const (
	// app
	APP_NAME_DFT  = "GAPI APP"
	APP_MODE_TEST = "test"
	APP_MODE_PROD = "prod"
	APP_MODE_DEV  = "dev"
	APP_MODE_DFT

	// httpService
	HTTP_ADDR_DFT         = ":8080"
	HTTP_ENABLE_DFT       = true
	HTTP_HTTP3_ENABLE_DFT = true
	HTTP_TLS_ENABLE_DFT   = false
	HTTP_TLS_ADDR_DFT     = ":8443"
	HTTP_TLS_CERTFILE_DFT = "localhost.pem"
	HTTP_TLS_KEYFILE_DFT  = "localhost-key.pem"

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

// app section
type app struct {
	Name string `yaml:"name,omitempty"`
	Mode string `yaml:"mode,omitempty"`
}

// http section
type http struct {
	Enable bool   `yaml:"enable,omitempty"`
	Addr   string `yaml:"addr,omitempty"`
	Http3  struct {
		Enable bool `yaml:"enable,omitempty"`
	} `yaml:"http3,omitempty"`
	Tls struct {
		Enable   bool   `yaml:"enable,omitempty"`
		Addr     string `yaml:"addr,omitempty"`
		CertFile string `yaml:"certfile,omitempty"`
		KeyFile  string `yaml:"keyfile,omitempty"`
	} `yaml:"tls,omitempty"`
}

// log section
type log struct {
	Format string `yaml:"format,omitempty"`
	Output string `yaml:"output,omitempty"`
	File   struct {
		Filename   string `yaml:"filename,omitempty"`
		MaxSize    int    `yaml:"max_size,omitempty"`
		MaxBackups int    `yaml:"max_backups,omitempty"`
		MaxAge     int    `yaml:"max_age,omitempty"`
		Compress   bool   `yaml:"compress,omitempty"`
	} `yaml:"file,omitempty"`
}

// mysql section
type mysql struct {
	DSN string `yaml:"dsn,omitempty"`
}

// all conf
type conf struct {
	App   app   `yaml:"app"`
	Http  http  `yaml:"http"`
	Log   log   `yaml:"log"`
	MySQL mysql `yaml:"mysql"`
}

var _instance *conf

func confSingle() *conf {
	if _instance == nil {
		c, err := confNew()
		if err != nil {
			slog.Warn("config failed", "error", err.Error())
		}
		_instance = c
	}

	return _instance
}

// 先使用默认值初始化
// 再解析配置文件内容，如果解析失败，返回默认值，否则返回解析后的值
func confNew() (*conf, error) {
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

	// 整理配置
	c.clean()

	return c, nil
}

// 整理Config, 以保证正确的配置
func (c *conf) clean() *conf {
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
