# 博客系统

一个前后端分离的博客系统，前端使用Vue 3，后端使用Go语言的Gin框架。

## 项目概述

这是一个完整的博客系统，具有用户管理、文章管理、评论互动等功能。系统采用前后端分离架构，前端使用Vue 3框架，后端使用Go语言的Gin框架。

### 技术栈

**前端**：
- Vue 3
- Vuex 4
- Vue Router 4
- Element Plus UI
- Axios

**后端**：
- Go1.24
- Gin框架
- GORM
- PostgreSQL
- Redis
- JWT认证

## 系统功能

- 用户注册与登录
- 文章发布与管理
- 标签管理
- 评论与回复
- 文章搜索
- 响应式设计

## 项目结构

```
blog/
  ├── backend/           # Go后端项目
  │   ├── api/           # API文档
  │   ├── config/        # 配置文件
  │   ├── controllers/   # 控制器
  │   ├── middlewares/   # 中间件
  │   ├── models/        # 数据模型
  │   ├── services/      # 业务逻辑
  │   ├── utils/         # 工具函数
  │   └── main.go        # 程序入口
  │
  ├── frontend/          # Vue前端项目
  │   ├── public/        # 静态资源
  │   ├── src/           # 源代码
  │   └── package.json   # 依赖配置
  │
  └── README.md          # 项目说明
```

## 运行项目

### 后端

1. 确保已安装PostgreSQL和Redis，并启动这些服务

2. 创建数据库
```bash
createdb blog
```

3. 进入后端目录
```bash
cd backend
```

4. 下载依赖
```bash
go mod tidy
```

5. 运行服务器
```bash
go run main.go
```

6. 初始化测试数据（可选）
```bash
go run scripts/init.go
```

### 前端

1. 进入前端目录
```bash
cd frontend
```

2. 安装依赖
```bash
npm install
```

3. 开发环境运行
```bash
npm run serve
```

4. 生产环境构建
```bash
npm run build
```

## 默认账号

- 管理员账号: admin@example.com / admin123
- 测试用户账号: demo@example.com / demo123

## 开发注意事项

- 后端API基础URL：`http://localhost:8080/api/v1`
- 前端开发服务器：`http://localhost:8080` 
