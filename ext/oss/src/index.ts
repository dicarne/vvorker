import { WorkerEntrypoint, env } from "cloudflare:workers";
export * from "./binding"

interface InitResult {
    UploadId: string;
}

interface UploadPartResult {
    ETag: string;
}

interface CompletePart {
    PartNumber: number;
    ETag: string;
}

interface CompleteUploadResult {
    message: string;
    Location: string;
    Bucket: string;
    Key: string;
    ETag: string;
}
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

	async uploadStreamFile(stream: ReadableStream<Uint8Array>, fileName: string, chunkSize = 32 * 1024 * 1024): Promise<CompleteUploadResult> { // 32MB chunks
        const reader = stream.getReader();
        let partNumber = 1;
        let uploadId: string = "";
        const parts: CompletePart[] = [];

        try {
            // 1. Initiate multipart upload
            const initResponse = await fetch(`${GO_API_URL}/api/ext/oss/initiate-multipart-upload`, {
                method: "POST",
                headers: {
                    ...commonConfig,
                    Object: fileName
                }
            });
            const initResult: InitResult = await initResponse.json();
            uploadId = initResult.UploadId;

            // 2. Upload chunks by aggregating smaller reads
            let chunkBuffer: Uint8Array[] = [];
            let bufferSize = 0;

            while (true) {
                const { done, value } = await reader.read();

                if (value) {
                    chunkBuffer.push(value);
                    bufferSize += value.length;
                }

                if (bufferSize >= chunkSize || (done && bufferSize > 0)) {
                    const combinedChunk = new Uint8Array(bufferSize);
                    let offset = 0;
                    for (const chunk of chunkBuffer) {
                        combinedChunk.set(chunk, offset);
                        offset += chunk.length;
                    }

                    const formData = new FormData();
                    formData.append("file", new Blob([combinedChunk]), `part-${partNumber}`);

                    const uploadResponse = await fetch(`${GO_API_URL}/api/ext/oss/upload-part`, {
                        method: "POST",
                        headers: {
                            ...commonConfig,
                            Object: fileName,
                            "x-amz-upload-id": uploadId,
                            "x-amz-part-number": partNumber.toString()
                        },
                        body: formData
                    });

                    const uploadResult: UploadPartResult = await uploadResponse.json();
                    parts.push({
                        PartNumber: partNumber,
                        ETag: uploadResult.ETag
                    });

                    partNumber++;
                    chunkBuffer = [];
                    bufferSize = 0;
                }

                if (done) {
                    break;
                }
            }

            // 3. Complete multipart upload
            const completeResponse = await fetch(`${GO_API_URL}/api/ext/oss/complete-multipart-upload`, {
                method: "POST",
                headers: {
                    ...commonConfig,
                    Object: fileName,
                    "x-amz-upload-id": uploadId
                },
                body: JSON.stringify({
                    Parts: parts
                })
            });

            return await completeResponse.json();
        } catch (error) {
            // If there's an error, try to abort the multipart upload
            if (uploadId) {
                try {
                    await fetch(`${GO_API_URL}/api/ext/oss/abort-multipart-upload`, {
                        method: "POST",
                        headers: {
                            ...commonConfig,
                            Object: fileName,
                            "x-amz-upload-id": uploadId
                        }
                    });
                } catch (abortError) {
                    console.error("Failed to abort multipart upload:", abortError);
                }
            }
            throw error;
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
