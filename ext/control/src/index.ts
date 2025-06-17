export default {
	fetch: async (request: Request, env: any) => {
		let body: any = await request.json()
		if (body.type === "scheduled") {
			await env.worker.scheduled(body)
		}
		return new Response(JSON.stringify({
			code: 0,
		}))
	}
}