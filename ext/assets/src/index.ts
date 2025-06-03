// filepath: src/index.ts

import { WorkerEntrypoint, env } from 'cloudflare:workers'

const eenv = env as unknown as any

eenv.RESOURCE_ID = eenv.RESOURCE_ID || ""

export default {
	async fetch(request: any, env: any) {
		return new Response("Hello, World!")
	},
};