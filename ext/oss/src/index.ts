import { WorkerEntrypoint, env } from "cloudflare:workers";
let env1 = env as unknown as any
// 假设Go的接口地址
let GO_API_URL = env1.OSS_AGENT_URL;
const {
	ENDPOINT,
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
	Endpoint: `${ENDPOINT}:${PORT}`,
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
