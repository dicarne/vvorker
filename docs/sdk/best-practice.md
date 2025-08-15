# 最佳实践

## 一、创建项目

```bash
vvcli init myproject
> ENTER VVORKER UID
```

选择所需模板后，等待片刻即可完成安装。

新项目默认使用`hono`作为后端，`vue3`作为前端。

## 二、配置tsconfig

我们需要使用`hono`rpc模式共享前后端类型，因此需要配置tsconfig。

首先，需要启用`strict`模式，这样`hono`前端才能获得正确的`input`类型。

其次，需要在前端访问后端的类型代码，设置include或references等选项达成此目的。

## 三、编写后端代码

首先我们可以创建一个上下文类，用于存储所需的请求信息，如用户ID、数据库等。

```typescript
import type { Context } from "hono";
import type { EnvBinding } from "../binding";
import { vvbind } from "@dicarne/vvorker-sdk";
import { drizzle } from 'drizzle-orm/mysql-proxy';
import { ForbiddenError } from "./error"

export class VContext {
    private db: ReturnType<typeof drizzle> | undefined
    private user: string | undefined
    env: EnvBinding

    constructor(private c?: Context) {
        this.user = this.c?.req.header("vv-sso-user-id") || import.meta.env.VITE_TEST_USER_ID
        this.env = this.c?.env
    }

    getUser() {
        if (!this.user || this.user === "") {
            throw new ForbiddenError()
        }
        const n = Number(this.user)
        if (isNaN(n) || n === 0) {
            throw new ForbiddenError()
        }
        return n
    }
    async getDb() {
        if (this.db) return this.db
        const mdb = vvbind({ env: this.env }).mysql("mysql")
        this.db = drizzle(async (sql, params, method) => {
            try {
                const rows: any = await mdb.query(sql, params.map(a => typeof a === "number" ? String(a) : a), method)
                if (rows.code) {
                    console.log({ sql, params, method })
                    console.error(rows.data.error)
                    throw new Error(rows.data.error)
                }
                return rows;
            } catch (e: any) {
                console.error('来自 mysql 代理服务器的错误：', e)
                return { rows: [] };
            }
        });
        return this.db
    }
}
```

也可以写一个辅助函数用于创建`hono`端点，避免重复配置类型。

```typescript
export function createHono() {
    const app = new Hono<{
        Bindings: EnvBinding, Variables: {
            ctx: VContext
        }
    }>();
    return app
}
```

并创建一个中间件用于自动创建上下文。

```typescript

app.use(createMiddleware<{
	Bindings: EnvBinding,
	Variables: {
		ctx: VContext
	}
}>(async (c, next) => {
	c.set("ctx", new VContext(c))
	await next()
}))

```

使用`zValidator`对入参进行验证（必须），并返回你的结果。

```typescript
import { zValidator } from "hono/zvalidator";
import { z } from "zod";

import { anotherHono } from "./another"

const app = createHono()
    .get("/user/info", async (c) => {
        const ctx = c.var.ctx // 这样就可以在代码中类型安全的访问上下文了。
        const userId = ctx.getUser()
        const user = (await (await ctx.getDb()).select().from(UserTable).where(eq(UserTable.id, userId)))[0]
        return c.json({ user })
    })
    .post("/createApp", zValidator("json", z.object({
        name: z.string(),
        workerUid: z.string(),
        description: z.string(),
    })), async (c) => {
        const ctx = c.var.ctx
        const body = c.req.valid("json")
        const { name, workerUid, description } = body
        const createBy = ctx.getUser()
        ...
        return c.json({ ... });
    })
    .post(...)
    .route("/appmgr", anotherHono) // 可以创建多个子路由
```

注意，为了类型推导，所有端点都需要链式调用，不能创建新的变量。

## 四、导出后端类型

直接导出app的**类型**即可。

```typescript
export type WebAPI = typeof app
```

## 五、前端RPC调用

利用`hc`创建客户端。

```typescript
import type WebAPI from 'server/api/webapi';
import { hc } from 'hono/client';

export const client = hc<WebAPI>("./api")
```


可以以对象的形式类型安全的直接调用后端API：

```typescript
export async function listApps() {
    return (await client.appmgr.listApps.$post({
        query: {
            offset: "0",
            limit: "10",
        }
    })).json()
}

export async function createApp(info: CreateAppRequest) {
    return (await client.appmgr.createApp.$post({ json: info })).json();
}
```

> [!warning]
> 注意，tsconfig需要开启strict模式，否则类型推导会失败。

