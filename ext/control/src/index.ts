export default {
	fetch: async (request: Request, env: any, ctx: ExecutionContext) => {
		let body: any = await request.json()
		if (body.type === "scheduled") {
			ctx.waitUntil(env.worker.scheduled(body))
		}
		return new Response(JSON.stringify({
			code: 0,
		}))
	}
}