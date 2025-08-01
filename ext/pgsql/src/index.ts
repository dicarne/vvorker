// filepath: src/index.ts
import { Client } from "pg";
export * from "./binding"
import { RpcTarget, WorkerEntrypoint, env } from 'cloudflare:workers'

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

class PGSQLTarget extends RpcTarget {
	client: Client
	constructor() {
		super()
		this.client = new Client({
			user: cfg.user,
			host: cfg.host,
			database: cfg.database,
			password: cfg.password,
			port: Number(cfg.port),
		});
	}
	async start() {
		await this.client.connect()
	}
	async query(sql: string, params: any[] = []) {
		const result = await this.client.query(sql, params)
		return {
			rows: result.rows,
			rowCount: result.rowCount,
			command: result.command,
			oid: result.oid,
		}
	}
}

export default class PGSQL extends WorkerEntrypoint {
	constructor(ctx: any, env: any) {
		super(ctx, env)
	}
	async client() {
		const target = new PGSQLTarget()
		await target.start()
		return target
	}
	connectionString() {
		return `postgres://${cfg.user}:${encodeURIComponent(cfg.password)}@${cfg.host}:${cfg.port}/${cfg.database}`;
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
		return (await rpc(sql, params, method, this.connectionString() + "?sslmode=disable")).json()
	}
}



async function rpc(sql: string, params: any, method: string, connection_string: string) {
	return fetch(`${eenv.MASTER_ENDPOINT}/api/ext/pgsql/query`, {
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
