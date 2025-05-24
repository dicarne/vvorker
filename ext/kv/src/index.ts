// filepath: src/index.ts
import { createClient } from "redis";

import { WorkerEntrypoint, env } from 'cloudflare:workers'

export default class Redis extends WorkerEntrypoint {

	constructor(ctx: any, env: any) {
		super(ctx, env)
	}
	async call() {
		let redis = createClient({
			url: "redis://127.0.0.1:6379/0"
		})
		let c = await redis.connect()
		await c.set("key", "hello redis set value")
		return await c.get("key")
	}
}