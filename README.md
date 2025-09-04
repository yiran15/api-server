# apiserver

## 简介

GO 编写的 WEB 框架，支持用户管理、角色管理、接口权限管理、OAuth2 登录等功能，此框架基于流行的`Gin`框架构建。

前端地址：[https://github.com/yiran15/ui.git](https://github.com/yiran15/ui.git)

预览地址：[清泉流响](https://qqlx.net/)

账号: `readonly@qqlx.net`
密码: `12345678`

## 技术栈

- gin ---> web 框架
- zap ---> 日志
- gorm ---> mysql 数据持久化
- go-redis ---> 缓存
- wire ---> 依赖注入
- casbin ---> 访问控制
- otel ---> 可观测性
- oauth2 ---> OAuth2 登录

## API 文档

- 支持 [swagger 文档](https://swagger.io/)
- [项目 API 文档](https://qqlx.net/swagger/index.html)

![alt text](docs/img/swagger.png)

## 功能

### 用户管理: 增删改查

![用户管理](docs/img/user.png)

### 角色管理: 增删改查

![角色管理](docs/img/role.png)

### 接口权限管理: 增删改查

![接口权限管理](docs/img/api.png)

### OAuth2 登录

支持 OAuth2 登录，目前支持飞书、keycloak。

![OAuth2 登录](docs/img/oauth2-1.png)
![OAuth2 登录](docs/img/oauth2-feishu.png)

## 可观测性

基于`otel`的可观测性，包括`trace`、`metrics`。

使用 [阿里云 otel](https://github.com/alibaba/loongsuite-go-agent) 构建镜像，自动注入`trace`、`metrics`。

### 配置 trace & metrics 导出

- OTEL_EXPORTER_OTLP_ENDPOINT: otlp 服务地址
- OTEL_EXPORTER_OTLP_PROTOCOL: otlp 服务协议
- OTEL_SERVICE_NAME: 服务名称
- OTEL_METRICS_EXPORTER: metrics 导出方式为`prometheus`格式
- OTEL_EXPORTER_PROMETHEUS_PORT: metrics 导出端口
- OTEL_EXPORTER_PROMETHEUS_HOST: metrics 导出主机

```bash
cat deploy/docker-compose.yaml
services:
  apiserver:
    image: api-server
    restart: always
    container_name: apiserver
    environment:
      - CONFIG_PATH=/app/config.yaml
      - OTEL_EXPORTER_OTLP_ENDPOINT=http://10.10.10.10:30001
      - OTEL_EXPORTER_OTLP_PROTOCOL=grpc
      - OTEL_SERVICE_NAME=api-server
      # - OTEL_METRICS_EXPORTER=prometheus
      # - OTEL_EXPORTER_PROMETHEUS_PORT=9464
      # - OTEL_EXPORTER_PROMETHEUS_HOST=0.0.0.0
```

### Trace

![Trace](docs/img/trace.png)

### Metrics

![Metrics](docs/img/metrics.png)

## 快速开始

### 配置

所有配置项均支持环境变量配置，环境变量前缀从环境变量`SERVICE_NAME`获取。

- 例如`SERVICE_NAME=qqlx`，则环境变量前缀为`QQLX`
- 例如`SERVICE_NAME=api-server`，则环境变量前缀为`API_SERVER`
- 如果未设置`SERVICE_NAME`，则环境变量前缀默认为`API_SERVER`。

配置中时间单位为“ns”、“us”（或“µs”）、“ms”、“s”、“m”、“h”。

指定配置文件路径:

- 命令行参数: `-c`或`--config-path`。
- 环境变量: `${SERVICE_NAME}_CONFIG_PATH`。

```yaml
server:
  bind: 0.0.0.0:8080
  timeZone: "Asia/Shanghai"
log:
  level: debug
mysql:
  # 开启后会打印 sql 语句
  debug: true
  username: xxx
  password: xxx
  host: xxx
  port: 3306
  database: xxx
  maxIdleConns: 10
  maxOpenConns: 20
  maxLifetime: 30m
redis:
  # single sentinel
  mode: single
  host: 127.0.0.1:6379
  password: 123456
  # 过期时间 3s 3m 3h
  expireTime: 300s
  # redis key 前缀
  keyPrefix: tutu
  db: 0
  poolSize: 20
  minIdleConns: 10
  connMaxLifetime: 30m
jwt:
  issuer: tutu
  # jwt secret
  secret: 123456
  # token 过期时间
  expireTime: 9999h
oauth2:
  # 是否启用 oauth2
  enable: true
  providers:
    feishu:
      clientId: xxx
      clientSecret: xxx
      # scopes: []
      authUrl: https://accounts.feishu.cn/open-apis/authen/v1/authorize
      tokenUrl: https://open.feishu.cn/open-apis/authen/v2/oauth/token
      userInfoUrl: https://open.feishu.cn/open-apis/authen/v1/user_info
      # 回调地址, host 为前端地址
      redirectUrl: http://10.0.0.10:5173/oauth/login
    keycloak:
      clientId: xxx
      clientSecret: xxx
      scopes:
        - openid
        - email
        - profile
        - roles
      authUrl: https://keycloak.qqlx.net/realms/qqlx/protocol/openid-connect/auth
      tokenUrl: https://keycloak.qqlx.net/realms/qqlx/protocol/openid-connect/token
      userInfoUrl: https://keycloak.qqlx.net/realms/qqlx/protocol/openid-connect/userinfo
      # 回调地址, host 为前端地址
      redirectUrl: http://10.0.0.10:5173/oauth/login
```

### 部署

```bash
cd deploy

# 初始化配置文件, 需要修改配置文件中的数据库信息
mv config-example.yaml config.yaml

# 初始化数据库
mysql -h 127.0.0.1 -P 3306 -u root -p my_database < schema.sql

# 构建镜像
export DOCKER_USERNAME=xxx
export DOCKER_PASSWORD=xxx

# 构建并启动容器
make all
```
