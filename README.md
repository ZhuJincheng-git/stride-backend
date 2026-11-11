# Stride: Machine Learning-Driven Planning Application

[English](README.md) | [简体中文](README_zh-CN.md)

_Let machine learning optimize our planning!_

## Table of Contents

- [Project Overview](#project-overview)
- [Quick Start](#quick-start)
  - [Requirements](#requirements)
  - [Run Locally](#run-locally)
- [Configuration](#configuration)
- [API Reference](#api-reference)

## Project Overview

A Stride backend implemented in Go. Its purpose is to make it convenient for users to decompose plans and quantify task duration. In the future, it will integrate machine learning to predict task duration and assist with planning.

- Stride allows goals and tasks to have parent-child relationships, enabling natural decomposition of goals or tasks.
- Stride gives every goal and task tags, making classification and statistics easier.
- Stride allows each task to point to a goal, making the T (goal relevance) in the SMART principle more obvious.
- Stride records the work time of each task, facilitating statistics and machine learning.

_Stride currently implements plan decomposition and tags. The machine learning part is coming soon._

## Quick Start

Follow these steps to run the project locally.

### Requirements

```bash
go >= 1.22
mysql
```

### Run Locally

```bash
# Clone the repository
git clone https://github.com/ZhuJincheng-git/stride-backend.git

# Enter the project directory
cd stride-backend

# Install dependencies
go tidy

# Test (optional)
go test ./...

# Generate .env from the template
cp .env.example .env

# Edit .env with an editor and change the configuration to your actual configuration (mainly the database configuration)
vim .env

# Run the development server
go run cmd/server/main.go

# Check whether it is running successfully
curl curl 127.0.0.1:8080/healthz # expected return: {"status":"ok"}
```

## Configuration

| Variable                      | Description                                                              | Default                                    |
| ----------------------------- | ------------------------------------------------------------------------ | ------------------------------------------ |
| `APP_PORT`                    | Application service port                                                 | `8080`                                     |
| `APP_ENV`                     | Application runtime environment; options: `develop` / `test` / `release` | `develop`                                  |
| `DB_HOST`                     | MySQL host address                                                       | `localhost`                                |
| `DB_PORT`                     | MySQL port                                                               | `3306`                                     |
| `DB_USER`                     | MySQL username                                                           | `root`                                     |
| `DB_PASSWORD`                 | MySQL password                                                           | `your-password`                            |
| `DB_NAME`                     | MySQL database name                                                      | `stride_db`                                |
| `DB_PARAMS`                   | MySQL connection parameters                                              | `charset=utf8mb4&parseTime=True&loc=Local` |
| `DB_MAX_IDLE`                 | Maximum number of idle database connections                              | `10`                                       |
| `DB_MAX_OPEN`                 | Maximum number of open database connections                              | `100`                                      |
| `DB_CONN_MAX_LIFETIME_SECOND` | Maximum lifetime of a database connection (seconds)                      | `3600`                                     |
| `JWT_SECRET`                  | JWT signing secret                                                       | `your-jwt-secret`                          |
| `JWT_EXPIRES_HOURS`           | JWT expiration time (hours)                                              | `1`                                        |
| `LOG_LEVEL`                   | Log level; options: `debug` / `info` / `warn` / `error`                  | `info`                                     |

See [`.env.example`](.env.example) for an example `.env`.

## API

See [`docs/api.md`](docs/api.md).
