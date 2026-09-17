# API 接口文档

基础路径：`/api/v1`（后端仅提供接口，不托管前端页面）。

统一响应格式：

```json
{ "code": 0, "msg": "success", "data": {} }
```

- `code = 0`：成功。
- `code = 1`：业务失败（如 `401 未登录` / `403 无权限` / 参数校验错误）。

---

## 一、鉴权

系统支持两种身份来源，**都经过 Casbin 接口级权限校验**：

1. **JWT（浏览器登录）**：`Authorization: Bearer <token>`
2. **API Key（机器 / 脚本 / Skill）**：`X-API-Key: gva_xxxx`

中间件顺序：`APIKeyAuth → CasbinAuth → OperationLog`。若请求携带 `X-API-Key` 则按 Key 查所属用户及其角色；否则回退到 JWT。两者身份都会交给 Casbin 做 `角色, 路径, 方法` 的 `Enforce` 校验。

白名单（跳过 Casbin）：`/auth/me`、`/menus/tree`、`/auth/logout`。

### 1.1 登录获取 Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
# 返回 token、菜单树(menus)、权限标识(permissions)
```

### 1.2 携带 Token 调用

```bash
curl http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer <token>"
```

### 1.3 API Key 调用

```bash
# 创建 Key（原始密钥仅返回一次）
curl -X POST http://localhost:8080/api/v1/apikeys \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"my-skill","scope":"all"}'

# 用 Key 调用，无需登录
curl http://localhost:8080/api/v1/users -H "X-API-Key: gva_xxxx"
```

**API Key 属性**：

- `scope`：`all`（默认，继承所属用户角色的完整增删改查）或 `readonly`（仅允许 GET，写操作返回 403）。
- `expiresAt`：创建时可传 `YYYY-MM-DD`，到期后调用返回 403；留空永不过期。
- 吊销：`DELETE /api/v1/apikeys/:id`（需登录态，只能吊销本人 Key）。

> 首次启动后端会在日志中打印一个演示 Key（`demo api key (admin): gva_...`），拥有管理员全部权限，便于联调。

---

## 二、接口一览

> 鉴权列：`否`=公开，`是`=需 JWT 或 API Key。

### 健康检查

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| GET | `/health` | 健康检查 | 否 |

### 认证与个人

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| POST | `/api/v1/auth/login` | 登录，返回 token / 菜单树 / 权限 | 否 |
| GET | `/api/v1/auth/me` | 当前用户信息（白名单） | 是 |
| POST | `/api/v1/auth/logout` | 登出（白名单） | 是 |
| GET | `/api/v1/menus/tree` | 当前角色菜单树（白名单） | 是 |
| GET | `/api/v1/dashboard` | 仪表盘统计 | 是 |

### 用户 / 角色 / 菜单

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| GET/POST/PUT/DELETE | `/api/v1/users` | 用户管理 | 是 |
| GET/POST/PUT/DELETE | `/api/v1/users/:id` | 单个用户 | 是 |
| GET/POST/PUT/DELETE | `/api/v1/roles` | 角色管理 | 是 |
| GET/POST/PUT/DELETE | `/api/v1/roles/:id` | 单个角色 | 是 |
| POST | `/api/v1/roles/menus` | 给角色分配菜单（自动同步 Casbin 策略） | 是 |
| GET/POST/PUT/DELETE | `/api/v1/menus` | 菜单管理（目录/菜单/按钮） | 是 |
| GET/POST/PUT/DELETE | `/api/v1/menus/:id` | 单个菜单 | 是 |

菜单三种类型：

- `type=1` 目录（Catalog）：仅侧边栏分组，无页面。
- `type=2` 菜单（Menu）：对应前端页面，`component` 指向 `web/src/views/{{component}}.vue`。
- `type=3` 按钮（Button）：对应权限标识（如 `user:add`），前端用于按钮显隐。

### API Key

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| GET | `/api/v1/apikeys` | 当前用户的 API Key 列表 | 是 |
| POST | `/api/v1/apikeys` | 创建 Key（原始密钥仅返回一次） | 是 |
| DELETE | `/api/v1/apikeys/:id` | 吊销 Key | 是 |

### 操作日志

系统自动记录每次已鉴权调用（`操作人 / 方法 / 接口 / IP / 状态码 / 耗时`），并跳过高频接口（`/menus/tree`、`/auth/me`、`/dashboard`、日志自身列表）。

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| GET | `/api/v1/operation-logs` | 日志列表（支持 `page`、`pageSize`、`method`、`username`、`path`） | 是 |
| DELETE | `/api/v1/operation-logs/:id` | 删除单条 | 是 |
| DELETE | `/api/v1/operation-logs` | 清空全部 | 是 |

### 字典

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| GET | `/api/v1/dict-types` | 字典类型列表（筛选/分页） | 是 |
| POST | `/api/v1/dict-types` | 新增字典类型 | 是 |
| PUT | `/api/v1/dict-types/:id` | 更新字典类型 | 是 |
| DELETE | `/api/v1/dict-types/:id` | 删除类型（级联删除其数据） | 是 |
| GET | `/api/v1/dict-data` | 字典数据列表（按 `typeId` / `typeCode` 过滤） | 是 |
| POST | `/api/v1/dict-data` | 新增字典数据 | 是 |
| PUT | `/api/v1/dict-data/:id` | 更新字典数据 | 是 |
| DELETE | `/api/v1/dict-data/:id` | 删除字典数据 | 是 |

### 定时任务

内置无外部依赖调度器，支持 HTTP / 函数两类任务。

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| GET | `/api/v1/tasks` | 任务列表（筛选/分页） | 是 |
| POST | `/api/v1/tasks` | 新增任务 | 是 |
| PUT | `/api/v1/tasks/:id` | 更新任务 | 是 |
| DELETE | `/api/v1/tasks/:id` | 删除任务（含其日志） | 是 |
| POST | `/api/v1/tasks/:id/toggle` | 启用/停用 | 是 |
| POST | `/api/v1/tasks/:id/run` | 手动执行一次 | 是 |
| GET | `/api/v1/task-logs` | 执行日志列表 | 是 |

### 系统参数

键值型参数（站点名称、登录标题、版权、备案号等），支持前端实时生效。

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| GET | `/api/v1/configs/map` | 公开：已启用参数的 `{key: value}` 映射 | 否 |
| GET | `/api/v1/configs` | 参数列表（支持 `key` / `name` / `group` 筛选与分页） | 是 |
| POST | `/api/v1/configs` | 新增参数（参数键唯一） | 是 |
| PUT | `/api/v1/configs/:id` | 更新参数（参数键不可改） | 是 |
| DELETE | `/api/v1/configs/:id` | 删除参数 | 是 |

---

## 三、通用约定

- **分页参数**：列表类接口支持 `page`（默认 1）、`pageSize`（默认 10），响应 `data` 含 `list` 与 `total`。
- **筛选参数**：以查询字符串传入，如 `?method=GET&username=admin&path=/api/v1/users`。
- **内容类型**：`POST/PUT` 请求需带 `Content-Type: application/json`。
- **时间格式**：`expiresAt` 等日期字段使用 `YYYY-MM-DD`。
