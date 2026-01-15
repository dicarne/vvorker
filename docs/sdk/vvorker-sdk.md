# VVORKER SDK

本 SDK 用于在本地调试时提供节点中绑定服务的代理服务，由于相关的数据库往往只能由主节点进行访问，本地调试无法访问，然而搭建一整套相关的服务又较为复杂，因此本 SDK 应运而生。

## 安装

```bash
pnpm install @dicarne/vvorker-sdk
```

## 使用

### 配置环境变量

在项目根目录下创建`.env`文件，用于配置环境变量。其中，URL 为需要连接的测试环境的对应服务地址。（不是平台根路径）

```env
VITE_VVORKER_BASE_URL=http://YOUR-DOMAIN/target-service-name/
```

接下来复制生产环境的配置文件，修改名称、uid、绑定的资源 id。

配置文件命名规则：`vvorker.环境名.json`。

每个环境将自动使用对应名称的配置文件，请勿搞混。

### 注册调试路由

在使用 sdk 之前，需要在 hono 主路由中对代理进行注册。

```typescript
import { useDebugEndpoint, init } from "@dicarne/vvorker-sdk";
import { env } from "cloudflare:workers";

init(env); // 必须，读取环境变量

useDebugEndpoint(app); // 启用调试路由，在产品模式下不会暴露，无需删除。
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

最后需要发布到节点，否则无法访问节点资源。但只需要在配置文件、环境变更后才需发布一次，本地代码变更则不需要进行发布。

> [!NOTE]
> 使用`vvcli env`切换到包含调试模式的环境后，再进行发布。

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
    var some_worker = await vvbind(c).service("SomeWorker")
    var res = await some_worker.fetch("/someapi", {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        }
    })
    // res 返回json，不需要手动转换成json。并且也不支持其他类型的响应。
    ...
});
```
