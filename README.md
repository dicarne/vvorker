# VVorker

VVorker 是一个简单强大的自部署 Cloudflare worker 替代系统。本项目基于 [Cloudflare Workerd](https://github.com/cloudflare/workerd) ，并在 [Vorker](https://github.com/VaalaCat/vorker) 的基础上进行改进。由于比 Vorker 多一点功能，因此本项目的名称为 VVorker。由于根据个人需求进行一定程度的魔改，因此难以向上游提交。

## 特色

- [x] 用户鉴权及多租户支持
- [x] 使用 API 控制
- [x] 基于域名或路径的多 Workers 路由
- [x] 简单的在线UI，用于配置资源及代码
- [x] 分布式多节点支持
- [x] litefs(HA) 支持（实验性）
- [x] Cloudflare Durable Objects (实验性)
- [x] 基于 PostgreSQL 的结构化数据库支持
- [x] 基于 Redis 的 KV 缓存支持
- [x] 基于 Minio 的对象存储支持
- [x] 快速绑定内部数据库资源而无需管理AccessKey与SecretKey，即插即用
- [x] 全局日志收集
- [x] 性能与状态监控
- [x] 命令行工具`vvcli`快速部署 Worker
- [ ] 对 SQL 的变更支持
- [ ] Worker 版本控制，包括灰度发布与测试分支
- [ ] Worker Debugging
- [ ] 打包某个服务及其所有依赖，用于迁移到其它系统

## 使用方法

控制面板：`http://localhost:8888/admin`

发出请求：

```bash
curl localhost:8080 -H "Host: workername.yourdomain.com" # replace workername with your worker name
```

或

```bash
curl  https://workername.yourdomain.com # replace workername with your worker name

```
或

```bash
curl localhost:8080 -H "Server-Host: workername.yourdomain.com" # replace workername with your worker name
```

通过环境变量进行配置，可以查看[env.go](./conf/env.go)文件了解更多信息。

## 安全性

本项目拥有一定的安全性，但由于主要目标是内部网络组网，而非 Cloudflare Workers 那样面向全球开发者服务，因此内部安全措施旨在防止开发者误用而非阻止开发者进行攻击。SQL 与 OSS 等数据库、用户、桶将会自动创建并赋予对应权限，KV 则仅通过Key前缀进行区分。考虑到内部网络的特殊性，Cloudflare Workers 默认只允许互联网访问的策略不太合适（内网服务全是本地地址），因此本项目默认开放一切网络访问权限（未来根据需要可能会提供配置方式）。

在设计上，目标内部网络不同服务器间仅有若干端口开放，其他节点无法直接连接数据库，因此数据库只会暴露在主节点中，子节点将通过代理访问数据库。

## Worker 开发
由于本项目基于 Workerd 生态，因此使用`wrangler`几乎是唯一选择（除非是非常简单的代码）。

### 使用 Cloudflare `wrangler`

以下是一些常用命令。

输出打包后的代码：
```bash
wrangler deploy --dry-run --outdir dist
```

`dist/index.js`应该包含所有你的代码，将其拷贝到 VVorker 中的代码编辑区即可，点击保存后自动生效。


## Screenshots

- Admin Page

![](./images/worker-admin.png)

- Worker Editor

![](./images/worker-edit.png)

- Worker Config

![](./images/worker-config.png)

- Agent Status

![](./images/status.png)

- Worker Execution

![](https://vaala.cat/images/vorkerexec.png)

## 其他

感谢 [Vorker](https://github.com/VaalaCat/vorker) 项目提供的良好基础！是一个我看得懂的好项目。

感谢 AI 提供的代码补全，在 GO 语言上工作良好，80%+ 代码由 AI 生成，请原谅其中的废话注释。