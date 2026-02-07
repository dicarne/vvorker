# VVorker 项目介绍

## 项目概述

VVorker 是一个简单强大的自部署 Cloudflare Worker 替代系统。本项目基于 Cloudflare 的 [Workerd](https://github.com/cloudflare/workerd) 项目开发。

### 核心定位

VVorker 定位为**企业级边缘计算与 Serverless 平台**，专为内部网络环境和私有化部署场景设计，提供了完整的 Worker 生命周期管理、多租户支持、分布式部署等能力。

## 核心功能

### 1. Worker 管理

- **完整的 CRUD 操作**：创建、读取、更新、删除 Worker
- **多种路由方式**：支持基于域名或路径的多 Workers 路由
- **Worker 模板系统**：内置常用 Worker 模板，快速部署

### 2. 多租户与权限

- **多租户支持**：不同用户/租户之间资源隔离
- **API 控制接口**：通过 RESTful API 管理所有资源
- **Access Token 管理**：为 Worker 创建访问令牌
- **单点登录（SSO）**：网关级的单点登录鉴权支持

### 3. 数据存储集成

- **结构化数据库支持**：
  - PostgreSQL
  - MySQL
  - SQLite
- **KV 缓存支持**：Redis
- **对象存储支持**：
  - Minio（兼容 S3）
  - 阿里云 OSS
- **智能资源绑定**：快速绑定内部数据库资源，无需手动管理 AccessKey/SecretKey，即插即用
- **兼容模式**：支持单数据库和单 Bucket 的兼容模式
- **SQL 变更管理**：对 SQL 的变更支持

### 4. 分布式架构

- **多节点支持**：支持主节点和子节点的分布式部署
- **负载均衡**：单节点多进程负载均衡，支持高并发
- **节点管理**：统一的节点管理和监控
- **代理访问**：子节点通过代理访问主节点数据库

### 5. 监控与日志

- **全局日志收集**：统一的日志收集和管理
- **状态监控**：实时查看 Workers 和节点运行状态
- **请求统计**：Worker 请求数量和时间分析

### 6. 开发工具

- **CLI 命令行工具**：`vvcli` 快速部署和管理 Worker
- **SDK 工具包**：提供 TypeScript/JavaScript SDK，便于集成开发
- **Wrangler 兼容**：完美兼容 Cloudflare Wrangler 工具链

### 7. Web 管理界面

- **简洁的在线 UI**：配置资源及代码
- **可视化配置**：直观的 Worker 和资源配置界面
- **实时状态显示**：显示 Worker 运行状态

## 系统架构

### 核心组件说明

#### 1. VVorker Core（主程序）

- **语言**：Go
- **Web 框架**：Gin
- **主要职责**：
  - HTTP API 服务器（端口 8888 / 8080）
  - Worker 生命周期管理
  - 资源管理和绑定
  - 用户鉴权和多租户支持
  - 节点管理和协调
  - 日志收集和监控

#### 2. Workerd Runtime

- **基础**：Cloudflare Workerd
- **职责**：执行 JavaScript/TypeScript Worker 代码
- **特性**：
  - 支持 ES Modules
  - 支持多种绑定（KV、MySQL 等）

#### 3. 前端管理界面

- **技术栈**：TypeScript + Vite + Vue
- **功能**：
  - Worker 在线编辑
  - 资源配置管理
  - 用户和权限管理
  - 监控和日志查看
  - 节点状态展示

#### 4. 扩展模块

- **数据库扩展**：
  - `ext/mysql`：MySQL 支持
  - `ext/pgsql`：PostgreSQL 支持
  - `ext/sqlite`：SQLite 支持
- **KV 存储**：`ext/kv`：Redis KV 存储
- **对象存储**：
  - `ext/oss`：S3 兼容存储（Minio）
  - `ext/oss/alioss`：阿里云 OSS
- **其他扩展**：RPC、Assets、Tunnel 等

## 使用场景

### 1. 内部网络服务部署

在内部网络环境中快速部署和扩展微服务，无需关心服务器资源管理。

### 2. API 网关和代理

使用 Worker 构建灵活的 API 网关，实现请求路由、代理、转换等功能。

### 3. 边缘计算

在分布式节点上部署计算任务，实现边缘计算能力。

### 4. Serverless 应用

无需管理服务器，专注于业务逻辑开发，自动弹性扩展。

### 5. 微服务架构

作为微服务架构的支撑平台，快速创建和管理各类服务。

### 6. Webhook 和事件处理

处理 Webhook 回调和事件驱动的业务逻辑。

## 与 Cloudflare Workers 的区别

| 特性        | Cloudflare Workers       | VVorker                                                             |
| ----------- | ------------------------ | ------------------------------------------------------------------- |
| 部署方式    | SaaS，托管在 Cloudflare  | 自部署，可运行在任意服务器                                          |
| 数据存储    | Cloudflare KV、D1、R2 等 | 支持多种数据库和存储（PostgreSQL、MySQL、Redis、Minio、阿里云 OSS） |
| 多租户      | 是                       | 是                                                                  |
| 分布式      | 全球边缘网络             | 可自定义分布式架构                                                  |
| 成本        | 按使用量付费             | 一次性部署，无额外成本                                              |
| Worker 调试 | 支持                     | 支持                                                                |
