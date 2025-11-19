package tmpls

const Configs = `app:
  name: "{{.App.Name}}" # 自定义应用名称
  mode: {{.App.Mode}} # dev prod test
httpService: # HTTP服务配置
  enable: {{.HttpService.Enable}} # 是否启用
  addr: {{.HttpService.Addr}} # 监听地址
  tls: {{.HttpService.Tls}} # 是否启用https
log: # 日志
  format: {{.Log.Format}} # text json
  output: {{.Log.Output}} # std file
  file: # 如果"output==file"则使用
    filename: {{.Log.File.Filename}} # 日志文件路径
    maxSize: {{.Log.File.MaxSize}} # 单个文件最大(MB)
    maxBackups: {{.Log.File.MaxBackups}} # 最多保留几个备份文件
    maxAge: {{.Log.File.MaxAge}} # 备份文件最大保存多少天
    compress: {{.Log.File.Compress}} # 压缩旧日志（gzip）
mysql:
  dsn: {{.MySQL.DSN}}
redis:
  dsn: localhost:6379
`
