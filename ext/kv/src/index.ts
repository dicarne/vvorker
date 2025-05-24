// filepath: src/index.ts
import { createClient } from "redis";

import { WorkerEntrypoint, env } from 'cloudflare:workers'

const eenv = env as unknown as any

eenv.RESOURCE_ID = eenv.RESOURCE_ID || ""
const prefix = eenv.RESOURCE_ID.length > 0 ? eenv.RESOURCE_ID + ":" : ""

export default class Redis extends WorkerEntrypoint {
	redis = createClient({
		url: `redis://${eenv.ENDPOINT}:${eenv.PORT}/0`
	})
	constructor(ctx: any, env: any) {
		super(ctx, env)
	}
	async start() {
		await this.redis.connect()
	}
	async end() {
		await this.redis.quit()
	}
	async get(key: string) {
		return await this.redis.get(prefix + key)
	}
	async set(key: string, value: string) {
		return await this.redis.set(prefix + key, value)
	}
	async del(key: string) {
		return await this.redis.del(prefix + key)
	}
	async keys(pattern: string) {
		return await this.redis.keys(prefix + pattern)
	}

}