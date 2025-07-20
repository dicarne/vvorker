import { KV } from "@dicarne/vvorker-kv";

export interface DistributedTaskStatus {
    /**
     * 任务状态
     * pending: 任务未开始
     * running: 任务正在运行
     * completed: 任务完成
     * failed: 任务失败（可重试）
     * error: 任务出错（不可重试）
     */
    status: "pending" | "running" | "completed" | "failed" | "error"
    /**
     * 任务ID
     */
    id: string
    /**
     * 任务创建时间
     */
    createdAt: number
    /**
     * 任务更新时间
     */
    updatedAt: number
    /**
     * 任务结果
     */
    result?: any

    payload?: any
}

export class DistributedTask {
    constructor(private kv: KV, public id: string) { }

    async status(): Promise<DistributedTaskStatus | null> {
        let r = await this.kv.get(`distributed-task:${this.id}`)
        if (!r) return null
        return JSON.parse(r)
    }

    async start(payload?: any) {
        if (!this.status()) {
            return false
        }
        await this.kv.set(`distributed-task:${this.id}`, JSON.stringify({
            id: this.id,
            status: "pending",
            createdAt: Date.now(),
            updatedAt: Date.now(),
            result: undefined,
            payload: payload
        }))
        return true
    }
}
