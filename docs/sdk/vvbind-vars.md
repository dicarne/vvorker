# Vars 绑定

Vars 绑定用于访问环境变量，提供统一的配置管理方式。

## 获取绑定

```typescript
import { vvbind } from "@dicarne/vvorker-sdk";

const vars = await vvbind(c).vars();
```

## 说明

`vars()` 方法返回 `vvorker.json` 配置文件中定义的所有环境变量。

**配置示例**

```json
// vvorker.json
{
    "vars": {
        "API_KEY": "your-api-key",
        "MAX_RETRY": 3,
        "DEBUG_MODE": true
    }
}
```

**使用示例**

```typescript
import { vvbind } from "@dicarne/vvorker-sdk";

honoApi.get("/config", async (c) => {
    const vars = await vvbind(c).vars();
    
    return c.json({
        apiKey: vars.API_KEY,
        maxRetry: vars.MAX_RETRY,
        debugMode: vars.DEBUG_MODE
    });
});
```

## 与直接访问 env.vars 的区别

在开发模式下，`vvbind(c).vars()` 会通过代理从远程节点获取环境变量，而直接访问 `c.env.vars` 只能获取本地配置的变量。

```typescript
// 推荐：在开发模式下也能正确获取远程节点的变量
const vars = await vvbind(c).vars();

// 不推荐：开发模式下无法获取远程节点的变量
const vars = c.env.vars;
```

## 完整示例

```typescript
import { vvbind } from "@dicarne/vvorker-sdk";

// 创建配置类
class AppConfig {
    private vars: any;
    
    async init(c: Context) {
        this.vars = await vvbind(c).vars();
    }
    
    get apiKey(): string {
        return this.vars.API_KEY || "";
    }
    
    get maxRetry(): number {
        return Number(this.vars.MAX_RETRY) || 3;
    }
    
    get debugMode(): boolean {
        return this.vars.DEBUG_MODE === true || this.vars.DEBUG_MODE === "true";
    }
}

// 在路由中使用
honoApi.get("/check-api", async (c) => {
    const vars = await vvbind(c).vars();
    
    // 使用环境变量配置第三方服务
    const apiKey = vars.EXTERNAL_API_KEY;
    
    const response = await fetch("https://api.example.com/data", {
        headers: {
            "Authorization": `Bearer ${apiKey}`
        }
    });
    
    return c.json(await response.json());
});

// 多环境配置
honoApi.post("/send-notification", async (c) => {
    const vars = await vvbind(c).vars();
    
    const notificationService = vars.NOTIFICATION_SERVICE_URL;
    const notificationKey = vars.NOTIFICATION_API_KEY;
    
    await fetch(`${notificationService}/send`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${notificationKey}`
        },
        body: JSON.stringify({
            message: "Hello from vvorker"
        })
    });
    
    return c.json({ success: true });
});
```
