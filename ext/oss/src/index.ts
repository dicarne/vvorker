import { WorkerEntrypoint, env } from "cloudflare:workers";
let env1 = env as unknown as any
// 假设Go的接口地址
let GO_API_URL = env1.OSS_AGENT_URL;
const {
	HOST,
	PORT,
	ACCESS_KEY_ID,
	ACCESS_KEY_SECRET,
	BUCKET,
	USE_SSL,
	REGION,
	RESOURCE_ID,
	X_SECRET,
	X_NODENAME
} = env1;

let commonConfig = {
	Endpoint: `${HOST}:${PORT}`,
	AccessKeyID: ACCESS_KEY_ID,
	SecretAccessKey: ACCESS_KEY_SECRET,
	UseSSL: USE_SSL,
	Region: REGION,
	Bucket: BUCKET,
	ResourceID: RESOURCE_ID,
	"x-secret": X_SECRET,
	"x-node-name": X_NODENAME,
}

export default class OSS extends WorkerEntrypoint {
	constructor(ctx: any, env: any) {
		super(ctx, env)
	}

	async listBuckets() {
		const response = await fetch(`${GO_API_URL}/api/ext/oss/list-buckets`, {
			method: "POST",
			headers: {
				...commonConfig
			},
		});
		return response.json();
	}


	async uploadFile(fileData: Uint8Array, fileName: string) {
		const formData = new FormData();
		// 将字节流转换为 Blob 再添加到 FormData
		const blob = new Blob([fileData]);
		formData.append("file", blob, fileName);
		const response = await fetch(`${GO_API_URL}/api/ext/oss/upload`, {
			method: "POST",
			headers: {
				...commonConfig,
				Object: fileName
			},
			body: formData,
		});
		return response.json();
	}

	async uploadStreamFile(stream: ReadableStream<Uint8Array>, fileName: string) {
		const formData = new FormData();
		const reader = stream.getReader();
		const chunks: Uint8Array[] = [];

		try {
			while (true) {
				const { done, value } = await reader.read();
				if (done) break;
				chunks.push(value);
			}

			// Combine all chunks into a single Uint8Array
			const totalLength = chunks.reduce((acc, chunk) => acc + chunk.length, 0);
			const combined = new Uint8Array(totalLength);
			let offset = 0;
			for (const chunk of chunks) {
				combined.set(chunk, offset);
				offset += chunk.length;
			}

			// Create blob from combined data
			const blob = new Blob([combined]);
			formData.append("file", blob, fileName);

			const response = await fetch(`${GO_API_URL}/api/ext/oss/upload`, {
				method: "POST",
				headers: {
					...commonConfig,
					Object: fileName
				},
				body: formData,
			});
			return response.json();
		} finally {
			reader.releaseLock();
		}
	}

	async downloadFile(fileName: string) {
		const response = await fetch(`${GO_API_URL}/api/ext/oss/download`, {
			method: "POST",
			headers: {
				...commonConfig,
				Object: fileName
			},
		});
		return response.bytes();
	}

	async downloadStreamFile(fileName: string) {
		const response = await fetch(`${GO_API_URL}/api/ext/oss/download`, {
			method: "POST",
			headers: {
				...commonConfig,
				Object: fileName
			},
		});
		return response.body;
	}

	async listObjects(path: string, recursive: boolean = false) {
		const response = await fetch(`${GO_API_URL}/api/ext/oss/list-objects`, {
			method: "POST",
			headers: {
				...commonConfig,
				Path: path,
				Recursive: recursive ? "true" : "false",
			},
		});
		return response.json();
	}

	async deleteObject(fileName: string) {
		const response = await fetch(`${GO_API_URL}/api/ext/oss/delete`, {
			method: "POST",
			headers: {
				...commonConfig,
				Object: fileName
			},
		});
		return response.json();
	}
}
