// filepath: src/index.ts
import { createClient } from "redis";
export * from "./binding"
import { RpcTarget, WorkerEntrypoint, env } from 'cloudflare:workers'

const eenv = env as unknown as any

eenv.RESOURCE_ID = eenv.RESOURCE_ID || ""
const prefix = eenv.RESOURCE_ID.length > 0 ? eenv.RESOURCE_ID + ":" : ""

async function client() {
	let redis = createClient({
		url: `redis://${eenv.HOST}:${eenv.PORT}/0`
	})
	await redis!.connect()
	return redis
}

class RedisTarget extends RpcTarget {
	constructor() {
		super()

	}
	redis: ReturnType<typeof createClient> | null = null

	async connect() {
		this.redis = await client()
	}

	async get(key: string) {
		return await this.redis!.get(prefix + key)
	}

	async set(key: string, value: string) {
		return await this.redis!.set(prefix + key, value)
	}

	async del(key: string) {
		return await this.redis!.del(prefix + key)
	}

	async keys(pattern: string) {
		return await this.redis!.keys(prefix + pattern)
	}

	async ping() {
		return await this.redis?.ping()
	}
}


export default class Redis extends WorkerEntrypoint {
	constructor(ctx: any, env: any) {
		super(ctx, env)
	}
	async client() {
		let r = new RedisTarget()
		await r.connect()
		return r
	}
}