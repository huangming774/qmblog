# 博客系统后端

这是一个基于Go语言Gin框架的博客系统后端。

## 技术栈

- Gin: 轻量级Web框架
- GORM: ORM库
- PostgreSQL: 关系型数据库
- Redis: 缓存数据库
- JWT: 用户认证

## 项目结构

```
backend/
  ├── api/            # API文档
  ├── config/         # 配置文件
  ├── controllers/    # 控制器
  ├── middlewares/    # 中间件
  ├── models/         # 数据模型
  ├── repositories/   # 数据仓库
  ├── services/       # 业务逻辑
  ├── utils/          # 工具函数
  ├── scripts/        # 脚本
  ├── main.go         # 程序入口
  └── go.mod          # Go模块文件
```

## 运行要求

- Go 1.18+
- PostgreSQL 13+
- Redis 6+

## 如何运行

1. 确保安装了PostgreSQL和Redis，并启动这些服务

2. 初始化数据库
```bash
createdb blog
```

3. 下载依赖
```bash
go mod tidy
```

4. 运行服务器
```bash
go run main.go
```

5. 初始化测试数据（可选）
```bash
go run scripts/init.go
```

## API文档

服务运行后，API接口列表：

- `POST /api/v1/auth/register`: 注册用户
- `POST /api/v1/auth/login`: 用户登录
- `GET /api/v1/posts`: 获取文章列表
- `GET /api/v1/posts/:id`: 获取文章详情
- `POST /api/v1/posts`: 创建文章
- `PUT /api/v1/posts/:id`: 更新文章
- `DELETE /api/v1/posts/:id`: 删除文章
- `GET /api/v1/posts/:id/comments`: 获取文章评论
- `POST /api/v1/posts/:id/comments`: 创建评论
- `DELETE /api/v1/comments/:id`: 删除评论

## 缓存策略

系统使用Redis实现多种缓存策略，提升性能：

1. **文章内容缓存**
   - 使用Hash结构存储文章信息
   - 缓存键格式：`post:{id}`
   - 包含字段：文章ID、标题、内容、摘要等
   - 过期时间：24小时
   - 缓存更新策略：文章更新或删除时主动删除缓存

2. **阅读计数缓存**
   - 使用String结构记录文章阅读次数
   - 缓存键格式：`post_view:{id}`
   - 每增加10次阅读，更新到数据库
   - 避免频繁数据库写操作

缓存数据使用延迟双删策略确保一致性，在高并发场景下提高读取性能。

## 并发处理

系统利用Go语言的goroutine特性实现高效并发处理：

1. **异步缓存操作**
   - Redis缓存读写与主业务逻辑并行执行
   - 阅读计数的增加和持久化异步进行
   - 缓存更新和删除操作不阻塞API响应

2. **并行数据库查询**
   - 文章列表查询中并行获取总数和分页数据
   - 使用WaitGroup同步多个并发操作结果

3. **性能优化**
   - 减少API响应时间，提高用户体验
   - 降低高并发场景下的系统负载
   - 非关键操作异步化，优先保障核心响应速度

这种并发设计使系统能够更有效地处理高流量请求，特别适合博客这类读多写少的应用场景。

## 默认用户

初始管理员账户：
- 用户名: admin
- 邮箱: admin@example.com
- 密码: admin123 