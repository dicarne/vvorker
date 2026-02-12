# 数据库连接池配置

## 环境变量说明

为了防止数据库连接超时和 `connection reset by peer` 错误，系统提供了数据库连接池配置。

### 连接池配置环境变量

| 环境变量 | 默认值 | 说明 |
|---------|--------|------|
| `DB_MAX_IDLE_CONNS` | `5` | 最大空闲连接数 |
| `DB_MAX_OPEN_CONNS` | `20` | 最大打开连接数 |
| `DB_CONN_MAX_LIFETIME` | `5` | 连接最大生命周期（分钟）|
| `DB_CONN_MAX_IDLE_TIME` | `1` | 空闲连接超时时间（分钟）|

### 配置说明

这些配置适用于 **PostgreSQL** 和 **MySQL** 数据库。

对于 **SQLite**，系统会自动使用更小的连接池配置：
- 最大空闲连接数：`max(1, DB_MAX_IDLE_CONNS / 2)`
- 最大打开连接数：`max(3, DB_MAX_OPEN_CONNS / 3)`
- 连接最大生命周期：`DB_CONN_MAX_LIFETIME * 2` 分钟
- 空闲连接超时时间：`DB_CONN_MAX_IDLE_TIME * 5` 分钟

### 工作原理

连接池通过以下机制防止连接重置问题：

1. **定期回收连接**：`DB_CONN_MAX_LIFETIME` 设置连接的最大生命周期，超过该时间的连接会被强制回收，下次使用时自动创建新连接
2. **空闲超时处理**：`DB_CONN_MAX_IDLE_TIME` 设置空闲连接的超时时间，超过该时间未使用的连接会被回收
3. **连接池管理**：保持适量的活跃连接，避免频繁创建和销毁连接带来的性能开销

### 推荐配置

对于使用 **frp 隧道**转发数据库连接的场景，建议使用以下配置：

```bash
# 保持连接活跃，减少空闲时间
DB_CONN_MAX_IDLE_TIME=1  # 1分钟后回收空闲连接
DB_CONN_MAX_LIFETIME=5   # 5分钟后强制创建新连接
DB_MAX_IDLE_CONNS=5      # 保持5个空闲连接
DB_MAX_OPEN_CONNS=20     # 最多20个并发连接
```

对于高并发场景，可以适当增加连接数：

```bash
DB_MAX_IDLE_CONNS=10
DB_MAX_OPEN_CONNS=50
DB_CONN_MAX_LIFETIME=10
DB_CONN_MAX_IDLE_TIME=2
```

### 配置示例

在 `.env` 文件中添加以下配置：

```bash
# 数据库连接池配置
DB_MAX_IDLE_CONNS=5
DB_MAX_OPEN_CONNS=20
DB_CONN_MAX_LIFETIME=5
DB_CONN_MAX_IDLE_TIME=1
```

或者在启动应用时设置环境变量：

```bash
export DB_MAX_IDLE_CONNS=5
export DB_MAX_OPEN_CONNS=20
export DB_CONN_MAX_LIFETIME=5
export DB_CONN_MAX_IDLE_TIME=1
```

### 监控和日志

应用启动时会输出连接池配置信息：

```
PostgreSQL database initialized with connection pool: max_idle=5, max_open=20, max_lifetime=5m0s, max_idle_time=1m0s
MySQL database initialized with connection pool: max_idle=5, max_open=20, max_lifetime=5m0s, max_idle_time=1m0s
SQLite database initialized with connection pool: max_idle=3, max_open=7, max_lifetime=10m0s, max_idle_time=5m0s
```

### 注意事项

- 修改这些配置需要重启应用才能生效
- 连接池配置仅在网络不稳定或使用隧道转发时需要调整
- 如果数据库本身运行正常但仍出现连接错误，建议调小 `DB_CONN_MAX_LIFETIME` 和 `DB_CONN_MAX_IDLE_TIME`
- 对于本地数据库（不使用隧道），可以使用默认值或更大的值
