package tmpls

const Configs = `app:
  name: "{{.appName}}" # 自定义应用名称
  mode: dev # dev prod test
api: # api服务配置
  addr: :8080 # 监听地址
  http3:
    enable: false # 是否启用HTTP/3
  tls: 
    enable: false # 是否启用https
    addr: :8443
    certfile: localhost.pem
    keyfile: localhost-key.pem
tast: # 任务服务配置
  enable: false # 是否启用
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

const DockerCompose = `version: "3"
services:
  mysql:  # docker exec -it  {{.dbName}}-mysql-1 mysql -psecret
    image: mysql:8
    command: mysqld --character-set-server=utf8mb4
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: secret
      MYSQL_DATABASE: {{.dbName}}
    volumes:
      - ./volumes/mysql/data:/var/lib/mysql
    ports:
      - '8306:3306'

`

const GitIgnore = `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with go test -c
*.test

# Code coverage profiles and other test artifacts
*.out
coverage.*
*.coverprofile
profile.cov

# Dependency directories (remove the comment below to include it)
# vendor/

# Go workspace file
go.work
go.work.sum

# env file
.env

# Editor/IDE
.idea/
.vscode/
`

const Dockerfile = ``
