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

## 代码内测试

也可以利用vitest对代码内的函数进行测试，实现如下效果：

```typescript
export function formatDate(date: Date): string {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');

  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}

// 源码内的测试套件
if (import.meta.vitest) {
  const { it, expect } = import.meta.vitest
  it('formatDate', () => {
    expect(formatDate(new Date("2025-08-22 10:14:41"))).toBe("2025-08-22 10:14:41")
  })
}
```

为此，我们需要配置tsconfig.json，添加vitest的类型

```json
// tsconfig.*.json
{
  "extends": "@vue/tsconfig/tsconfig.dom.json",
  "include": [
    "env.d.ts",
    "src/**/*",
    "src/**/*.vue",
    "types/*"
  ],
  "exclude": [
    "src/**/__tests__/*"
  ],
  "compilerOptions": {
    "tsBuildInfoFile": "./node_modules/.tmp/tsconfig.app.tsbuildinfo",
    "allowJs": true,
    "paths": {
      "@/*": [
        "./src/*"
      ],
      "server/*": [
        "./server/*"
      ]
    },
    "types": [
      "./worker-configuration.d.ts",
      "vite/client",
      "node",
      "vitest/importMeta" // 获取import.meta.vitest类型
    ],
    "strict": true
  }
}
```

并且通过配置文件启用代码内的测试函数

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
        includeSource: ['server/**/*.{js,ts}'], // 指定需要测试的源码文件
    },
});
```

最后，配置主vite配置文件，以避免打包中包含测试代码。

```typescript
// vite.config.ts
import { defineConfig } from 'vite'


// https://vite.dev/config/
export default defineConfig({
	// ...
	define: {
		'import.meta.vitest': 'undefined',
	},
})

```