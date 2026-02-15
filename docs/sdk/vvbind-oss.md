# OSS 绑定

OSS 绑定提供对象存储服务功能，支持文件的上传、下载、列表和删除操作。

## 获取绑定

```typescript
import { vvbind } from "@dicarne/vvorker-sdk";

const oss = vvbind(c).oss("ossBindName");
```

## 方法

### listBuckets

列出所有存储桶。

```typescript
listBuckets(): Promise<any>
```

**返回值**

- `Promise<any>` - 存储桶列表

**示例**

```typescript
const buckets = await oss.listBuckets();
console.log(buckets);
```

---

### listObjects

列出指定存储桶中的所有对象。

```typescript
listObjects(bucket: string): Promise<any>
```

**参数**

| 参数 | 类型 | 描述 |
|------|------|------|
| bucket | string | 存储桶名称 |

**返回值**

- `Promise<any>` - 对象列表

**示例**

```typescript
const objects = await oss.listObjects("my-bucket");
console.log(objects);
```

---

### downloadFile

下载文件，返回文件的二进制数据。

```typescript
downloadFile(fileName: string): Promise<Uint8Array>
```

**参数**

| 参数 | 类型 | 描述 |
|------|------|------|
| fileName | string | 文件名（包含路径） |

**返回值**

- `Promise<Uint8Array>` - 文件的二进制数据

**示例**

```typescript
const fileData = await oss.downloadFile("images/avatar.png");
// 转换为 Blob 或进行其他处理
const blob = new Blob([fileData], { type: "image/png" });
```

---

### uploadFile

上传文件。

```typescript
uploadFile(data: Uint8Array, fileName: string): Promise<any>
```

**参数**

| 参数 | 类型 | 描述 |
|------|------|------|
| data | Uint8Array | 文件的二进制数据 |
| fileName | string | 目标文件名（包含路径） |

**返回值**

- `Promise<any>` - 上传结果

**示例**

```typescript
// 从字符串创建文件
const text = "Hello, World!";
const data = new TextEncoder().encode(text);
await oss.uploadFile(data, "documents/hello.txt");

// 上传图片
const imageBuffer = await fetch("https://example.com/image.png")
    .then(r => r.arrayBuffer());
await oss.uploadFile(new Uint8Array(imageBuffer), "images/downloaded.png");
```

---

### uploadStreamFile

以流的方式上传文件。

```typescript
uploadStreamFile(stream: ReadableStream<Uint8Array>, fileName: string): Promise<any>
```

**参数**

| 参数 | 类型 | 描述 |
|------|------|------|
| stream | `ReadableStream<Uint8Array>` | 可读流 |
| fileName | string | 目标文件名（包含路径） |

**返回值**

- `Promise<any>` - 上传结果

**示例**

```typescript
// 从请求中获取文件流并上传
const fileStream = c.req.body;
await oss.uploadStreamFile(fileStream, "uploads/user-file.bin");
```

---

### downloadStreamFile

以流的方式下载文件。

```typescript
downloadStreamFile(fileName: string): Promise<ReadableStream<Uint8Array>>
```

**参数**

| 参数 | 类型 | 描述 |
|------|------|------|
| fileName | string | 文件名（包含路径） |

**返回值**

- `Promise<ReadableStream<Uint8Array>>` - 文件的可读流

**示例**

```typescript
const stream = await oss.downloadStreamFile("videos/demo.mp4");
// 流式处理大文件
const reader = stream.getReader();
while (true) {
    const { done, value } = await reader.read();
    if (done) break;
    // 处理数据块
    console.log("Received chunk:", value.length);
}
```

---

### deleteObject

删除指定文件。

```typescript
deleteObject(fileName: string): Promise<any>
```

**参数**

| 参数 | 类型 | 描述 |
|------|------|------|
| fileName | string | 要删除的文件名（包含路径） |

**返回值**

- `Promise<any>` - 删除结果

**示例**

```typescript
await oss.deleteObject("temp/old-file.txt");
```

## 完整示例

```typescript
import { vvbind } from "@dicarne/vvorker-sdk";

// 文件上传接口
honoApi.post("/upload", async (c) => {
    const oss = vvbind(c).oss("myoss");
    const formData = await c.req.formData();
    const file = formData.get("file") as File;

    // 读取文件内容
    const buffer = await file.arrayBuffer();
    const fileName = `uploads/${Date.now()}-${file.name}`;

    // 上传到 OSS
    await oss.uploadFile(new Uint8Array(buffer), fileName);

    return c.json({
        success: true,
        path: fileName
    });
});

// 文件下载接口
honoApi.get("/download/:filename", async (c) => {
    const oss = vvbind(c).oss("myoss");
    const filename = c.req.param("filename");

    const data = await oss.downloadFile(`uploads/${filename}`);

    return new Response(data, {
        headers: {
            "Content-Type": "application/octet-stream",
            "Content-Disposition": `attachment; filename="${filename}"`
        }
    });
});
```
