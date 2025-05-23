
import OpenAI from "openai";

// const res = {
// 	async createOpenAIClient(API_KEY: string, BASE_URL: string, MODEL: string) {
// 		return async (messages: any[]) => {
// 			const openai = new OpenAI({
// 				apiKey: API_KEY,
// 				baseURL: BASE_URL,
// 			})
// 			let r = await openai.chat.completions.create({
// 				messages: messages,
// 				model: MODEL,
// 			})
// 			return r.choices[0].message.content
// 		}
// 	},
// }

// export default function () {
// 	return res
// }

import { WorkerEntrypoint } from 'cloudflare:workers'

// 转换为函数形式
function createChainExample() {
  let count = 0;

  // 返回一个对象，包含 add 和 result 方法
  return {
    add(num: number) {
      count += num;
      return this;
    },
    result() {
      return count;
    }
  };
}

// 原有的类定义可以移除
// class ChainExcample {
//   count: number
//   constructor() {
//     this.count = 0;
//   }
//   add(num: number) {
//     this.count += num;
//     return this
//   }
//   result() {
//     return this.count
//   }
// }

export default class OpenAIWorker extends WorkerEntrypoint {
	count: number
	constructor(ctx: any, env: any) {
		super(ctx, env)
		this.count = 0;
	}
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
	}

	anotherFunction() {
		return createChainExample()
	}

}