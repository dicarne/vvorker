# 请求加密

VVorker 支持对 Worker 的请求体进行 AES-GCM 加密，确保客户端与 Worker 之间的数据在传输过程中不被窃取。

## 工作原理

1. 客户端使用 AES-GCM 算法加密请求体，并设置 `X-Encrypted-Data: true` 请求头
2. 代理层检测到加密头后，查找该 Worker 的加密规则，获取密钥
3. 代理层解密请求体，将明文转发给 Worker
4. Worker 收到的是正常的明文请求，无需做任何解密处理

```
客户端 ──[加密请求]──> 代理层 ──[解密后明文]──> Worker
```

## 配置步骤

### 1. 添加加密规则

1. 进入 Worker 编辑页面，切换到 **规则** 标签
2. 点击 **添加规则**
3. 选择规则类型为 **加密**
4. 路由前缀填写 `/`（对所有请求生效，或填写特定路径前缀）
5. 点击 **自动生成** 按钮生成一个 32 字节的 AES-256 密钥，也可手动输入 16/24/32 字节的密钥
6. 填写描述，点击确认保存

> [!NOTE]
> 加密规则独立于访问控制，不需要开启"启用访问控制"开关即可使用。

### 2. 客户端使用

#### 使用 SDK（推荐）

SDK 提供了 `encryptedFetch` 函数，是 `fetch` 的直接替代，自动完成加密：

```bash
# 设置环境变量
export VITE_VVORKER_ENCRYPTION_KEY="你的32字节密钥"
```

```typescript
import { encryptedFetch } from "@dicarne/vvorker-sdk";

// 替换原来的 fetch 调用
const res = await encryptedFetch("https://your-worker.example.com/api/data", {
  method: "POST",
  body: JSON.stringify({ message: "hello" }),
});
```

#### 使用 curl

手动加密后发送请求：

```bash
# 需要先用工具将 body 加密为 Base64 格式
curl -X POST https://your-worker.example.com/api/data \
  -H "X-Encrypted-Data: true" \
  -H "Content-Type: application/json" \
  -d '"加密后的Base64数据"'
```

## SDK 函数

### `encryptedFetch(url, init?)`

`fetch` 的加密替代版本。从 `VITE_VVORKER_ENCRYPTION_KEY` 环境变量读取密钥，自动加密请求体。

```typescript
import { encryptedFetch } from "@dicarne/vvorker-sdk";

// 用法与 fetch 完全一致
const res = await encryptedFetch("/api/data", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({ key: "value" }),
});
```

支持的 body 类型：`string`、`ArrayBuffer`、`ArrayBufferView`。

### `encryptData(data, key)`

底层加密函数，使用 AES-GCM 加密数据并返回 Base64 字符串。

```typescript
import { encryptData } from "@dicarne/vvorker-sdk";

const encrypted = await encryptData("原始数据", "你的32字节密钥");
// 返回 Base64 编码的加密数据
```

## 安全建议

- 使用 **32 字节**密钥（AES-256）以获得最高安全性
- 通过环境变量传递密钥，不要硬编码在代码中
- 密钥应当保密，不要提交到版本控制系统
- 建议定期轮换密钥
