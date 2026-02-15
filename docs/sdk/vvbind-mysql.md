# MySQL 绑定

MySQL 绑定提供 MySQL 数据库访问功能，支持连接信息获取和 SQL 查询执行。

## 获取绑定

```typescript
import { vvbind } from "@dicarne/vvorker-sdk";

const mysql = vvbind(c).mysql("mysqlBindName");
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
const connStr = await mysql.connectionString();
console.log(connStr); // mysql://user:password@host:port/database
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
const info = await mysql.connectionInfo();
console.log(`连接到 ${info.host}:${info.port} 的 ${info.database} 数据库`);
```

---

### query

执行 SQL 查询。

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
| method | string | 查询方法（如 `"all"`, `"get"`, `"execute"`） |

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
// 查询所有记录
const result = await mysql.query(
    "SELECT * FROM users WHERE status = ?",
    ["active"],
    "all"
);

if (result.code === 0) {
    console.log("查询结果:", result.rows);
}

// 插入记录
const insertResult = await mysql.query(
    "INSERT INTO users (name, email) VALUES (?, ?)",
    ["John", "john@example.com"],
    "execute"
);
```

## 与 Drizzle ORM 集成

MySQL 绑定可以与 Drizzle ORM 结合使用，提供更好的类型安全性和开发体验。

```typescript
import { vvbind } from "@dicarne/vvorker-sdk";
import { drizzle } from "drizzle-orm/mysql-proxy";

async function getDb(c: Context) {
    const mdb = vvbind(c).mysql("mysql");
    
    return drizzle(async (sql, params, method) => {
        try {
            const rows: any = await mdb.query(
                sql,
                params.map((a) => (typeof a === "number" ? String(a) : a)),
                method
            );
            
            if (rows.code) {
                console.log({ sql, params, method });
                console.error(rows.msg);
                throw new Error(rows.msg);
            }
            return rows;
        } catch (e: any) {
            console.error("MySQL 代理服务器错误：", e);
            return { rows: [] };
        }
    });
}
```

## 完整示例

```typescript
import { vvbind } from "@dicarne/vvorker-sdk";

honoApi.get("/users", async (c) => {
    const mysql = vvbind(c).mysql("mysql");
    
    const result = await mysql.query(
        "SELECT id, name, email FROM users LIMIT 10",
        [],
        "all"
    );
    
    if (result.code !== 0) {
        return c.json({ error: result.msg }, 500);
    }
    
    return c.json({ users: result.rows });
});

honoApi.post("/users", async (c) => {
    const mysql = vvbind(c).mysql("mysql");
    const body = await c.req.json();
    
    const result = await mysql.query(
        "INSERT INTO users (name, email, created_at) VALUES (?, ?, NOW())",
        [body.name, body.email],
        "execute"
    );
    
    if (result.code !== 0) {
        return c.json({ error: result.msg }, 500);
    }
    
    return c.json({ success: true });
});
```
