# FastSprint - 敏捷开发协作平台

一个多人协作的敏捷开发平台，参考 Jira 和 Confluence，同时具备类似 Obsidian 的知识库功能。

## 技术栈

### 前端
- **框架**: React 18 + TypeScript
- **构建工具**: Vite
- **UI 组件库**: Ant Design 5
- **状态管理**: Redux Toolkit
- **路由**: React Router v6
- **HTTP 客户端**: Axios
- **日期处理**: Day.js

### 后端
- **语言**: Go 1.21
- **Web 框架**: Gin
- **ORM**: GORM
- **微服务架构**: 按业务领域拆分服务

### 数据存储
- **关系型数据库**: PostgreSQL 15
- **缓存**: Redis 7

## 项目结构

```
fastsprint/
├── backend/                    # 后端服务
│   ├── common/                 # 公共模块
│   │   ├── config/            # 配置管理
│   │   ├── database/          # 数据库连接
│   │   ├── response/          # 统一响应格式
│   │   ├── middleware/        # 中间件（鉴权、CORS等）
│   │   ├── utils/             # 工具函数
│   │   └── models/            # 基础模型
│   ├── user-service/          # 用户服务 (端口: 8001)
│   │   ├── cmd/               # 入口
│   │   ├── internal/
│   │   │   ├── handler/       # HTTP 处理器
│   │   │   ├── service/       # 业务逻辑
│   │   │   ├── repository/    # 数据访问
│   │   │   └── model/         # 数据模型
│   │   └── config/            # 配置文件
│   ├── project-service/       # 项目服务 (端口: 8002)
│   └── task-service/          # 任务服务 (端口: 8003)
├── frontend/                  # 前端应用
│   ├── src/
│   │   ├── api/               # API 接口
│   │   ├── components/        # 公共组件
│   │   ├── pages/             # 页面组件
│   │   ├── store/             # Redux 状态管理
│   │   ├── router/            # 路由配置
│   │   ├── layouts/           # 布局组件
│   │   ├── utils/             # 工具函数
│   │   └── types/             # TypeScript 类型定义
│   └── ...
├── deploy/                    # 部署配置
│   ├── docker/                # Dockerfile
│   ├── sql/                   # 数据库初始化脚本
│   ├── docker-compose.yml     # Docker Compose 配置
│   └── nginx.conf             # Nginx 配置
└── docs/                      # 文档
```

## 核心功能

### 第一阶段：项目管理（Jira 部分）

- **用户管理**
  - 用户注册、登录
  - JWT 身份认证
  - 用户信息管理

- **项目管理**
  - 项目创建、编辑、删除
  - 项目成员管理（角色：管理员、成员、查看者）
  - 项目状态配置（待办、进行中、已完成等）
  - 任务类型配置（任务、故事、缺陷、子任务）
  - 优先级配置（最高、高、中、低、最低）

- **任务管理**
  - 任务创建、编辑、删除
  - 任务状态流转
  - 任务分配
  - 任务评论
  - 任务变更记录
  - 任务附件（预留）

- **看板视图**
  - 按状态分列展示任务
  - 快速创建任务
  - 点击查看任务详情

- **冲刺管理**
  - 冲刺创建、编辑、删除
  - 冲刺启动、完成
  - 冲刺任务管理

### 第二阶段：知识库（Confluence + Obsidian 部分）
- （待开发）

## 快速开始

### 前置要求
- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- Redis 7+
- Docker & Docker Compose (推荐)

### 使用 Docker Compose 启动（推荐）

```bash
cd deploy
docker-compose up -d
```

访问地址:
- 前端: http://localhost:3000
- 用户服务: http://localhost:8001
- 项目服务: http://localhost:8002
- 任务服务: http://localhost:8003

### 本地开发

#### 1. 启动数据库

```bash
cd deploy
docker-compose up -d postgres redis
```

#### 2. 初始化数据库

```bash
# 执行 deploy/sql/init.sql
```

#### 3. 启动后端服务

```bash
# 用户服务
cd backend/user-service
go mod tidy
go run cmd/main.go

# 项目服务 (新终端)
cd backend/project-service
go mod tidy
go run cmd/main.go

# 任务服务 (新终端)
cd backend/task-service
go mod tidy
go run cmd/main.go
```

#### 4. 启动前端

```bash
cd frontend
npm install
npm run dev
```

访问 http://localhost:3000

## API 文档

### 用户服务 (port: 8001)

#### 认证接口
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录

#### 用户接口
- `GET /api/v1/users/me` - 获取当前用户信息
- `PUT /api/v1/users/me` - 更新当前用户信息
- `POST /api/v1/users/change-password` - 修改密码
- `GET /api/v1/users/:id` - 获取指定用户信息
- `GET /api/v1/users` - 获取用户列表

### 项目服务 (port: 8002)

#### 项目接口
- `POST /api/v1/projects` - 创建项目
- `GET /api/v1/projects` - 获取项目列表
- `GET /api/v1/projects/:id` - 获取项目详情
- `PUT /api/v1/projects/:id` - 更新项目
- `DELETE /api/v1/projects/:id` - 删除项目
- `GET /api/v1/projects/key/:key` - 按 key 获取项目

#### 成员接口
- `GET /api/v1/projects/:id/members` - 获取成员列表
- `POST /api/v1/projects/:id/members` - 添加成员
- `DELETE /api/v1/projects/:id/members/:userId` - 移除成员
- `PUT /api/v1/projects/:id/members/:userId/role` - 修改成员角色

#### 配置接口
- `GET /api/v1/projects/:id/statuses` - 获取状态列表
- `POST /api/v1/projects/:id/statuses` - 创建状态
- `GET /api/v1/projects/:id/issue-types` - 获取任务类型列表
- `POST /api/v1/projects/:id/issue-types` - 创建任务类型
- `GET /api/v1/projects/:id/priorities` - 获取优先级列表
- `POST /api/v1/projects/:id/priorities` - 创建优先级

### 任务服务 (port: 8003)

#### 任务接口
- `POST /api/v1/issues` - 创建任务
- `GET /api/v1/issues` - 获取任务列表
- `GET /api/v1/issues/:id` - 获取任务详情
- `PUT /api/v1/issues/:id` - 更新任务
- `DELETE /api/v1/issues/:id` - 删除任务
- `GET /api/v1/issues/key/:key` - 按 key 获取任务

#### 评论接口
- `GET /api/v1/issues/:id/comments` - 获取评论列表
- `POST /api/v1/issues/:id/comments` - 添加评论
- `DELETE /api/v1/issues/:id/comments/:commentId` - 删除评论

#### 变更记录接口
- `GET /api/v1/issues/:id/changelogs` - 获取变更记录

#### 冲刺接口
- `POST /api/v1/sprints` - 创建冲刺
- `GET /api/v1/sprints` - 获取冲刺列表
- `GET /api/v1/sprints/:id` - 获取冲刺详情
- `PUT /api/v1/sprints/:id` - 更新冲刺
- `DELETE /api/v1/sprints/:id` - 删除冲刺
- `POST /api/v1/sprints/:id/start` - 启动冲刺
- `POST /api/v1/sprints/:id/complete` - 完成冲刺

#### 看板接口
- `GET /api/v1/board` - 获取看板数据

## 数据库设计

数据库采用 Schema 隔离的方式，每个微服务拥有独立的 Schema：

- `user_service` - 用户服务 Schema
  - `users` - 用户表
  - `roles` - 角色表
  - `user_roles` - 用户角色关联表

- `project_service` - 项目服务 Schema
  - `projects` - 项目表
  - `project_members` - 项目成员表
  - `statuses` - 状态配置表
  - `issue_types` - 任务类型表
  - `priorities` - 优先级表

- `task_service` - 任务服务 Schema
  - `issues` - 任务表
  - `issue_comments` - 任务评论表
  - `issue_attachments` - 任务附件表
  - `issue_changelogs` - 任务变更记录表
  - `issue_watchers` - 任务关注表
  - `sprints` - 冲刺表

## 开发计划

### 已完成
- [x] 项目架构设计
- [x] 数据库设计
- [x] 后端公共模块
- [x] 用户服务（注册、登录、用户管理）
- [x] 项目服务（项目CRUD、成员管理、配置管理）
- [x] 任务服务（任务CRUD、评论、变更记录、看板）
- [x] 冲刺服务
- [x] 前端基础框架
- [x] 登录注册页面
- [x] 项目列表页面
- [x] 看板页面
- [x] 任务详情页面
- [x] Docker Compose 部署配置

### 待开发
- [ ] API 网关服务
- [ ] 服务间通信（gRPC/HTTP）
- [ ] 任务拖拽（看板拖拽排序）
- [ ] 高级搜索和筛选
- [ ] 数据统计和报表
- [ ] 通知服务
- [ ] 文件上传服务
- [ ] 知识库模块（Wiki）
- [ ] 双向链接（Obsidian 风格）
- [ ] 知识图谱
- [ ] 单元测试
- [ ] 集成测试
- [ ] CI/CD 流水线

## 贡献指南

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License
