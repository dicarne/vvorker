
import { env, WorkerEntrypoint } from 'cloudflare:workers'
import { v4 } from 'uuid'

const eenv = env as unknown as any

const {
	MASTER_ENDPOINT,
	WORKER_UID,
	X_SECRET,
	X_NODENAME
} = eenv
let commonConfig = {
	"x-secret": X_SECRET,
	"x-node-name": X_NODENAME,
	"Content-Type": "application/json",
}
export default class Task extends WorkerEntrypoint {
	async start() {
		const id = v4()
		let c1 = await fetch(`${MASTER_ENDPOINT}/api/ext/task/create`, {
			method: "POST",
			headers: {
				...commonConfig
			},
			body: JSON.stringify({ trace_id: id, worker_uid: WORKER_UID })
		})
		let c2 = await c1.json() as any
		if (c2.code != 0) {
			return undefined
		}
		return {
			async should_exit() {
				let c2 = await (await fetch(`${MASTER_ENDPOINT}/api/ext/task/check`, {
					method: "POST",
					headers: {
						...commonConfig
					},
					body: JSON.stringify({ trace_id: id, worker_uid: WORKER_UID })
				})).json() as any
				if (c2.code != 0) {
					return undefined
				}
				return c2.data.status === "canceled"
			},
			async complete() {
				let c3 = await (await fetch(`${MASTER_ENDPOINT}/api/ext/task/complete`, {
					method: "POST",
					headers: {
						...commonConfig
					},
					body: JSON.stringify({ trace_id: id, worker_uid: WORKER_UID })
				}))

			},
			async log(text: string) {
				let c4 = await (await fetch(`${MASTER_ENDPOINT}/api/ext/task/log`, {
					method: "POST",
					headers: {
						...commonConfig
					},
					body: JSON.stringify({ trace_id: id, worker_uid: WORKER_UID, log: String(text) })
				}))
			}
		}
	}
}
