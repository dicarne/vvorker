
import OpenAI from "openai";

const res = {
	async createOpenAIClient(API_KEY: string, BASE_URL: string, MODEL: string) {
		return async (messages: any[]) => {
			const openai = new OpenAI({
				apiKey: API_KEY,
				baseURL: BASE_URL,
			})
			let r = await openai.chat.completions.create({
				messages: messages,
				model: MODEL,
			})
			return r.choices[0].message.content
		}
	},
}

export default function () {
	return res
}
