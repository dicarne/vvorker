import { WorkerEntrypoint, env } from "cloudflare:workers";

export default class extends WorkerEntrypoint {
	async invoke(url: string, init: any) {
		return await (await env.internalNet.fetch(url, init)).text();
	}
}