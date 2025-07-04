// filepath: src/index.ts
import { createClient } from "redis";
export * from "./binding"
import { RpcTarget, WorkerEntrypoint, env } from 'cloudflare:workers'
import { KV } from "./binding";

const eenv = env as unknown as any

eenv.RESOURCE_ID = eenv.RESOURCE_ID || ""
const prefix = eenv.RESOURCE_ID.length > 0 ? eenv.RESOURCE_ID + ":" : ""

const masterEndpoint = eenv.MASTER_ENDPOINT
const provider = eenv.KVPROVIDER
const commonConfig = {
	"x-secret": eenv.X_SECRET,
	"x-node-name": eenv.X_NODENAME,
	"resource-id": eenv.RESOURCE_ID,
}

async function client() {
	let redis = createClient({
		url: `redis://${eenv.HOST}:${eenv.PORT}/0`
	})
	await redis!.connect()
	return redis
}

class NutsDBTarget extends RpcTarget implements KV {
	async invoke<T>(config: any) {
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
			console.error(r)
			return undefined
		}
		return r.data
	}

	async get(key: string): Promise<string> {
		return await this.invoke<string>({
			method: "get",
			key: key
		}) ?? ""
	}

	async set(key: string, value: string, ttl?: number): Promise<void> {
		await this.invoke({
			method: "set",
			key: key,
			value: value,
			ttl: ttl
		})
	}

	async del(key: string): Promise<void> {
		await this.invoke({
			method: "del",
			key: key
		})
	}

	async keys(pattern: string, offset: number, size: number): Promise<string[]> {
		return await this.invoke<string[]>({
			method: "keys",
			pattern: pattern,
			offset: offset,
			size: size
		}) ?? []
	}
}

class RedisTarget extends RpcTarget implements KV {
	constructor() {
		super()
	}
	redis: ReturnType<typeof createClient> | null = null

	async connect() {
		this.redis = await client()
	}

	async get(key: string): Promise<string> {
		return await this.redis!.get(prefix + key)
	}

	async set(key: string, value: string, ttl?: number): Promise<void> {
		if(!ttl) return await this.redis!.set(prefix + key, value)
		return await this.redis!.setEx(prefix + key, ttl, value)
	}

	async del(key: string): Promise<void> {
		return await this.redis!.del(prefix + key)
	}

	async keys(pattern: string, offset: number, size: number): Promise<string[]> {
		return await this.redis!.keys(prefix + pattern)
	}
}


export default class Redis extends WorkerEntrypoint {
	constructor(ctx: any, env: any) {
		super(ctx, env)
	}
	async client(): Promise<KV> {
		if (provider === "redis") {
			let r = new RedisTarget()
			await r.connect()
			return r
		}
		let r = new NutsDBTarget()
		return r
	}
}