
import OpenAI from "openai";
import { WorkerEntrypoint } from 'cloudflare:workers'

export default class OpenAIWorker extends WorkerEntrypoint {
	count: number
	API_KEY: string
	BASE_URL: string
	MODEL: string
	constructor(ctx: any, env: any) {
		super(ctx, env)
		this.count = 0;
		this.API_KEY = env.API_KEY;
		this.BASE_URL = env.BASE_URL;
		this.MODEL = env.MODEL;
	}
	async ask(messages: any[]) {
		const openai = new OpenAI({
			apiKey: this.API_KEY,
			baseURL: this.BASE_URL,
		})
		let r = await openai.chat.completions.create({
			messages: messages,
			model: this.MODEL,
		})
		return r.choices[0].message.content
	}
}