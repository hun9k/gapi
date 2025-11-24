package conf

// public funcs

func App() *app {
	return &confSingle().App
}

func Http() *http {
	return &confSingle().Http
}

func Log() *log {
	return &confSingle().Log
}

func MySQL() *mysql {
	return &confSingle().MySQL
}

// new default conf
func NewDefaultConf() *conf {
	c := &conf{}
	c.App.Name = APP_NAME_DFT
	c.App.Mode = APP_MODE_DFT

	c.Http.Enable = HTTP_ENABLE_DFT
	c.Http.Addr = HTTP_ADDR_DFT
	c.Http.Http3.Enable = HTTP_HTTP3_ENABLE_DFT
	c.Http.Tls.Enable = HTTP_TLS_ENABLE_DFT
	c.Http.Tls.Addr = HTTP_TLS_ADDR_DFT
	c.Http.Tls.CertFile = HTTP_TLS_CERTFILE_DFT
	c.Http.Tls.KeyFile = HTTP_TLS_KEYFILE_DFT

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
