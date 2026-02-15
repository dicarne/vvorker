# Proxy 绑定

Proxy 绑定提供网络代理功能，用于控制网络请求的路由方式。

## 获取绑定

```typescript
import { vvbind } from "@dicarne/vvorker-sdk";

const proxy = vvbind(c).proxy("proxyBindName");
```

## 方法

### fetch

通过代理发送网络请求。

```typescript
fetch(url: string, init?: RequestInit): Promise<Response>
```

**参数**

| 参数 | 类型 | 描述 |
|------|------|------|
| url | string | 请求的完整 URL |
| init | RequestInit | 请求配置（可选） |

**返回值**

- `Promise<Response>` - HTTP 响应对象

**示例**

```typescript
const proxy = vvbind(c).proxy("myproxy");

// 发送 GET 请求
const res = await proxy.fetch("https://api.example.com/data");
const data = await res.json();

// 发送 POST 请求
const res = await proxy.fetch("https://api.example.com/submit", {
    method: "POST",
    headers: {
        "Content-Type": "application/json"
    },
    body: JSON.stringify({ key: "value" })
});
```

## 使用场景

### 1. 访问外部 API

通过代理访问外部 API，可以绕过 CORS 限制或隐藏真实服务器地址。

```typescript
honoApi.get("/external-data", async (c) => {
    const proxy = vvbind(c).proxy("externalProxy");
    
    const res = await proxy.fetch("https://external-api.example.com/data", {
        headers: {
            "Authorization": "Bearer secret-key"
        }
    });
    
    return c.json(await res.json());
});
```

### 2. 本地开发模式

在本地开发模式下，proxy 的行为取决于配置：

- 如果配置为远程代理，请求会通过节点转发
- 如果配置为本地代理，直接使用标准 fetch

```typescript
// 本地开发时，某些请求可能需要通过节点代理
const proxy = vvbind(c).proxy("myproxy");
const res = await proxy.fetch("https://internal-service/data");
```

## 完整示例

```typescript
import { vvbind } from "@dicarne/vvorker-sdk";

// 代理转发外部请求
honoApi.get("/proxy/*", async (c) => {
    const proxy = vvbind(c).proxy("myproxy");
    const targetUrl = c.req.path.replace("/proxy/", "");
    
    const res = await proxy.fetch(`https://${targetUrl}`);
    
    return res;
});

// 调用第三方 API
honoApi.post("/send-sms", async (c) => {
    const proxy = vvbind(c).proxy("smsProxy");
    const body = await c.req.json();
    
    const res = await proxy.fetch("https://sms-provider.example.com/send", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "X-API-Key": "your-api-key"
        },
        body: JSON.stringify({
            phone: body.phone,
            message: body.message
        })
    });
    
    if (!res.ok) {
        return c.json({ error: "发送失败" }, 500);
    }
    
    return c.json(await res.json());
});
```

## 配置说明

在 `vvorker.json` 中配置 Proxy 绑定：

```json
{
    "proxy": {
        "myproxy": {
            "type": "proxy",
            "remote": true
        }
    }
}
```

| 字段 | 类型 | 描述 |
|------|------|------|
| type | string | 绑定类型，固定为 `proxy` |
| remote | boolean | 是否使用远程代理（本地开发时通过节点转发） |
