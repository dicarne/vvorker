# Service 绑定

Service 绑定用于调用其他 Worker 服务，支持 HTTP 请求和服务间 RPC 调用。

## 获取绑定

```typescript
import { vvbind } from "@dicarne/vvorker-sdk";

const someWorker = vvbind(c).service("SomeWorker");
```

## 方法

### fetch

向目标服务发送 HTTP 请求。

```typescript
fetch(path: string, init?: RequestInit): Promise<Response>
```

**参数**

| 参数 | 类型 | 描述 |
|------|------|------|
| path | string | 请求路径（会自动添加 `http://vvorker.local` 前缀） |
| init | RequestInit | 请求配置（可选） |

**RequestInit 常用字段**

| 字段 | 类型 | 描述 |
|------|------|------|
| method | string | HTTP 方法（GET、POST、PUT、DELETE 等） |
| headers | object | 请求头 |
| body | string | 请求体 |

**返回值**

- `Promise<Response>` - HTTP 响应对象

**示例**

```typescript
// GET 请求
const res = await someWorker.fetch("/api/users");
const data = await res.json();

// POST 请求
const res = await someWorker.fetch("/api/users", {
    method: "POST",
    headers: {
        "Content-Type": "application/json"
    },
    body: JSON.stringify({ name: "John" })
});
```

---

### call

简化调用方法，直接以 JSON 形式调用服务。

```typescript
call(path: string, data?: any): Promise<any>
```

**参数**

| 参数 | 类型 | 描述 |
|------|------|------|
| path | string | 请求路径 |
| data | any | 请求数据（会被序列化为 JSON） |

**返回值**

- `Promise<any>` - 响应体的 JSON 数据

> 注意：此方法使用 POST 请求，Content-Type 为 application/json

**示例**

```typescript
// 简单调用
const result = await someWorker.call("/api/users", {
    name: "John",
    email: "john@example.com"
});
console.log(result);
```

## 完整示例

```typescript
import { vvbind } from "@dicarne/vvorker-sdk";

// 调用用户服务获取用户信息
honoApi.get("/user/:id/profile", async (c) => {
    const userId = c.req.param("id");
    const userService = vvbind(c).service("UserService");
    
    const res = await userService.fetch(`/api/user/${userId}`);
    
    if (!res.ok) {
        return c.json({ error: "获取用户信息失败" }, 500);
    }
    
    const user = await res.json();
    return c.json(user);
});

// 调用订单服务创建订单
honoApi.post("/order", async (c) => {
    const orderService = vvbind(c).service("OrderService");
    const body = await c.req.json();
    
    try {
        const result = await orderService.call("/api/order/create", {
            userId: body.userId,
            items: body.items
        });
        return c.json(result);
    } catch (e) {
        return c.json({ error: e.message }, 500);
    }
});

// 服务间链式调用
honoApi.get("/order/:id/detail", async (c) => {
    const orderId = c.req.param("id");
    const orderService = vvbind(c).service("OrderService");
    const userService = vvbind(c).service("UserService");
    
    // 获取订单信息
    const order = await orderService.call(`/api/order/${orderId}`);
    
    // 获取用户信息
    const user = await userService.call(`/api/user/${order.userId}`);
    
    return c.json({
        order,
        user
    });
});
```

## 开发模式说明

在本地开发模式下，service 绑定会通过代理访问远程节点的服务。此时会自动添加 `vvorker-worker-uid` 请求头，用于身份验证。

```typescript
// 本地开发时，可以通过环境变量指定模拟的用户 UID
// .env 文件
VITE_APP_UID=your-test-uid
```
