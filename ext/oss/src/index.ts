import { WorkerEntrypoint, env } from "cloudflare:workers";


export default class OSS extends WorkerEntrypoint {
	constructor(ctx: any, env: any) {
		super(ctx, env)

	}
	async call() {

		return "hi"
	}
}