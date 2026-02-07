# 快速开始

本指南将带你快速上手使用 VVorker 平台进行项目开发和部署。


## 0. 平台部署

### 环境要求

- Go 1.25.4+
- Docker（推荐）
- Node.js 18+（用于前端开发）

### 快速部署

#### 使用 Docker Compose（推荐）

```bash
# 克隆项目
git clone https://github.com/dicarne/vvorker.git
cd vvorker

# 启动服务
docker-compose up -d

# 访问管理界面
# http://localhost:8888/admin
```

#### 访问 Worker

```bash
# 方式 1：通过 Host 头
curl localhost:8080 -H "Host: workername.yourdomain.com"

# 方式 2：通过 Server-Host 头
curl localhost:8080 -H "Server-Host: workername.yourdomain.com"

# 方式 3：通过 HTTPS（配置了域名后）
curl https://workername.yourdomain.com
```

## 生态工具

### CLI 工具（vvcli）

命令行工具，用于快速部署和管理 Worker：

- 创建本地项目
- 部署代码
- 切换多个开发环境
- 生成类型

### SDK（TypeScript）

提供 TypeScript/JavaScript SDK，便于在项目中集成 VVorker：

- Worker 代码打包
- API 调用封装
- 资源管理接口


## 1. 创建 Worker

### 访问管理控制台

首先打开 VVorker 平台管理控制台，默认地址为：

```
http://platform-url/admin
```

### 创建 Worker

1. 点击"创建"按钮创建新的 Worker
2. 填写必要的 Worker 信息
3. 点击"编辑"按钮进入 Worker 编辑界面
4. 复制 **UID**（该 UID 将用于本地项目绑定）

## 2. 安装 CLI 工具

VVorker 提供了 `vvcli` 命令行工具来简化项目开发和部署流程。

> [!IMPORTANT]
> 本项目均使用 `pnpm` 作为包管理器，请勿使用 `npm`！

```bash
pnpm install -g @dicarne/vvcli
```

## 3. 初始化项目

在本地创建新项目：

```bash
vvcli init <项目名>
```

执行过程中需要输入：

- **平台 UID**：粘贴从控制台复制的 UID
- **项目类型**：选择 `vue` 或 `api`
  - `vue`：前后端分离项目（包含前端网页）
  - `api`：纯后端 API 项目

## 4. 安装依赖

进入项目目录并安装依赖：

```bash
cd <项目名>
pnpm i
```

## 5. 本地开发

启动开发服务器进行本地开发和调试：

```bash
vvcli dev
```

该命令等同于 `pnpm run dev`，便于统一命令行操作。

在进行[一些配置](../sdk/vvorker-sdk.md)后，本地 SDK 可以连接使用远程 VVorker 平台绑定的资源。注意，请连接开发环境而不是生产环境以避免不必要的错误！

## 6. 发布部署

开发完成后，使用以下命令将项目部署到 VVorker 平台：

```bash
vvcli deploy
```

该命令会自动编译代码、打包静态资源并上传到平台。如果有使用 drizzle 数据库，也将自动完成数据库迁移工作。

## 配置说明

### 环境管理

`vvcli` 支持多环境管理：

- `vvcli create`：交互式创建新环境
- `vvcli env`：切换当前环境

## 完整示例

```bash
# 1. 全局安装 vvcli
pnpm install -g @dicarne/vvcli

# 2. 初始化项目
vvcli init my-project
# 输入 UID: abc123def456
# 选择类型: vue

# 3. 安装依赖
cd my-project
pnpm i

# 4. 本地开发
vvcli dev

# 5. 部署发布
vvcli deploy
```
