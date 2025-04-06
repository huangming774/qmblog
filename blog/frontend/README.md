# 博客系统前端

这是一个基于Vue 3的博客系统前端项目。

## 技术栈

- Vue 3
- Vuex 4
- Vue Router 4
- Element Plus UI
- Axios
- Marked (Markdown渲染)
- Highlight.js (代码高亮)

## 项目结构

```
frontend/
  ├── public/              # 静态资源目录
  ├── src/                 # 源代码
  │   ├── assets/          # 资源文件(图片、样式等)
  │   ├── components/      # 公共组件
  │   ├── router/          # 路由配置
  │   ├── store/           # 状态管理
  │   ├── views/           # 页面视图
  │   ├── App.vue          # 根组件
  │   └── main.js          # 入口文件
  ├── .gitignore           # Git忽略文件
  ├── babel.config.js      # Babel配置
  ├── package.json         # 依赖配置
  └── README.md            # 项目说明
```

## 功能特性

- 文章列表展示
- 文章详情阅读
- 用户登录注册
- 文章评论互动
- 文章管理(创建、编辑、删除)
- 响应式设计，适配移动端和桌面端

## 安装和使用

1. 安装依赖
```bash
npm install
```

2. 开发环境运行
```bash
npm run serve
```

3. 生产环境构建
```bash
npm run build
```

## 后端API接口

本项目需要与后端API配合使用，API接口基础URL为：`http://localhost:8080/api/v1`

## 默认账号

- 管理员账号: admin@example.com / admin123
- 测试用户账号: demo@example.com / demo123 