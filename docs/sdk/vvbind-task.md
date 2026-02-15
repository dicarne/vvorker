# Task 绑定

Task 绑定提供分布式任务管理功能，支持任务创建、状态检查、日志记录和完成通知。

## 获取绑定

```typescript
import { vvbind } from "@dicarne/vvorker-sdk";

const task = vvbind(c).task("taskBindName");
```

## 方法

### create

创建一个新的任务实例。

```typescript
create(name: string, trace_id?: string): Promise<string | undefined>
```

**参数**

| 参数 | 类型 | 描述 |
|------|------|------|
| name | string | 任务名称 |
| trace_id | string | 可选的追踪 ID，用于关联任务 |

**返回值**

- `Promise<string | undefined>` - 返回任务的 trace_id，用于后续操作

**示例**

```typescript
// 创建任务
const traceId = await task.create("process-order");
console.log("任务已创建:", traceId);

// 创建带自定义 trace_id 的任务
const customTraceId = await task.create("sync-data", "custom-trace-123");
```

---

### should_exit

检查任务是否应该退出。用于长时间运行的任务定期检查是否需要中断。

```typescript
should_exit(trace_id: string): Promise<boolean>
```

**参数**

| 参数 | 类型 | 描述 |
|------|------|------|
| trace_id | string | 任务的 trace_id |

**返回值**

- `Promise<boolean>` - 如果应该退出返回 `true`，否则返回 `false`

**示例**

```typescript
const traceId = await task.create("long-running-task");

while (true) {
    // 检查是否需要退出
    if (await task.should_exit(traceId)) {
        console.log("任务被要求退出");
        break;
    }
    
    // 执行任务逻辑
    await doSomeWork();
}
```

---

### log

记录任务日志。

```typescript
log(trace_id: string, text: string): Promise<void>
```

**参数**

| 参数 | 类型 | 描述 |
|------|------|------|
| trace_id | string | 任务的 trace_id |
| text | string | 日志内容 |

**返回值**

- `Promise<void>`

**示例**

```typescript
const traceId = await task.create("data-import");

await task.log(traceId, "开始导入数据...");

for (let i = 0; i < items.length; i++) {
    await processItem(items[i]);
    await task.log(traceId, `已处理 ${i + 1}/${items.length} 项`);
}

await task.log(traceId, "导入完成");
```

---

### complete

标记任务完成。

```typescript
complete(trace_id: string): Promise<void>
```

**参数**

| 参数 | 类型 | 描述 |
|------|------|------|
| trace_id | string | 任务的 trace_id |

**返回值**

- `Promise<void>`

**示例**

```typescript
const traceId = await task.create("cleanup-job");

try {
    await doCleanup();
    await task.log(traceId, "清理完成");
} catch (e) {
    await task.log(traceId, `清理失败: ${e.message}`);
} finally {
    await task.complete(traceId);
}
```

## 完整示例

```typescript
import { vvbind } from "@dicarne/vvorker-sdk";

// 批量处理任务
honoApi.post("/process-batch", async (c) => {
    const task = vvbind(c).task("task");
    const body = await c.req.json();
    const items = body.items;
    
    // 创建任务
    const traceId = await task.create("batch-process");
    await task.log(traceId, `开始处理 ${items.length} 个项目`);
    
    // 异步处理
    c.executionCtx.waitUntil(
        (async () => {
            let processed = 0;
            
            for (const item of items) {
                // 检查是否需要退出
                if (await task.should_exit(traceId)) {
                    await task.log(traceId, "任务被中断");
                    break;
                }
                
                try {
                    await processItem(item);
                    processed++;
                    
                    // 每10个项目记录一次进度
                    if (processed % 10 === 0) {
                        await task.log(traceId, `已处理 ${processed}/${items.length}`);
                    }
                } catch (e) {
                    await task.log(traceId, `处理项目失败: ${e.message}`);
                }
            }
            
            await task.log(traceId, `处理完成，共 ${processed} 个项目`);
            await task.complete(traceId);
        })()
    );
    
    return c.json({ 
        success: true, 
        traceId,
        message: "任务已开始处理" 
    });
});

// 定时任务示例
honoApi.get("/scheduled-cleanup", async (c) => {
    const task = vvbind(c).task("task");
    const traceId = await task.create("scheduled-cleanup");
    
    await task.log(traceId, "开始清理过期数据...");
    
    const expiredRecords = await getExpiredRecords();
    
    for (const record of expiredRecords) {
        if (await task.should_exit(traceId)) {
            await task.log(traceId, "清理任务被中断");
            break;
        }
        
        await deleteRecord(record.id);
    }
    
    await task.log(traceId, "清理任务完成");
    await task.complete(traceId);
    
    return c.json({ success: true });
});
```

## 最佳实践

1. **定期检查退出状态**：长时间运行的任务应该定期调用 `should_exit` 检查是否需要中断。

2. **记录关键日志**：在任务的关键节点调用 `log` 记录进度，便于调试和监控。

3. **确保调用 complete**：无论任务成功还是失败，都应该调用 `complete` 标记任务结束。

4. **使用 try-finally**：确保在异常情况下也能正确完成任务。
