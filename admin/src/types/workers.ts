export interface WorkerItem {
  UID: string
  ExternalPath?: string
  HostName?: string
  NodeName?: string
  Port?: number
  Entry?: string
  Code: string
  Name: string
  Template: string
}

export interface WorkerItemProperties {
  item: WorkerItem
}

// name is generated on server side
// @ts-expect-error
export const DEFAULT_WORKER_ITEM: WorkerItem = {
  UID: 'worker',
  Code: btoa(`export default {
  async fetch(req, env) {
    try {
		let resp = new Response("worker: " + req.url + " is online! -- " + new Date())
		return resp
	} catch(e) {
		return new Response(e.stack, { status: 500 })
	}
  }
};`),
  Template: `
{
    "name": "worker",
    "version": "1.0.0",
    "extensions": [],
    "services": [],
    "vars": {},
    "ai": [],
    "oss": [],
    "pgsql": [],
    "kv": []
}
`
}

export interface WorkerEditorProperties {
  item: string
}

export interface VorkerSettingsProperties {
  WorkerURLSuffix: string
  Scheme: string
  EnableRegister: boolean
  UrlType: string
  ApiUrl: string
}

export interface Task {
  worker_uid: string
  trace_id: string
  status: "completed" | "running" | "canceled" | "failed"
  start_time: string
  end_time: string
  worker_name: string
}


export interface TaskLog {
  time: string
  content: string
  type: string
}

// 参考 Go 语言的 WorkerLogData 结构体定义 TypeScript 接口
export interface WorkerLog {
  // 工作者的唯一标识符
  uid: string;
  // 日志输出内容
  output: string;
  // 日志记录时间，在 TypeScript 里用字符串表示日期时间
  time: string;
  // 日志类型
  type: string;
  // 日志的唯一标识符
  log_uid: string;
}