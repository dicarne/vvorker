// filepath: src/index.ts
export * from "./binding"
import { RpcTarget, WorkerEntrypoint, env } from 'cloudflare:workers'

const eenv = env as unknown as any
function config() {
	// 从环境变量中获取配置信息
	let cfg = {
		"user": eenv.USER,
		"host": eenv.HOST,
		"port": eenv.PORT,
		"password": eenv.PASSWORD,
		"database": eenv.DATABASE,
	}

	// 遍历配置对象，检查每个属性是否为空
	for (const [key, value] of Object.entries(cfg)) {
		if (!value) {
			throw new Error(`Environment variable ${key.toUpperCase()} is missing or empty`);
		}
	}

	return cfg
}

const cfg = config()

export default class MySQL extends WorkerEntrypoint {
	constructor(ctx: any, env: any) {
		super(ctx, env)
	}

	connectionString() {
		return `mysql://${cfg.user}:${cfg.password}@${cfg.host}:${cfg.port}/${cfg.database}`;
	}
}