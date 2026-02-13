# VVorker CLI

Vvorker CLI 是用于管理 VVorker 平台项目的命令行工具，提供项目初始化、部署、环境管理等功能。

## 安装

```bash
pnpm install -g @dicarne/vvcli
```

## 命令列表

### `init` - 初始化项目

创建一个新的 VVorker 项目，并自动配置相关文件。

```bash
vvcli init <projectName>
```

该命令会：
1. 交互式输入 VVorker 平台的 Worker UID
2. 选择项目类型：
   - **Worker**：复杂的后端项目，包含完整功能
   - **Vue**：前后端分离项目（前端 + 后端）
   - **Simple**：简单 Worker 项目，不包含数据库
3. 从模板仓库克隆项目代码
4. 自动生成 `vvorker.{env}.json` 配置文件
5. 安装依赖、初始化 Git、生成类型提示

#### 项目类型说明

| 类型 | 说明 | 适用场景 |
|------|------|----------|
| Worker | 完整后端项目 | 包含数据库、KV、OSS 等资源的复杂业务 |
| Vue | 前后端分离 | 需要 UI 界面的应用，自动配置 Assets 绑定 |
| Simple | 简单 Worker | 轻量级服务，无需数据库支持 |

### `deploy` - 部署项目

将项目部署到 VVorker 平台。

```bash
vvcli deploy [options]
```

该命令会：
1. 读取 `vvorker.{env}.json` 配置
2. 运行构建命令（`pnpm run build` 或 `npm run build`）
3. 上传 Assets 文件（增量差分，跳过未修改文件）
4. 执行数据库迁移（MySQL/PostgreSQL）
5. 上传 Worker 代码
6. 等待部署完成并显示结果

#### 选项

| 选项 | 说明 |
|------|------|
| `--skip-assets` | 跳过 Assets 文件上传 |
| `--skip-build` | 跳过构建步骤 |
| `-f, --force` | 强制上传所有 Assets 文件（跳过差分校验） |

#### 示例

```bash
# 正常部署（完整流程）
vvcli deploy

# 跳过构建（如果已手动构建）
vvcli deploy --skip-build

# 强制上传所有资源文件
vvcli deploy --force
```

### `types` - 生成类型提示

根据 `vvorker.{env}.json` 配置自动生成 TypeScript 类型定义。

```bash
vvcli types
```

生成的文件：
- `src/binding.ts` 或 `server/binding.ts` - 类型定义文件
- `vvorker-schema.json` - 配置 Schema 文件

### `dev` - 开发模式

启动开发服务器，等同于运行 `pnpm run dev`。

```bash
vvcli dev
```

### `logs` - 查看日志

获取 Worker 运行日志。

```bash
vvcli logs [options]
```

#### 选项

| 选项 | 说明 |
|------|------|
| `-f, --follow` | 实时显示日志（类似 `tail -f`） |
| `-p, --page <number>` | 分页页码（从 1 开始），默认 1 |
| `--page-size <number>` | 每页记录数，默认 50 |

#### 示例

```bash
# 查看最新日志（默认 50 条）
vvcli logs

# 查看第 2 页日志
vvcli logs -p 2

# 实时追踪日志
vvcli logs -f
```

### `create` - 创建环境

交互式创建新的部署环境配置。

```bash
vvcli create
```

需要输入：
- **环境名称**：如 `dev`、`test`、`prod`
- **URL**：VVorker 平台地址，如 `https://vvorker.cloud.zhishudali.ink`
- **Token**：API 访问令牌

> 首次创建的环境会自动设置为当前环境。

### `env` - 切换环境

交互式切换当前使用的部署环境。

```bash
vvcli env
```

列出所有已配置的环境，选择后切换。

### `config` - 配置管理

管理环境配置参数。

#### `config set` - 设置配置项

```bash
vvcli config set <key> <value>
```

支持的配置项：

| key | 说明 | 示例 |
|-----|------|------|
| `env` | 设置当前环境 | `vvcli config set env prod` |
| `url` | 设置当前环境的 URL | `vvcli config set url https://example.com` |
| `token` | 设置当前环境的 Token | `vvcli config set token xxxxxx` |

#### `config show` - 显示当前配置

```bash
vvcli config show
```

显示：
- 当前环境名称
- Token（部分隐藏）
- URL

### `ver` - 版本管理

更新项目版本号。

```bash
vvcli ver [options]
```

#### 选项

| 选项 | 说明 | 版本变化示例 |
|------|------|--------------|
| `-m, --major` | 主版本号+1，次版本和修订归零 | `1.2.3` → `2.0.0` |
| `-s, --minor` | 次版本号+1，修订归零 | `1.2.3` → `1.3.0` |
| （无） | 修订号+1 | `1.2.3` → `1.2.4` |

版本号会写入项目根目录的 `VERSION.txt` 文件，部署时自动更新到配置中。

#### 示例

```bash
# 增加修订号（1.2.3 -> 1.2.4）
vvcli ver

# 增加次版本号（1.2.3 -> 1.3.0）
vvcli ver -s

# 增加主版本号（1.2.3 -> 2.0.0）
vvcli ver -m
```

### `upgrade` - 升级 CLI

升级 vvcli 到最新版本。

```bash
vvcli upgrade
```

---

## 环境配置说明

配置文件位置：`~/.vvcli/config.json`

配置结构：
```json
{
  "current_env": "dev",
  "env": {
    "dev": {
      "url": "https://vvorker.cloud.zhishudali.ink",
      "token": "your_token_here"
    },
    "prod": {
      "url": "https://production.example.com",
      "token": "prod_token_here"
    }
  }
}
```

详细配置说明请参考 [CLI 环境配置](/config/env)。

## 工作流程示例

### 典型开发流程

```bash
# 1. 创建环境
vvcli create

# 2. 初始化项目
vvcli init my-project

# 3. 开发
vvcli dev

# 4. 查看日志
vvcli logs -f

# 5. 更新版本并部署
vvcli ver -s
vvcli deploy
```

### 多环境部署

```bash
# 创建开发环境
vvcli create  # 输入 dev 环境信息
vvcli init my-project  # 生成 vvorker.dev.json
vvcli deploy  # 部署到 dev

# 切换到生产环境
vvcli env  # 选择 prod
# 或
vvcli config set env prod

# 部署到生产环境
vvcli deploy
```
