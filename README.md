# go-vue-admin

基于 **Gin + GORM + Vue3 (Element Plus) + Casbin RBAC + JWT** 的后台权限管理脚手架，支持「动态菜单路由 + 前后端双重权限校验」。

## 技术栈

| 层 | 技术 |
| --- | --- |
| 后端 | Gin、GORM、glebarez/sqlite（或 MySQL）、casbin/gorm-adapter、golang-jwt、bcrypt |
| 前端 | Vue3、Vue Router4、Pinia、Element Plus、Axios、Vite |
| 鉴权 | JWT（登录态）+ Casbin（RBAC 接口级权限） |

## 功能特性

- ✅ 账号密码登录，JWT 签发与校验
- ✅ 动态菜单路由：登录后端返回当前角色的**菜单树**，前端据此**动态注册路由并渲染侧边栏**
- ✅ 按钮级权限：`v-permission` 指令由后端返回的**权限标识列表**驱动
- ✅ 接口鉴权：JWT 中间件 + Casbin 中间件**双重校验**每一个受保护接口
- ✅ 用户 / 角色 / 菜单（目录、菜单、按钮）管理
- ✅ 角色分配菜单后**自动同步 Casbin 策略**（由菜单的 `api` + `method` 生成）
- ✅ **API Key 机制**：登录用户可自助生成 API Key，供 skill / 脚本 / 机器通过 `X-API-Key` 头调用接口，复用所属用户 Casbin 权限
- ✅ **操作日志**：中间件自动记录所有已鉴权请求，支持筛选、分页、清空
- ✅ **个人设置**：主题切换（浅色 / 深色 / 跟随系统）
- ✅ **字典管理**：字典类型与字典数据维护
- ✅ **定时任务**：内置无外部依赖调度器，支持 HTTP / 函数两类任务
- ✅ **标签页（Tag Views）**：内容区标签栏，支持关闭/右键批量操作
- ✅ **参数设置（系统参数）**：键值型系统参数，前端实时生效
- ⬜ 代码生成器（预留扩展位）

## 目录结构

```
go-vue-admin/
├── config/                 # 运行配置 + Casbin 模型
│   ├── config.yaml
│   └── rbac_model.conf
├── internal/
│   ├── config/            # 配置加载
│   ├── db/                # GORM 初始化 + 自动迁移
│   ├── casbin/            # Casbin Enforcer 初始化
│   ├── global/            # 全局单例（DB / Enforcer / Config）
│   ├── models/            # 实体模型
│   ├── utils/             # JWT、密码、统一响应
│   ├── middleware/        # JWT 鉴权、Casbin 鉴权、操作日志
│   ├── controller/        # 各业务控制器
│   ├── task/              # 定时任务调度器
│   ├── router/            # 路由注册
│   └── initialize/        # 种子数据
├── web/                   # 前端（Vue3）
├── docs/                   # 文档（部署、API）
├── main.go
└── go.mod
```

## 快速开始

### 后端

```bash
# 需要 Go 1.22+
go mod tidy
go run main.go
# 或编译后运行
go build -o server . && ./server
```

首次启动自动创建 `./data/app.db`（SQLite）并执行种子数据，监听 `:8080`。

### 前端

```bash
cd web
npm install
npm run dev      # 开发服务器 :3000，已配置 /api 代理到 :8080
```

默认管理员账号：**admin / admin123**

## 配置 `config/config.yaml`

```yaml
server:
  port: 8080
  mode: debug            # debug / release（生产务必 release）
jwt:
  secret: "go-vue-admin-secret-change-me"   # 生产请修改
  expire: 72h
database:
  driver: sqlite         # sqlite | mysql
  dsn: "./data/app.db"
  # 使用 MySQL：driver: mysql，dsn: "root:123456@tcp(127.0.0.1:3306)/go_vue_admin?charset=utf8mb4&parseTime=True&loc=Local"
casbin:
  model: "./config/rbac_model.conf"
```

## 权限模型

### 菜单（Menu）三种类型

- `type=1` 目录（Catalog）：仅侧边栏分组，无页面
- `type=2` 菜单（Menu）：对应一个前端页面，`component` 指向 `web/src/views/{{component}}.vue`
- `type=3` 按钮（Button）：对应权限标识（如 `user:add`），用于按钮显隐

### 双重校验流程

1. 登录后后端按角色返回 **菜单树** 与 **权限标识列表**（`permissions`）。
2. 前端用菜单树 **动态注册路由** 并渲染侧边栏；权限列表驱动 `v-permission` 指令。
3. 每个受保护接口经过：
   - **JWT 中间件**：校验 `Authorization: Bearer <token>`，取出 `userId / username / roles`；
   - **Casbin 中间件**：以 `角色, 路径, 方法` 做 `Enforce` 校验。
4. 角色分配菜单时，自动把菜单上的 `api + method` 写入 Casbin 策略，**无需手动维护**。

> 白名单（`/auth/me`、`/menus/tree`、`/auth/logout`）跳过 Casbin，保证登录后必定能拿到菜单与用户信息。

## API Key 与机器调用

除浏览器登录（JWT）外，支持 **API Key 机器鉴权**：Key 归属于某用户，并以该用户角色的 Casbin 策略参与鉴权，skill / 自动化脚本可直接对资源做查询、新增、修改、删除。

- 生成：个人中心 → API 密钥 → 创建；或调用 `POST /api/v1/apikeys`。
- 调用：携带 `X-API-Key: gva_xxxx` 即可，无需登录拿 token。
- 属性：`scope`（`all` / `readonly`）、`expiresAt`（到期返回 403）；吊销 `DELETE /api/v1/apikeys/:id`。
- 首次启动日志会打印一个 **demo admin Key**，便于联调。

完整接口列表、鉴权示例请参阅 [docs/api.md](docs/api.md)。

## 文档

- **部署指南**（构建、Nginx + 二进制、Docker、生产配置、验证）：[docs/deployment.md](docs/deployment.md)
- **API 接口文档**（鉴权、全量接口、通用约定）：[docs/api.md](docs/api.md)

## TODO（预留扩展）

代码生成器等可在 `internal/controller` 与 `web/src/views` 中按现有 CRUD 模式扩展；字典、定时任务、参数设置等内置模块已完成，可作为新功能范本。
