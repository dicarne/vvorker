export default {
	fetch: (request: Request, env: any) => {
		return new Response(JSON.stringify(!!env.worker))
	}
}