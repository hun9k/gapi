package conf

import (
	"log/slog"
	"time"

	"github.com/spf13/viper"
)

func Get[T any](path string, ns ...string) T {
	var result T
	switch any(result).(type) {
	case string:
		return any(Inst(ns...).GetString(path)).(T)
	case bool:
		return any(Inst(ns...).GetBool(path)).(T)
	case int:
		return any(Inst(ns...).GetInt(path)).(T)
	case int64:
		return any(Inst(ns...).GetInt64(path)).(T)
	case float64:
		return any(Inst(ns...).GetFloat64(path)).(T)
	case time.Duration:
		return any(Inst(ns...).GetDuration(path)).(T)
	case []string:
		return any(Inst(ns...).GetStringSlice(path)).(T)
	case map[string]interface{}:
		return any(Inst(ns...).GetStringMap(path)).(T)
	case map[string]string:
		return any(Inst(ns...).GetStringMapString(path)).(T)
	default:
		// 对于自定义类型，尝试使用 Unmarshal
		var target T
		if err := Inst(ns...).UnmarshalKey(path, &target); err != nil {
			slog.Error("failed to unmarshal key", "key", path, "error", err)
		}
		return target
	}
}

var confs = map[string]*viper.Viper{}

const (
	EMPTY_NAME   = "configs"
	DEFAULT_NAME = ""
)

// Inst 是一个获取配置实例的函数
// 它使用单例模式确保全局只有一个配置实例
// 返回值: *conf - 返回配置实例的指针
func Inst(ns ...string) *viper.Viper {
	name := EMPTY_NAME
	if len(ns) > 0 {
		name = ns[0]
	}
	// 检查配置实例是否已经被初始化
	// 如果未初始化(_c == nil)，则创建一个新的配置实例
	if confs[name] == nil {
		confs[name] = newConf(name)
	}

	// 返回配置实例
	return confs[name]
}

func Default() *viper.Viper {
	if confs[DEFAULT_NAME] == nil {
		confs[DEFAULT_NAME] = viper.New()
		// 默认配置
		setDefault(confs[DEFAULT_NAME])
	}

	return confs[DEFAULT_NAME]
}

/**
 * 初始化配置文件
 * @return error 初始化过程中遇到的错误，若无错误则返回nil
 */
func newConf(key string) *viper.Viper {
	c := viper.New()
	// 默认配置
	setDefault(c)

	// 初始化配置选项
	if err := read(c, key); err != nil {
		slog.Warn("failed to read config file", "error", err)
		return c
	}

	// 解析为特定类型
	if err := viper.Unmarshal(&c); err != nil {
		slog.Warn("failed to unmarshal config file", "error", err)
		return c
	}

	// 整理配置
	clean(c)

	return c
}

// config constants
const (
	// app
	APP_NAME_DFT  = "GAPI APP"
	APP_MODE_TEST = "test"
	APP_MODE_PROD = "prod"
	APP_MODE_DEV  = "dev"
	APP_MODE_DFT

	// apiService
	API_ENABLE_DFT       = true
	API_ADDR_DFT         = ":8080"
	API_HTTP3_ENABLE_DFT = true
	API_TLS_ENABLE_DFT   = false
	API_TLS_ADDR_DFT     = ":8443"
	API_TLS_CERTFILE_DFT = "localhost.pem"
	API_TLS_KEYFILE_DFT  = "localhost-key.pem"

	// log
	LOG_FORMAT_JSON   = "json"
	LOG_FORMAT_TEXT   = "text"
	LOG_LEVEL_DEBUG   = slog.LevelDebug
	LOG_LEVEL_WARN    = slog.LevelWarn
	LOG_LEVEL_ERROR   = slog.LevelError
	LOG_LEVEL_INFO    = slog.LevelInfo
	LOG_WRITER_STDOUT = "stdout"
	// default log
	LOG_DEFAULT_FORMAT = LOG_FORMAT_TEXT
	LOG_DEFAULT_LEVEL  = LOG_LEVEL_INFO
	LOG_DEFAULT_WRITER = LOG_WRITER_STDOUT

	// DB default
	DB_DEFAULT_NAME   = ""
	DB_DEFAULT_DRIVER = "mysql"
	DB_DEFAULT_DSN    = "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
)

func setDefault(c *viper.Viper) *viper.Viper {
	c.SetDefault("app.name", APP_NAME_DFT)
	c.SetDefault("app.mode", APP_MODE_DFT)

	c.SetDefault("api.enable", API_ENABLE_DFT)
	c.SetDefault("api.addr", API_ADDR_DFT)
	c.SetDefault("api.http3.enable", API_HTTP3_ENABLE_DFT)
	c.SetDefault("api.tls.enable", API_TLS_ENABLE_DFT)
	c.SetDefault("api.tls.addr", API_TLS_ADDR_DFT)
	c.SetDefault("api.tls.certfile", API_TLS_CERTFILE_DFT)
	c.SetDefault("api.tls.keyfile", API_TLS_KEYFILE_DFT)

	c.SetDefault("log.default.format", LOG_DEFAULT_FORMAT)
	c.SetDefault("log.default.level", LOG_DEFAULT_LEVEL)
	c.SetDefault("log.default.writer", LOG_DEFAULT_WRITER)

	c.SetDefault("db.default.driver", DB_DEFAULT_DRIVER)
	c.SetDefault("db.default.dsn", DB_DEFAULT_DSN)

	return c
}

// 整理Config, 以保证正确的配置
func clean(c *viper.Viper) *viper.Viper {
	// app.mode
	modes := map[string]struct{}{
		APP_MODE_DEV: {}, APP_MODE_PROD: {}, APP_MODE_TEST: {},
	}
	if _, exists := modes[c.GetString("app.mode")]; !exists {
		c.Set("app.mode", APP_MODE_DFT)
	}

	// log.format
	formats := map[string]struct{}{
		LOG_FORMAT_JSON: {}, LOG_FORMAT_TEXT: {},
	}
	if _, exists := formats[c.GetString("log.default.format")]; !exists {
		c.Set("log.default.format", LOG_DEFAULT_FORMAT)
	}

	return c
}

func read(c *viper.Viper, key string) error {
	// 解析配置文件
	// 配置文件名称(无扩展名)
	c.SetConfigName(key)
	// 查找配置文件所在的路径
	c.AddConfigPath(".")

	// 查找并读取配置文件
	if err := c.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
