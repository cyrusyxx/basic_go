# basic_go

## 项目简介

basic_go 是一个基于Go语言的Web应用程序示例，实现了用户管理、文章管理等功能。该项目采用了现代Go Web开发的最佳实践，包括依赖注入、领域驱动设计、分层架构等。

## 项目结构

```
.
├── webook/            # 主应用程序目录
│   ├── main.go        # 应用入口
│   ├── app.go         # 应用结构定义
│   ├── wire.go        # 依赖注入配置
│   ├── config/        # 配置管理
│   ├── internal/      # 内部模块
│   │   ├── domain/    # 领域模型
│   │   ├── repository/# 数据访问层
│   │   ├── service/   # 业务逻辑层
│   │   ├── web/       # Web控制器
│   │   ├── events/    # 事件管理
│   │   └── job/       # 计划任务
│   └── pkg/           # 公共包
├── homework/          # 作业提交目录
└── script/            # 脚本文件
```

## 技术栈

- Go 1.21
- Web框架: Gin
- 数据库: MySQL (GORM)
- 缓存: Redis
- 消息队列: Kafka
- 依赖注入: Wire
- 日志: Zap
- 配置管理: Viper
- 计划任务: Cron
- 监控: Prometheus

## 核心功能

1. **用户管理**：注册、登录、个人信息更新
2. **文章管理**：创建、编辑、发布、查询文章
3. **事件处理**：使用Kafka处理异步事件
4. **计划任务**：定时执行特定任务

## 运行项目

### 前置条件

确保已安装以下服务：
- MySQL
- Redis
- Kafka

### 配置

修改 `webook/config/config.yaml` 配置文件，设置正确的数据库、Redis和Kafka连接信息。

### 启动应用

```bash
cd webook
go run .
```

应用将在 `:8080` 端口启动，监控指标在 `:8081/metrics` 提供。

## 测试

```bash
go test ./...
```

## 作业提交

所有作业应提交到 `homework/` 目录。