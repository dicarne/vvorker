// filepath: src/index.ts
export * from "./binding"
import { WorkerEntrypoint, env } from 'cloudflare:workers'

const eenv = env as unknown as any

eenv.RESOURCE_ID = eenv.RESOURCE_ID || ""
const prefix = eenv.RESOURCE_ID.length > 0 ? eenv.RESOURCE_ID + ":" : ""

const masterEndpoint = eenv.MASTER_ENDPOINT
const commonConfig = {
	"x-secret": eenv.X_SECRET,
	"x-node-name": eenv.X_NODENAME,
	"resource-id": eenv.RESOURCE_ID,
}

async function invoke<T>(config: any) {
	const r = await (await fetch(`${masterEndpoint}/api/ext/kv/invoke`, {
		method: "POST",
		headers: {
			...commonConfig
		},
		body: JSON.stringify({
			rid: eenv.RESOURCE_ID,
			...config
		})
	})).json() as {
		code: number,
		data: T
	}
	if (r.code !== 0) {
		return undefined
	}
	return r.data
}

export default class KV extends WorkerEntrypoint {
	constructor(ctx: any, env: any) {
		super(ctx, env)
	}


	async get(key: string): Promise<string | null> {
		return await invoke<string>({
			method: "get",
			key: key
		}) ?? null
	}

	async set(key: string, value: string, options?: {
		EX?: number,
		NX?: boolean,
		XX?: boolean
	} | number): Promise<number> {
		return await invoke<number>({
			method: "set",
			key: key,
			value: value,
			options: typeof options === "number" ? {
				EX: options
			} : options
		}) ?? -1
	}

	async del(key: string): Promise<void> {
		await invoke({
			method: "del",
			key: key
		})
	}

	async keys(pattern: string, offset: number, size: number): Promise<string[]> {
		return await invoke<string[]>({
			method: "keys",
			pattern: pattern,
			offset: offset,
			size: size
		}) ?? []
	}
}