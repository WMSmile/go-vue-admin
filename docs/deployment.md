# 部署指南

本系统由两部分组成：

- **前端**：Vue3 + Vite 构建出的纯静态 SPA（`web/dist/`）。
- **后端**：Go 编译出的单一二进制，仅提供 `/api/v1` 接口（**不托管前端静态文件**）。

因此生产环境需要一个反向代理（推荐 Nginx）同时完成两件事：**托管前端静态资源** + **将 `/api/` 反代到后端**。

> 默认后端监听 `:8080`，接口前缀 `/api/v1`，配置来自 `config/config.yaml` 与 `config/rbac_model.conf`（相对路径加载）。

---

## 一、构建

### 1.1 前端

```bash
cd web
npm install
npm run build        # 产物输出到 web/dist/
```

- 若部署到子路径（如 `https://example.com/admin/`），需 `vite build --base=/admin/` 并同步调整前端路由 `base`，否则使用根路径部署无需改动。
- 前端 API 基址写死为 `/api/v1`（`web/src/utils/request.js`），与部署域名同源，无需额外环境变量。

### 1.2 后端

```bash
# 仓库根目录
go mod tidy
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/go-vue-admin .
```

- 数据库驱动 `glebarez/sqlite` 为纯 Go 实现，`CGO_ENABLED=0` 可得到静态二进制，免装 glibc 兼容包。
- 其他平台：`GOOS=darwin`（macOS）、`GOARCH=arm64`（ARM 服务器）。

构建产物需与以下文件一起部署：

```
go-vue-admin            # 二进制
config/config.yaml
config/rbac_model.conf
```

运行时会在工作目录下生成 `./data/app.db`（SQLite，默认）。

---

## 二、方案 A：Nginx + Go 二进制（单机 VPS 主流做法）

### 2.1 目录布局

```
/opt/go-vue-admin/
├── go-vue-admin            # 后端二进制
├── config/
│   ├── config.yaml
│   └── rbac_model.conf
├── data/                   # 运行时生成的 app.db（需写权限）
└── web/dist/               # 前端构建产物
```

### 2.2 Nginx 配置

```nginx
server {
  listen 80;
  server_name your.domain.com;

  root /opt/go-vue-admin/web/dist;
  index index.html;

  # SPA 回退：刷新 /system/role 这类前端路由时返回 index.html，交给前端路由处理。
  # 缺少这行会导致直接刷新深链接时出现 404（文件不存在）。
  location / {
    try_files $uri $uri/ /index.html;
  }

  # 反代 API 到 Go 后端
  location /api/ {
    proxy_pass http://127.0.0.1:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
  }
}
```

### 2.3 后端守护进程（systemd）

`/etc/systemd/system/go-vue-admin.service`：

```ini
[Unit]
Description=go-vue-admin
After=network.target

[Service]
WorkingDirectory=/opt/go-vue-admin
ExecStart=/opt/go-vue-admin/go-vue-admin
Restart=always
User=www-data

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now go-vue-admin
```

---

## 三、方案 B：Docker（容器化 / CI-CD）

推荐**前后端两个容器** + `docker-compose`，Nginx 容器负责静态托管与反代。

### 3.1 后端 Dockerfile（`Dockerfile.backend`）

```dockerfile
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/go-vue-admin .

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/go-vue-admin /app/go-vue-admin
COPY config /app/config
RUN apk add --no-cache ca-certificates
EXPOSE 8080
CMD ["/app/go-vue-admin"]
```

### 3.2 前端 + Nginx Dockerfile（`Dockerfile.frontend`）

```dockerfile
FROM node:20-alpine AS build
WORKDIR /web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM nginx:alpine
COPY --from=build /web/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
```

`nginx.conf`（容器内）：

```nginx
server {
  listen 80;
  root /usr/share/nginx/html;
  index index.html;
  location / { try_files $uri $uri/ /index.html; }
  location /api/ {
    proxy_pass http://backend:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
  }
}
```

### 3.3 `docker-compose.yml`

```yaml
services:
  backend:
    build:
      context: .
      dockerfile: Dockerfile.backend
    volumes:
      - ./data:/app/data
    restart: always

  frontend:
    build:
      context: .
      dockerfile: Dockerfile.frontend
    ports:
      - "80:80"
    depends_on:
      - backend
    restart: always
```

```bash
docker compose up -d --build
```

---

## 四、生产配置建议

修改 `config/config.yaml`：

```yaml
server:
  port: 8080
  mode: release            # 务必 release，关闭 Gin 调试模式
jwt:
  secret: "换成随机长字符串"   # 不要用默认值
database:
  driver: sqlite
  dsn: "/opt/go-vue-admin/data/app.db"   # 改为绝对路径，避免 cwd 问题
casbin:
  model: "/opt/go-vue-admin/config/rbac_model.conf"  # 生产建议绝对路径
```

要点：

1. **`mode: release`**：关闭 Gin 调试输出。
2. **`jwt.secret`**：必须修改为强随机值，否则 token 可被伪造。
3. **SQLite 相对路径**：默认 `./data/app.db` 依赖进程工作目录；用 systemd 已固定 `WorkingDirectory`，容器内需挂卷。多实例请切换为 MySQL。
4. **CORS**：代码里 `AllowAllOrigins: true`，同域（Nginx 反代）下无影响；跨域部署时再按需收紧。
5. **HTTPS**：在 Nginx 前置一层 Let's Encrypt（`certbot`）即可，后端无需改动。

---

## 五、验证

```bash
# 后端健康检查
curl http://localhost:8080/health

# SPA 深链接回退（应返回 200 且内容为 index.html）
curl -I http://your.domain.com/system/role

# 浏览器打开域名，直接刷新 /system/role、/system/menu 应正常，不再 404
```

默认管理员账号：**admin / admin123**（首次启动种子数据写入）。
