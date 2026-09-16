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
- ✅ **API Key 机制**：登录用户可在「个人中心 → API 密钥」中**自助生成** API Key，外部 **skill / 脚本 / 机器** 通过 `X-API-Key` 头调用接口（查询、新增、修改、删除等），复用所属用户的 Casbin 权限做双重校验
- ✅ **操作日志**：中间件自动记录所有已鉴权请求（操作人、方法、接口、IP、状态码、耗时），管理员可在「系统管理 → 操作日志」查看、筛选与清空
- ✅ **个人设置**：右上角头像下拉「设置」可切换主题（浅色 / 深色 / 跟随系统）
- ✅ **字典管理**：维护字典类型与字典数据（枚举值），可供下拉、状态等场景复用
- ⬜ 代码生成器、定时任务（预留扩展位）

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
│   ├── models/            # User / Role / Menu / ApiKey / OperationLog / DictType / DictData 模型
│   ├── utils/             # JWT、密码、统一响应
│   ├── middleware/        # JWT 鉴权、Casbin 鉴权、操作日志
│   ├── controller/        # 登录、用户、角色、菜单、仪表盘、操作日志
│   ├── router/            # 路由注册
│   └── initialize/        # 种子数据（默认角色/管理员/菜单/策略）
├── web/                   # 前端（Vue3）
│   └── src/
│       ├── api/           # 接口封装
│       ├── store/         # Pinia（登录、动态路由生成）
│       ├── router/        # 静态路由 + 路由守卫
│       ├── layout/        # 布局 + 递归侧边栏
│       ├── utils/         # request、permission 指令
│       └── views/         # 登录、仪表盘、系统管理页面
├── main.go
└── go.mod
```

## 快速开始

### 1. 后端

```bash
# 需要 Go 1.22+
go mod tidy
go run main.go
# 或编译后运行
go build -o server . && ./server
```

首次启动会自动创建 `./data/app.db`（SQLite）并执行种子数据。**监听 `:8080`**。

### 2. 前端

```bash
cd web
npm install
npm run dev      # 开发服务器 :3000，已配置 /api 代理到 :8080
# 生产构建
npm run build    # 产物输出到 web/dist
```

默认账号：**admin / admin123**

## 配置文件 `config/config.yaml`

```yaml
server:
  port: 8080
  mode: debug            # debug / release

jwt:
  secret: "go-vue-admin-secret-change-me"
  expire: 72h

database:
  driver: sqlite         # sqlite | mysql
  dsn: "./data/app.db"
  # 使用 MySQL 时改为：
  # driver: mysql
  # dsn: "root:123456@tcp(127.0.0.1:3306)/go_vue_admin?charset=utf8mb4&parseTime=True&loc=Local"

casbin:
  model: "./config/rbac_model.conf"
```

> 切换为 MySQL：将 `driver` 改为 `mysql` 并填写 `dsn`，GORM 会自动建表。

## 权限模型

### 菜单（Menu）三种类型
- `type=1` 目录（Catalog）：仅用于侧边栏分组，无页面
- `type=2` 菜单（Menu）：对应一个前端页面，`component` 指向 `web/src/views/{{component}}.vue`
- `type=3` 按钮（Button）：对应一个权限标识（如 `user:add`），前端用于按钮显隐

### 双重校验流程
1. 登录成功后，后端根据角色查到菜单，返回 **菜单树** 与 **权限标识列表**（`permissions`）。
2. 前端用菜单树 **动态注册路由** 并渲染侧边栏；用权限列表驱动 `v-permission` 指令做按钮级控制。
3. 每个受保护接口都经过：
   - **JWT 中间件**：校验 `Authorization: Bearer <token>`，取出 `userId / username / roles`；
   - **Casbin 中间件**：以 `角色, 请求路径, 请求方法` 在 Casbin 中做 `Enforce` 校验。
4. 角色分配菜单时，后端自动把菜单上的 `api + method` 写入 Casbin 策略，**无需手动维护策略**。

> 白名单（`/auth/me`、`/menus/tree`、`/auth/logout`）跳过 Casbin 校验，保证登录后必定能拿到菜单与用户信息。

## API Key 与机器调用（Skill / 脚本）

除浏览器登录（JWT）外，本系统支持 **API Key 机器鉴权**：一个 API Key 归属于某个用户，并以该用户角色的 Casbin 策略参与鉴权，因此 skill / 自动化脚本可以直接对资源做**查询、新增、修改、删除**等操作，且受同一套权限约束。

### 1. 生成 API Key
登录后进入「个人中心 → API 密钥」点击「创建密钥」，或在任意已登录会话调用：

```bash
curl -X POST http://localhost:8080/api/v1/apikeys \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"my-skill"}'
# 返回：{ "code":0, "data": { "id":1, "name":"my-skill", "key":"gva_xxxx", "secret":"gva_xxxx", "hint":"请妥善保存，密钥仅展示一次" } }
```

`key` 仅在此刻展示一次，请妥善保存。也可在后端日志中看到首次启动时生成的演示 Key（`demo api key (admin): gva_...`），该 Key 拥有管理员全部权限，便于快速联调。

### 2. Skill 调用接口
携带 `X-API-Key` 头即可，无需走登录拿 token：

```bash
# 查询用户列表（GET）
curl http://localhost:8080/api/v1/users -H "X-API-Key: gva_xxxx"

# 新增用户（POST）
curl -X POST http://localhost:8080/api/v1/users \
  -H "X-API-Key: gva_xxxx" -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"alice123","nickname":"Alice"}'

# 修改用户（PUT）
curl -X PUT http://localhost:8080/api/v1/users/2 \
  -H "X-API-Key: gva_xxxx" -H "Content-Type: application/json" \
  -d '{"nickname":"Alice2"}'

# 删除用户（DELETE）
curl -X DELETE http://localhost:8080/api/v1/users/2 -H "X-API-Key: gva_xxxx"
```

### 3. 鉴权与权限
- 中间件 `APIKeyAuth`：若请求携带 `X-API-Key`，则按 Key 查出所属用户及其角色，注入上下文；否则回退到 `JWTAuth`（Bearer Token）。两种方式的身份都会交给 `CasbinAuth` 做接口级校验。
- 受角色菜单中的 `api + method` 约束：skill 只能调用其所属用户角色被授权的接口（管理员 Key 默认拥有 `/* *`）。
- 吊销：`DELETE /api/v1/apikeys/:id`（需登录态，且只能吊销本人 Key）。
- **权限范围 `scope`**：`all`（默认，继承所属用户角色的完整增删改查权限）或 `readonly`（仅允许 GET 查询，写操作返回 403），创建时指定。
- **过期时间 `expiresAt`**：创建时可传入 `YYYY-MM-DD`，到期后调用返回 403；留空则永不过期。
- 范围与过期时间都可在「个人中心 → API 密钥」创建时设置，列表中展示范围、过期时间及是否已过期。

## 操作日志

系统通过 `OperationLog` 中间件**自动记录每一次已通过鉴权的接口调用**，用于审计与排障。日志记录在响应返回后落库，并自动跳过高频/框架类接口（`/menus/tree`、`/auth/me`、`/dashboard`、操作日志自身列表）以避免刷屏。

记录的字段：`操作人`、`请求方法`、`接口路径`、`IP`、`状态码`、`耗时(ms)`、`操作时间`。

### 查看与管理
使用管理员登录后，进入「系统管理 → 操作日志」：
- 按 **方法 / 操作人 / 接口路径** 筛选；
- 支持分页查看；
- 可单条删除或一键清空。

### 接口
- 列表：`GET /api/v1/operation-logs`（支持 `page`、`pageSize`、`method`、`username`、`path` 参数）
- 删除单条：`DELETE /api/v1/operation-logs/:id`
- 清空全部：`DELETE /api/v1/operation-logs`

## 接口一览

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| GET | `/health` | 健康检查 | 否 |
| POST | `/api/v1/auth/login` | 登录，返回 token / 菜单树 / 权限 | 否 |
| GET | `/api/v1/auth/me` | 当前用户 | 是 |
| POST | `/api/v1/auth/logout` | 登出 | 是 |
| GET | `/api/v1/menus/tree` | 当前角色菜单树 | 是 |
| GET | `/api/v1/dashboard` | 仪表盘统计 | 是 |
| GET/POST/PUT/DELETE | `/api/v1/users` | 用户管理 | 是 |
| GET/POST/PUT/DELETE | `/api/v1/roles` | 角色管理 | 是 |
| POST | `/api/v1/roles/menus` | 给角色分配菜单（同步 Casbin） | 是 |
| GET/POST/PUT/DELETE | `/api/v1/menus` | 菜单管理 | 是 |
| GET | `/api/v1/apikeys` | 当前用户的 API Key 列表 | 是 |
| POST | `/api/v1/apikeys` | 创建 API Key（原始密钥仅返回一次） | 是 |
| DELETE | `/api/v1/apikeys/:id` | 吊销 API Key | 是 |
| GET | `/api/v1/operation-logs` | 操作日志列表（支持筛选/分页） | 是 |
| DELETE | `/api/v1/operation-logs/:id` | 删除单条日志 | 是 |
| DELETE | `/api/v1/operation-logs` | 清空全部日志 | 是 |
| GET | `/api/v1/dict-types` | 字典类型列表（支持筛选/分页） | 是 |
| POST | `/api/v1/dict-types` | 新增字典类型 | 是 |
| PUT | `/api/v1/dict-types/:id` | 更新字典类型 | 是 |
| DELETE | `/api/v1/dict-types/:id` | 删除字典类型（级联删除其数据） | 是 |
| GET | `/api/v1/dict-data` | 字典数据列表（按 typeId / typeCode 过滤） | 是 |
| POST | `/api/v1/dict-data` | 新增字典数据 | 是 |
| PUT | `/api/v1/dict-data/:id` | 更新字典数据 | 是 |
| DELETE | `/api/v1/dict-data/:id` | 删除字典数据 | 是 |

统一响应格式：

```json
{ "code": 0, "msg": "success", "data": {} }
```

`code=0` 表示成功；`code=1` 表示业务失败（如 `401 未登录 / 403 无权限`）。

## 部署提示

- 生产请将 `server.mode` 改为 `release`。
- 修改 `jwt.secret` 为强随机值。
- 前端 `npm run build` 后，可将 `web/dist` 交由 Nginx 托管，并反向代理 `/api` 到后端 `:8080`。
- 侧边栏递归组件依赖菜单的 `path`（目录为绝对路径如 `/system`，子菜单为相对路径如 `user`），新增菜单时请保持一致。

## TODO（预留扩展）

代码生成器、字典管理、定时任务等内置模块可在 `internal/controller` 与 `web/src/views` 中按现有 CRUD 模式继续扩展。
