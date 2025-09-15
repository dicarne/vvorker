// filepath: src/index.ts
export * from "./binding"
import { WorkerEntrypoint, env } from 'cloudflare:workers'

const eenv = env as unknown as any

let commonConfig = {
	"x-secret": eenv.X_SECRET,
	"x-node-name": eenv.X_NODENAME,
}

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
		return `mysql://${cfg.user}:${encodeURIComponent(cfg.password)}@${cfg.host}:${cfg.port}/${cfg.database}`;
	}
	connectionInfo() {
		return {
			user: cfg.user,
			host: cfg.host,
			database: cfg.database,
			password: cfg.password,
			port: Number(cfg.port),
		}
	}
	async query(sql: string, params: any, method: string) {
		return (await rpc(sql, params, method,
			`${cfg.user}:${cfg.password}@tcp(${cfg.host}:${cfg.port})/${cfg.database}`
		)).json()
	}
}




async function rpc(sql: string, params: any, method: string, connection_string: string) {
	return fetch(`${eenv.MASTER_ENDPOINT}/api/ext/mysql/query`, {
		method: "POST",
		headers: {
			...commonConfig,
		},
		body: JSON.stringify({
			sql,
			params,
			method,
			connection_string,
		})
	})
}
