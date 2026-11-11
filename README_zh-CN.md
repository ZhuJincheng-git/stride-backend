# Stride: 机器学习驱动的规划应用

[English](README.md) | [简体中文](README_zh-CN.md)

_让机器学习优化我们的规划！_

## 目录

- [项目概述](#项目概述)
- [快速开始](#快速开始)
  - [环境要求](#环境要求)
  - [本地运行](#本地运行)
- [配置说明](#配置说明)
- [API 参考](#api-参考)

## 项目概述

使用Go实现的Stride后端，目的是让用户便捷地将计划分解，量化任务耗时，后续将集成机器学习预测任务耗时，辅助规划。

- Stride让目标（Goal）和任务（Task）都有父子关系，实现目标或任务的自然分解。
- Stride让每个目标和任务都有标签，便于分类和统计。
- Stride让每个任务可以指向一个目标，让SMART原则的T（目标相关）更加显而易见。
- Stride让每个任务的工作时间得以记录，便于统计与机器学习。

_Stride目前已经实现了计划分解和标签，机器学习部分敬请期待。_

## 快速开始

按照以下步骤在本地运行该项目。

### 环境要求

```bash
go >= 1.22
mysql
```

### 本地运行

```bash
# 克隆仓库
git clone https://github.com/ZhuJincheng-git/stride-backend.git

# 进入项目目录
cd stride-backend

# 安装依赖
go tidy

# 测试（可选）
go test ./...

# 基于模板生成.env
cp .env.example .env

# 用编辑器修改.env内配置为实际配置（主要为数据库配置）
vim .env

# 运行开发服务器
go run cmd/server/main.go

# 检查是否运行成功
curl curl 127.0.0.1:8080/healthz # 预期返回：{"status":"ok"}
```

## 配置说明

| 变量名                        | 说明                                                | 默认值                                     |
| ----------------------------- | --------------------------------------------------- | ------------------------------------------ |
| `APP_PORT`                    | 应用服务端口                                        | `8080`                                     |
| `APP_ENV`                     | 应用运行环境，可选：`develop` / `test` / `release`  | `develop`                                  |
| `DB_HOST`                     | MySQL 主机地址                                      | `localhost`                                |
| `DB_PORT`                     | MySQL 端口                                          | `3306`                                     |
| `DB_USER`                     | MySQL 用户名                                        | `root`                                     |
| `DB_PASSWORD`                 | MySQL 密码                                          | `your-password`                            |
| `DB_NAME`                     | MySQL 数据库名                                      | `stride_db`                                |
| `DB_PARAMS`                   | MySQL 连接参数                                      | `charset=utf8mb4&parseTime=True&loc=Local` |
| `DB_MAX_IDLE`                 | 数据库最大空闲连接数                                | `10`                                       |
| `DB_MAX_OPEN`                 | 数据库最大打开连接数                                | `100`                                      |
| `DB_CONN_MAX_LIFETIME_SECOND` | 数据库连接最大生命周期（秒）                        | `3600`                                     |
| `JWT_SECRET`                  | JWT 签名密钥                                        | `your-jwt-secret`                          |
| `JWT_EXPIRES_HOURS`           | JWT 过期时间（小时）                                | `1`                                        |
| `LOG_LEVEL`                   | 日志级别，可选：`debug` / `info` / `warn` / `error` | `info`                                     |

示例 `.env` 见文件 [`.env.example`](.env.example)。

## API 参考

详见 [`docs/api.md`](docs/api.md)。
