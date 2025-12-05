# github.com/hun9k/gapi/conf

## 使用方法

## configs.yaml 示例

```yaml
app:
  name: "GAPI APP" # 自定义应用名称
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
  format: json # text json
  level: info # debug info warn error
  output: file://?filename=logs/app.log&maxSize=100&maxBackups=7&maxAge=30&compress=true
db:
  - driver: mysql
    dsn: root:secret@tcp(localhost:8306)/gapidemo?charset=utf8mb4&parseTime=True&loc=Local
```
