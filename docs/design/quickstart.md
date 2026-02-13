# 快速开始

本指南将带你从零开始完成 VVorker 平台的部署、配置、开发和部署全流程。

## 1. 平台部署

### 环境要求

- Go 1.25.4+
- Docker（推荐）
- Docker Compose
- Node.js 18+（用于前端开发）
- 域名（可选，用于生产环境）

### 部署方式

#### 使用 Docker Compose（推荐）

VVorker 支持两种运行模式：**主节点（master）** 和 **子节点（agent）**。

```bash
# 编写 docker-compose.yml 文件

# 配置环境变量
# 请查看下一章节的详细说明

# 启动服务（默认主节点模式）
docker-compose up -d
```

docker-compose.yml 示例

```yaml
services:
  vvorker:
    image: git.cloud.zhishudali.ink/dicarne/vvorker:latest
    restart: always
    volumes:
      - ./vvorker-data:/app/data
      - /etc/localtime:/etc/localtime:ro
      - /etc/timezone:/etc/timezone:ro
    ports:
      - "8888:8888"
      - "8080:8080"
      - "18080:18080"
      - "10080:10080"
    depends_on:
      vvorker-pgsql:
        condition: service_healthy
  redis:
    image: git.cloud.zhishudali.ink/dicarne/redis:8
    restart: always
    ports:
      - "6379:6379"
  mysql:
    image: git.cloud.zhishudali.ink/dicarne/mysql:8.0.39
    restart: always
    volumes:
      - ./mysql:/var/lib/mysql
    environment:
      - MYSQL_ROOT_PASSWORD=<YOUR_ROOT_PASSWORD>
    ports:
      - 3306:3306

```

#### 仅启动主节点

如果只需要单节点部署，可以只启动主节点：

```bash
docker-compose up -d
```

#### 访问管理界面

服务启动后，访问管理控制台：

```
http://localhost:8888/admin
```

首次访问需要注册管理员账户，注册后可登录。

> [!TIP]
> 默认注册功能是开启的。如需禁用公开注册，仅允许管理员创建用户，请设置环境变量 `ENABLE_REGISTER=false`。

## 2. 配置平台环境变量

VVorker 通过环境变量进行灵活配置。以下是最关键的环境变量配置：

### 基础配置

在 `docker-compose.yml` 或 `.env` 文件中配置：

| 环境变量 | 说明 | 默认值 | 示例 |
|---------|------|--------|------|
| `RUN_MODE` | 运行模式 | `master` | `master` / `agent` |
| `WORKER_URL_SUFFIX` | Worker URL 后缀 | `.vvorker.local` | `.example.com` |
| `SCHEME` | 使用的协议 | `http` | `http` / `https` |
| `COOKIE_DOMAIN` | Cookie 域名 | `vvorker.local` | `example.com` |
| `JWT_SECRET` | JWT 签名密钥 | `secret` | `your-secret-key` |
| `AGENT_SECRET` | 节点间通信密钥 | - | `123123` |
| `ENABLE_REGISTER` | 是否允许公开注册 | `true` | `true` / `false` |

### 数据库配置（可选）

VVorker 默认使用 SQLite，也可配置使用 MySQL 或 PostgreSQL：

#### MySQL 配置

```bash
ENABLE_MYSQL=true
SERVER_MYSQL_HOST=localhost
SERVER_MYSQL_PORT=3306
SERVER_MYSQL_USER=root
SERVER_MYSQL_PASSWORD=your-password
```

#### PostgreSQL 配置

```bash
ENABLE_PGSQL=true
SERVER_POSTGRE_HOST=localhost
SERVER_POSTGRE_PORT=5432
SERVER_POSTGRE_USER=postgres
SERVER_POSTGRE_PASSWORD=your-password
```

### KV 存储配置（可选）

默认使用 Redis：

```bash
ENABLE_REDIS=true
SERVER_REDIS_HOST=localhost
SERVER_REDIS_PORT=6379
SERVER_REDIS_PASSWORD=your-password
```

### OSS 存储配置（可选）

支持 Minio 或阿里云 OSS：

#### Minio 配置

```bash
ENABLE_MINIO=true
SERVER_MINIO_HOST=localhost
SERVER_MINIO_PORT=9000
SERVER_MINIO_ACCESS=minioadmin
SERVER_MINIO_SECRET=minioadmin
```

#### 阿里云 OSS 配置

```bash
SERVER_OSS_TYPE=aliyun
# 其他阿里云 OSS 相关配置...
```

### 端口映射

| 端口 | 说明 |
|------|------|
| `8080` | Worker 服务端口（反向代理） |
| `8888` | 管理 API 和 Web UI 端口 |
| `10080` | Tunnel 入口端口（内部） |
| `18080` | Tunnel API 端口（公开） |


> [!WARNING]
> 生产环境部署时，请务必修改默认密钥（`JWT_SECRET`、`AGENT_SECRET`），并使用 HTTPS 协议。

## 3. 启动 Docker 服务

### 启动服务

```bash
# 启动所有服务
docker compose up -d

# 查看服务状态
docker compose ps

# 查看日志
docker compose logs -f vvorker
```

### 验证服务

```bash
# 检查管理界面是否可访问
curl http://localhost:8888/admin

# 检查 API 是否正常
curl http://localhost:8888/api/health
```

## 4. 创建 Worker

### 访问管理控制台

打开浏览器访问：

```
http://localhost:8888/admin
```

### 注册/登录

1. 首次访问点击"注册"
2. 填写用户名、密码等信息
3. 注册成功后自动登录

### 创建 Worker

1. 点击左侧菜单"Workers"
2. 点击"创建 Worker"按钮
3. 填写必要信息：
   - **名称**：Worker 的唯一标识（如 `myapp`）
   - **描述**：Worker 的说明信息
   - **路由类型**：选择 `host`（基于域名）或 `path`（基于路径）
4. 点击"保存"

### 获取 Worker UID

1. Worker 创建成功后，点击 Worker 名称进入详情页
2. 复制 **UID**（格式类似 `6eb40d7735f7499dae93d40a31912905`）
3. 该 UID 将用于本地项目绑定

### 配置资源绑定

在 VVorker 控制台中，可以绑定以下资源：

- **数据库**：MySQL、PostgreSQL
- **KV 存储**：Redis、NutsDB
- **对象存储**：Minio、阿里云 OSS
- **内部服务**：绑定其他 Worker 进行内部调用

## 5. 本地创建环境

### 安装 CLI 工具

VVorker 提供了 `vvcli` 命令行工具来简化项目开发和部署流程。

> [!IMPORTANT]
> 本项目均使用 `pnpm` 作为包管理器，请勿使用 `npm`！

```bash
# 全局安装 vvcli
pnpm install -g @dicarne/vvcli
```

### 创建 CLI 环境

`vvcli` 支持多环境管理，可以配置不同的平台 URL 和 API 密钥。

```bash
# 创建新环境
vvcli create
```

根据提示输入：

1. **环境名称**：如 `dev`（开发环境）、`test`（测试环境）、`prod`（生产环境）
2. **URL**：VVorker 平台地址
   - 本地开发：`http://localhost:8888`
   - 远程服务器：`https://your-domain.com`
3. **API 密钥**：API 访问令牌
   - 从 VVorker 控制台"用户"页面获取

### 切换环境

```bash
# 切换环境
vvcli env
```

## 6. 本地开发

### 初始化项目

```bash
# 初始化新项目
vvcli init <项目名>
```

执行过程中需要输入：

- **平台 UID**：粘贴从控制台复制的 Worker UID
- **项目类型**：选择项目类型
  - **Worker**：完整的后端项目，包含数据库、KV、OSS 等资源
  - **Vue**：前后端分离项目（前端 + 后端），自动配置 Assets 绑定
  - **Simple**：简单 Worker 项目，不包含数据库

### 安装依赖

```bash
# 进入项目目录
cd <项目名>

# 安装依赖
pnpm i
```

### 启动开发服务器

```bash
# 启动开发模式
vvcli dev
```

该命令等同于 `pnpm run dev`。

### 生成类型提示

当修改了 `vvorker.{env}.json` 配置文件后，需要重新生成类型提示：

```bash
# 生成类型文件
vvcli types
```

生成的文件：
- `src/binding.ts` 或 `server/binding.ts` - 类型定义文件
- `vvorker-schema.json` - 配置 Schema 文件

### 本地开发调试

本地开发时，SDK 可以连接远程 VVorker 平台绑定的资源：

```typescript
// src/index.ts
export default {
  async fetch(request: Request, env: Env) {
    // 访问绑定的 KV 存储
    const value = await env.MYKV.get('key');

    // 访问绑定的数据库
    const result = await env.MYSQL.query('SELECT * FROM users');

    // 调用绑定的内部服务
    const response = await env.MyService.fetch(request);

    return new Response('Hello from VVorker!');
  }
};
```

> [!WARNING]
> 本地开发时，请连接开发环境而不是生产环境，以避免误操作生产数据！



## 7. 部署到服务器

### 更新版本号（可选）

```bash
# 更新修订号（1.2.3 -> 1.2.4）
vvcli ver

# 更新次版本号（1.2.3 -> 1.3.0）
vvcli ver -s

# 更新主版本号（1.2.3 -> 2.0.0）
vvcli ver -m
```

### 部署到平台

```bash
# 部署到当前环境
vvcli deploy
```

该命令会：
1. 运行构建命令（`pnpm run build`）
2. 上传 Assets 文件（增量差分，跳过未修改文件）
3. 执行数据库迁移（MySQL/PostgreSQL）
4. 上传 Worker 代码
5. 等待部署完成并显示结果

### 部署选项

```bash
# 正常部署（完整流程）
vvcli deploy

# 跳过构建（如果已手动构建）
vvcli deploy --skip-build

# 跳过 Assets 上传
vvcli deploy --skip-assets

# 强制上传所有 Assets 文件（跳过差分校验）
vvcli deploy --force
```

### 查看线上日志

```bash
# 查看最新日志
vvcli logs

# 实时追踪日志
vvcli logs -f

# 查看第 2 页日志
vvcli logs -p 2
```

### 访问部署的 Worker

#### 通过 Host 头访问（推荐用于测试）

```bash
curl localhost:8080 -H "Host: workername.yourdomain.com"
```

#### 通过 Server-Host 头访问

```bash
curl localhost:8080 -H "Server-Host: workername.yourdomain.com"
```

#### 通过域名访问（需要配置 DNS）

```bash
curl https://workername.yourdomain.com
```

### 多环境部署

```bash
# 1. 创建开发环境
vvcli create
# 输入环境名称: dev
# 输入 URL: http://localhost:8888
# 输入 Token: dev-token

# 2. 创建生产环境
vvcli create
# 输入环境名称: prod
# 输入 URL: https://your-domain.com
# 输入 Token: prod-token

# 3. 部署到开发环境
vvcli env  # 选择 dev
vvcli deploy

# 4. 部署到生产环境
vvcli env  # 选择 prod
vvcli deploy
```

## 生态工具

### CLI 工具（vvcli）

命令行工具，用于快速部署和管理 Worker：

- 创建本地项目
- 部署代码
- 切换多个开发环境
- 生成类型
- 查看日志

### SDK（TypeScript）

提供 TypeScript/JavaScript SDK，便于在项目中集成 VVorker：

- API 调用封装
- 资源管理接口
