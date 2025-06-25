## VVORKER SDK

本SDK用于在本地调试时提供节点中绑定服务的代理服务，由于相关的数据库往往只能由主节点进行访问，本地调试无法访问，然而搭建一整套相关的服务又较为复杂，因此本SDK应运而生。

## 安装

```bash
pnpm install @dicarne/vvorker-sdk
```

## 使用


### 注册调试路由
在使用sdk之前，需要在hono主路由中对代理进行注册。

```typescript
import { useDebugEndpoint, init } from "@dicarne/vvorker-sdk";
import { env } from "cloudflare:workers"

init(env) // 必须，读取环境变量

useDebugEndpoint(app)  // 启用调试路由，在产品模式下不会暴露
```

接着需要配置环境变量，以启用服务端的代理功能。
```json
// vvorker.json
{
    ...
    "vars": {
        "MODE": "development"
    }
    ...
}

```

最后需要发布到节点，否则无法访问节点资源。但只需要在环境变量变更后发布，本地代码变更则不需要发布。
```bash
vvcli deploy
```

### 类型生成
可以根据`vvorker.json`中的绑定信息生成`typescript`类型文件。
```
vvcli types
```
生成的文件在`server/src/binding.ts`(vite)下。

### 访问绑定

进行本地调试，可以快速重载代码变更。
```bash
pnpm run dev
```

通过`vvbind`访问绑定的资源。由于部分资源只能通过节点访问，因此只能通过`vvbind`代理对节点资源的访问。`vvbind`会根据开发还是产品模式决定是代理还是直接访问绑定资源。

```typescript
import { vvbind } from "@dicarne/vvorker-sdk";

honoApi.get("/someapi", async (c) => {
    var vars = await vvbind(c).vars()
    var oss = await vvbind(c).oss("ossBindName")
    var pgsql = await vvbind(c).pgsql("pgsqlBindName")
    var kv = await vvbind(c).kv("kvBindName")

    ...
});
```