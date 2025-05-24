// filepath: src/index.ts
import { createClient } from "redis";

import { WorkerEntrypoint, env } from 'cloudflare:workers'

export default class Redis extends WorkerEntrypoint {

	constructor(ctx: any, env: any) {
		super(ctx, env)
	}
	async call() {
		let redis = createClient({
			url: `redis://${env.ENDPOINT}:${env.PORT}/0`
		})
		let c = await redis.connect()
		await c.set("key", "hello redis set value")
		let r = await c.get("key")
		redis.close()
		return {
			result: r,
			endpoint: env.ENDPOINT,
			port: env.PORT
		}
	}
}