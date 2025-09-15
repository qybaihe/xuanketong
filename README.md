# 选课通系统 - 项目说明

本项目是一个前后端分离的在线课程选择与管理平台，旨在为学生提供一个方便的课程信息浏览、评分和评论渠道，并为管理员提供一套高效的管理工具。

## 系统架构

系统整体采用前后端分离的架构模式，由三个独立的部分组成：

1.  **后端服务 (Backend)**: 使用 **Go** 语言开发，提供 RESTful API 接口，负责处理所有业务逻辑、数据存储与用户认证。
2.  **前端应用 (Frontend)**: 使用 **Vue 3** 和 **Vite** 构建的单页面应用 (SPA)，主要面向学生用户，提供课程浏览、搜索、评分、评论等功能。
3.  **管理后台 (Admin Panel)**: 基于 **AdminLTE** (Bootstrap 5) 模板，为管理员提供一个独立的后台管理界面，用于管理用户、课程和评论。

### 架构图

```text
   [学生浏览器]         [管理员浏览器]
        |                    |
        v                    v
[前端 (Vue 3)]    [管理后台 (AdminLTE)]
        |                    |
        +-------+    +-------+
                |    |
                v    v
           [后端 API (Go)]
                |
        +-------+-------+
        |               |
        v               v
    [数据库]         [文件存储]
(SQLite/MySQL)
```

---

## 主要功能

### 学生端
- 用户注册与登录
- 课程列表浏览与关键词搜索
- 查看课程详细信息
- 对课程进行评分
- 发表、查看和管理自己的课程评论

### 管理员端
- **用户管理**: 查看、禁用或删除用户账户。
- **课程管理**: 添加、编辑、删除课程信息。
- **评论管理**: 查看和删除所有用户的评论。

---

## 技术栈

- **后端**: Go
- **前端**: Vue 3, Vite, TypeScript
- **管理后台**: AdminLTE v4, Bootstrap 5, jQuery
- **数据库**: 使用 `test.db` 文件，推测为 SQLite。

---

## 快速开始

在开始之前，请确保您的开发环境中已安装以下软件：
- [Go](https://golang.org/) (建议版本 1.18+)
- [Node.js](https://nodejs.org/) (建议版本 16+)

### 1. 启动后端服务

```bash
# 进入后端项目目录
cd backend

# 安装依赖 (如果需要)
go mod tidy

# 启动 Go 服务 (通常在 http://localhost:8080)
go run main.go
```

### 2. 启动前端应用 (学生端)

```bash
# 进入前端项目目录
cd frontend

# 安装 npm 依赖
npm install

# 启动开发服务器 (通常在 http://localhost:5173)
npm run dev
```

### 3. 启动管理后台 (管理员端)

```bash
# 进入 AdminLTE 目录
cd AdminLTE

# 安装 npm 依赖
npm install

# 启动开发服务器 (通常在 http://localhost:3000)
npm start
```

---

## 目录结构

```
.
├── AdminLTE/      # 管理后台项目 (AdminLTE)
├── backend/       # 后端 Go 项目
├── frontend/      # 前端 Vue 3 项目 (学生端)
├── AGENT.md       # 系统开发文档
└── README.md      # 本项目说明文件
```
