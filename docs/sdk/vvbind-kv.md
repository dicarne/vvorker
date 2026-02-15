# KV 绑定

KV 绑定提供键值存储功能，支持基本的 CRUD 操作和模式匹配查询。

## 获取绑定

```typescript
import { vvbind } from "@dicarne/vvorker-sdk";

const kv = vvbind(c).kv("kvBindName");
```

## 方法

### get

获取指定键的值。

```typescript
get(key: string): Promise<string | null>
```

**参数**

| 参数 | 类型 | 描述 |
|------|------|------|
| key | string | 要获取的键名 |

**返回值**

- `Promise<string | null>` - 返回键对应的值，如果键不存在则返回 `null`

**示例**

```typescript
const value = await kv.get("mykey");
if (value) {
    console.log(value);
}
```

---

### set

设置键值对。

```typescript
set(key: string, value: string, options?: {
    EX?: number,
    NX?: boolean,
    XX?: boolean
} | number): Promise<number>
```

**参数**

| 参数 | 类型 | 描述 |
|------|------|------|
| key | string | 键名 |
| value | string | 键值 |
| options | object \| number | 可选配置 |

**options 参数说明**

| 参数 | 类型 | 描述 |
|------|------|------|
| EX | number | 设置过期时间（秒） |
| NX | boolean | 仅当键不存在时设置 |
| XX | boolean | 仅当键存在时设置 |

> 也可以直接传入一个数字作为过期时间（秒）

**返回值**

- `Promise<number>` - 操作结果

**示例**

```typescript
// 简单设置
await kv.set("mykey", "myvalue");

// 设置过期时间（10秒）
await kv.set("mykey", "myvalue", 10);

// 仅当键不存在时设置
await kv.set("mykey", "myvalue", { NX: true });

// 设置过期时间并仅当键存在时更新
await kv.set("mykey", "myvalue", { EX: 60, XX: true });
```

---

### del

删除指定键。

```typescript
del(key: string): Promise<void>
```

**参数**

| 参数 | 类型 | 描述 |
|------|------|------|
| key | string | 要删除的键名 |

**返回值**

- `Promise<void>`

**示例**

```typescript
await kv.del("mykey");
```

---

### keys

根据模式匹配获取键列表。

```typescript
keys(pattern: string, offset: number, size: number): Promise<string[]>
```

**参数**

| 参数 | 类型 | 描述 |
|------|------|------|
| pattern | string | 匹配模式（支持通配符 `*`） |
| offset | number | 偏移量 |
| size | number | 返回数量 |

**返回值**

- `Promise<string[]>` - 匹配的键名数组

**示例**

```typescript
// 获取所有以 "user:" 开头的键
const keys = await kv.keys("user:*", 0, 100);
console.log(keys); // ["user:1", "user:2", ...]
```

## 完整示例

```typescript
import { vvbind } from "@dicarne/vvorker-sdk";

honoApi.get("/cache/user/:id", async (c) => {
    const kv = vvbind(c).kv("mykv");
    const userId = c.req.param("id");

    // 尝试从缓存获取
    const cacheKey = `user:${userId}`;
    let userData = await kv.get(cacheKey);

    if (!userData) {
        // 缓存不存在，从数据库获取
        userData = JSON.stringify({ id: userId, name: "John" });

        // 存入缓存，设置60秒过期
        await kv.set(cacheKey, userData, 60);
    }

    return c.json(JSON.parse(userData));
});
```
