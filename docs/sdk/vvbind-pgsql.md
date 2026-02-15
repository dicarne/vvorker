# PostgreSQL 绑定

PostgreSQL 绑定提供 PostgreSQL 数据库访问功能，支持连接信息获取、SQL 查询执行和客户端模式。

## 获取绑定

```typescript
import { vvbind } from "@dicarne/vvorker-sdk";

const pgsql = vvbind(c).pgsql("pgsqlBindName");
```

## 方法

### connectionString

获取数据库连接字符串。

```typescript
connectionString(): Promise<string>
```

**返回值**

- `Promise<string>` - 数据库连接字符串

**示例**

```typescript
const connStr = await pgsql.connectionString();
console.log(connStr); // postgresql://user:password@host:port/database
```

---

### connectionInfo

获取数据库连接详细信息。

```typescript
connectionInfo(): Promise<{
    user: string,
    host: string,
    database: string,
    password: string,
    port: number
}>
```

**返回值**

- `Promise<object>` - 连接信息对象

| 字段 | 类型 | 描述 |
|------|------|------|
| user | string | 数据库用户名 |
| host | string | 数据库主机地址 |
| database | string | 数据库名称 |
| password | string | 数据库密码 |
| port | number | 数据库端口 |

**示例**

```typescript
const info = await pgsql.connectionInfo();
console.log(`连接到 ${info.host}:${info.port} 的 ${info.database} 数据库`);
```

---

### client

获取数据库客户端实例，用于执行查询。

```typescript
client(): Promise<PGSQLClient>
```

**PGSQLClient 接口**

```typescript
interface PGSQLClient {
    query(sql: string): Promise<{
        rows: any[],
        rowCount: number,
        command: string,
        oid: number
    }>;
}
```

**返回值**

- `Promise<PGSQLClient>` - 数据库客户端

**示例**

```typescript
const client = await pgsql.client();

const result = await client.query("SELECT * FROM users WHERE id = $1", [1]);
console.log("查询到", result.rowCount, "条记录");
console.log(result.rows);
```

---

### query

执行 SQL 查询（代理模式）。

```typescript
query(
    sql: string,
    params: any,
    method: string
): Promise<
    ({ rows: string[] } | { rows: string[][] }) & {
        types: string[];
        code: number;
        msg?: string;
    }
>
```

**参数**

| 参数 | 类型 | 描述 |
|------|------|------|
| sql | string | SQL 查询语句 |
| params | any | 查询参数 |
| method | string | 查询方法 |

**返回值**

- `Promise<object>` - 查询结果

| 字段 | 类型 | 描述 |
|------|------|------|
| rows | string[] \| string[][] | 查询结果行 |
| types | string[] | 字段类型 |
| code | number | 状态码（0 表示成功） |
| msg | string | 错误消息（可选） |

**示例**

```typescript
const result = await pgsql.query(
    "SELECT * FROM users WHERE status = $1",
    ["active"],
    "all"
);

if (result.code === 0) {
    console.log("查询结果:", result.rows);
}
```

## 完整示例

```typescript
import { vvbind } from "@dicarne/vvorker-sdk";

// 使用 client 方式
honoApi.get("/users/:id", async (c) => {
    const pgsql = vvbind(c).pgsql("pgsql");
    const userId = c.req.param("id");
    
    const client = await pgsql.client();
    const result = await client.query(
        `SELECT id, name, email FROM users WHERE id = ${userId}`
    );
    
    if (result.rowCount === 0) {
        return c.json({ error: "用户不存在" }, 404);
    }
    
    return c.json(result.rows[0]);
});

// 使用 query 方式
honoApi.post("/users", async (c) => {
    const pgsql = vvbind(c).pgsql("pgsql");
    const body = await c.req.json();
    
    const result = await pgsql.query(
        "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
        [body.name, body.email],
        "all"
    );
    
    if (result.code !== 0) {
        return c.json({ error: result.msg }, 500);
    }
    
    return c.json({ success: true, id: result.rows[0] });
});
```

## 与 ORM 集成

PostgreSQL 绑定同样可以与 Drizzle ORM 等 ORM 库结合使用：

```typescript
import { vvbind } from "@dicarne/vvorker-sdk";
import { drizzle } from "drizzle-orm/postgres-js";

async function getDb(c: Context) {
    const pg = vvbind(c).pgsql("pgsql");
    const client = await pg.client();
    
    return drizzle(client);
}
```
