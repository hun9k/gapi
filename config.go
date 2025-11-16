package gapi

import (
	"log/slog"

	"github.com/spf13/viper"
)

var _config *Config

func Conf() *Config {
	if _config != nil {
		return _config
	}

	_config, err := newConf()
	if err != nil {
		slog.Warn("config failed", "error", err.Error())
	}

	return _config
}

func newConf() (*Config, error) {
	// 设置默认值
	configSetDefault()

	// 解析配置
	c := Config{}
	// 初始化配置选项
	if err := configInit(); err != nil {
		return &c, err
	}

	// 解析为特定类型
	if err := viper.Unmarshal(&c); err != nil {
		return &c, err
	}

	return &c, nil
}

type Config struct {
	App struct {
		Name string `yaml:"name,omitempty" json:"name,omitempty"`
		Mode string `yaml:"mode,omitempty" json:"mode,omitempty"`
	} `yaml:"app,omitempty" json:"app,omitempty"`
	HttpService struct {
		Enabled bool   `yaml:"enabled,omitempty" json:"enabled,omitempty"`
		Addr    string `yaml:"addr,omitempty" json:"addr,omitempty"`
		Tls     bool   `yaml:"tls,omitempty" json:"tls,omitempty"`
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
		CONF_APP_MODE_DEV: {}, CONF_APP_MODE_PROD: {}, CONF_APP_MODE_TEST: {},
	}
	if _, exists := modes[c.App.Mode]; !exists {
		c.App.Mode = CONF_APP_MODE_DEFAULT
		viper.Set("app.mode", CONF_APP_MODE_DEFAULT)
	}

	// log.format
	formats := map[string]struct{}{
		CONF_LOG_FORMAT_JSON: {}, CONF_LOG_FORMAT_TEXT: {},
	}
	if _, exists := formats[c.Log.Format]; !exists {
		c.Log.Format = CONF_LOG_FORMAT_DEFAULT
		viper.Set("log.format", CONF_LOG_FORMAT_DEFAULT)
	}

	return c
}

func configInit() error {
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

const (
	CONF_APP_MODE_TEST = "test"
	CONF_APP_MODE_PROD = "prod"
	CONF_APP_MODE_DEV  = "dev"
	CONF_APP_MODE_DEFAULT

	CONF_LOG_FORMAT_JSON = "json"
	CONF_LOG_FORMAT_TEXT = "text"
	CONF_LOG_FORMAT_DEFAULT
)

// 设置默认值
func configSetDefault() {
	viper.SetDefault("app.name", "GAPI APP")
	viper.SetDefault("app.mode", CONF_APP_MODE_DEFAULT)
	viper.SetDefault("httpService.enabled", true)
	viper.SetDefault("httpService.addr", ":8080")
	viper.SetDefault("httpService.tls", false)
	viper.SetDefault("log.format", CONF_LOG_FORMAT_DEFAULT)
	viper.SetDefault("log.file.filename", "./app.log")
	viper.SetDefault("log.file.maxSize", 100)
	viper.SetDefault("log.file.maxBackups", 7)
	viper.SetDefault("log.file.maxAge", 30)
	viper.SetDefault("log.file.compress", true)
	viper.SetDefault("mysql.dsn", "root:@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local")
}
