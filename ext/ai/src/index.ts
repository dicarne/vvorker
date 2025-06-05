
import OpenAI from "openai";
import { WorkerEntrypoint } from 'cloudflare:workers'

export default class OpenAIWorker extends WorkerEntrypoint {
	API_KEY: string
	BASE_URL: string
	MODEL: string
	constructor(ctx: any, env: any) {
		super(ctx, env)
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
	async chat(messages: any[]) {
		const openai = new OpenAI({
			apiKey: this.API_KEY,
			baseURL: this.BASE_URL,
		})
		let r = await openai.chat.completions.create({
			messages: messages,
			model: this.MODEL,
		})
		return r
	}
	async chatStream(messages: any[]) {
		const openai = new OpenAI({
			apiKey: this.API_KEY,
			baseURL: this.BASE_URL,
		})
		let r = await openai.chat.completions.create({
			messages: messages,
			model: this.MODEL,
			stream: true,
		})

		let chatStreamIterator: AsyncIterator<any> = r[Symbol.asyncIterator]();
		let err: any = null;
		return {

			async next() {
				try {
					if (chatStreamIterator) {
						return chatStreamIterator.next();
					}
				} catch (error) {
					err = error;
				} finally {
					return { done: true, err };
				}
			}
		}
	}
}