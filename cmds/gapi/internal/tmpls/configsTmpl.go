package tmpls

const Configs = `app:
  name: "YOUR APP NAME" # 自定义应用名称
  mode: dev # dev prod test
api: # api服务配置
  enable: true # 是否启用
  addr: :8080 # 监听地址
  http3:
    enable: true # 是否启用HTTP/3
  tls: 
    enable: true # 是否启用https
    addr: :8443
    certfile: localhost.pem
    keyfile: localhost-key.pem
tast: # 任务服务配置
  enable: true # 是否启用
log: # 日志
  default: # 默认日志配置
    format: text # text json
    level: info # debug info warn error
    writer: stdout # stdout
    # writer: file://?filename=logs/app.log&maxSize=100&maxBackups=7&maxAge=30&compress=true # url style
db:
  default: 
    driver: mysql
    dsn: root:secret@tcp(localhost:8306)/gapidemo?charset=utf8mb4&parseTime=True&loc=Local
`
