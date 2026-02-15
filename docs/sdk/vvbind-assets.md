# Assets 绑定

Assets 绑定用于访问静态资源文件，支持从资源目录获取文件内容。

## 获取绑定

```typescript
import { vvbind } from "@dicarne/vvorker-sdk";

const assets = vvbind(c).assets("assetsBindName");
```

## 方法

### fetch

获取静态资源文件。

```typescript
fetch(path: string, init?: RequestInit): Promise<Response>
```

**参数**

| 参数 | 类型 | 描述 |
|------|------|------|
| path | string | 资源路径（会自动添加 `http://vvorker.local` 前缀） |
| init | RequestInit | 请求配置（可选） |

**返回值**

- `Promise<Response>` - HTTP 响应对象

**示例**

```typescript
// 获取静态文件
const res = await assets.fetch("/index.html");
const html = await res.text();

// 获取图片
const imgRes = await assets.fetch("/images/logo.png");
const imageBuffer = await imgRes.arrayBuffer();
```

## 完整示例

```typescript
import { vvbind } from "@dicarne/vvorker-sdk";

// 返回静态 HTML 页面
honoApi.get("/page", async (c) => {
    const assets = vvbind(c).assets("assets");
    
    const res = await assets.fetch("/pages/index.html");
    
    return new Response(await res.text(), {
        headers: {
            "Content-Type": "text/html; charset=utf-8"
        }
    });
});

// 返回静态资源（代理模式）
honoApi.get("/static/*", async (c) => {
    const assets = vvbind(c).assets("assets");
    const path = c.req.path.replace("/static", "");
    
    const res = await assets.fetch(path);
    
    // 直接返回响应
    return res;
});

// 根据文件类型设置 Content-Type
honoApi.get("/file/:filename", async (c) => {
    const assets = vvbind(c).assets("assets");
    const filename = c.req.param("filename");
    
    const res = await assets.fetch(`/files/${filename}`);
    
    if (!res.ok) {
        return c.json({ error: "文件不存在" }, 404);
    }
    
    // 获取原始响应
    return res;
});
```

## 配置说明

在 `vvorker.json` 中配置 Assets 绑定：

```json
{
    "assets": {
        "assetsBindName": {
            "type": "assets",
            "directory": "./public"
        }
    }
}
```

## 与前端静态资源的配合

Assets 绑定通常用于：

1. **托管前端应用**：返回 Vue/React 等前端构建产物
2. **静态文件服务**：提供图片、文档等静态资源下载
3. **模板文件**：提供 HTML 模板用于服务端渲染

```typescript
// 托管 SPA 应用
honoApi.get("/*", async (c) => {
    const assets = vvbind(c).assets("assets");
    
    // 尝试获取请求的文件
    let res = await assets.fetch(c.req.path);
    
    // 如果文件不存在，返回 index.html（SPA 路由）
    if (!res.ok) {
        res = await assets.fetch("/index.html");
    }
    
    return res;
});
```
