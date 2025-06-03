// filepath: src/index.ts

import { env } from 'cloudflare:workers'

const eenv = env as unknown as any

let commonConfig = {
	"x-secret": eenv.X_SECRET,
	"x-node-name": eenv.X_NODENAME,
}

export default {
	async fetch(request: any, env: any) {
		const url = new URL(request.url);
		return fetch(`${eenv.MASTER_ENDPOINT}/api/ext/assets/get-assets`, {
			method: "GET",
			headers: {
				...commonConfig,
				"vvorker-asset-path": url.pathname,
				"vvorker-asset-worker-uid": eenv.WORKER_UID,
			}
		})
	},
};