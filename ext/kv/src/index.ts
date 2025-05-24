// filepath: src/index.ts
import { createClient } from "redis";

import { WorkerEntrypoint, env } from 'cloudflare:workers'

const eenv = env as unknown as any

eenv.RESOURCE_ID = eenv.RESOURCE_ID || ""
const prefix = eenv.RESOURCE_ID.length > 0 ? eenv.RESOURCE_ID + ":" : ""

export default class Redis extends WorkerEntrypoint {
	constructor(ctx: any, env: any) {
		super(ctx, env)
	}
	async start() {
		const redis = createClient({
			url: `redis://${eenv.ENDPOINT}:${eenv.PORT}/0`
		})
		await redis.connect()
		return {
			end: async () => {
				await redis.quit()
			},
			get: async (key: string) => {
				return await redis.get(prefix + key)
			},
			set: async (key: string, value: string) => {
				return await redis.set(prefix + key, value)
			},
			del: async (key: string) => {
				return await redis.del(prefix + key)
			},
			keys: async (pattern: string) => {
				return await redis.keys(prefix + pattern)
			}
		}
	}
}