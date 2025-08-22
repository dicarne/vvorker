# 测试

## 安装依赖

```bash
pnpm add -D vitest@~3.2.0 @cloudflare/vitest-pool-workers
```

## 创建配置文件

```typescript
// vitest.config.ts
import { defineWorkersConfig } from "@cloudflare/vitest-pool-workers/config";

export default defineWorkersConfig({
  test: {
    poolOptions: {
      workers: {
        wrangler: { configPath: "./wrangler.jsonc" },
      },
    },
  },
});
```

```json
// test/tsconfig.json 
{
  "extends": "../tsconfig.json",
  "compilerOptions": {
    "moduleResolution": "bundler",
    "types": [
      "@cloudflare/vitest-pool-workers", // provides `cloudflare:test` types
    ],
  },
  "include": [
    "./**/*.ts",
    "../worker-configuration.d.ts", // output of `wrangler types`
  ],
}
```

## 单元测试
```typescript
import {
  env,
  createExecutionContext,
  waitOnExecutionContext,
} from "cloudflare:test";
import { describe, it, expect } from "vitest";
// Import your worker so you can unit test it
import worker from "../server"; // ../src

describe("Hello World worker", () => {
  it("responds with Hello World!", async () => {
    const request = new Request("http://example.com/404");
    // Create an empty context to pass to `worker.fetch()`
    const ctx = createExecutionContext();
    const response = await worker.fetch(request, env, ctx);
    // Wait for all `Promise`s passed to `ctx.waitUntil()` to settle before running test assertions
    await waitOnExecutionContext(ctx);
    expect(response.status).toBe(404);
    expect(await response.text()).toBe("Not found");
  });
});
```