// filepath: src/index.ts
import { Client } from "pg";

import { WorkerEntrypoint, env } from 'cloudflare:workers'

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

export default class PGSQL extends WorkerEntrypoint {
	constructor(ctx: any, env: any) {
		super(ctx, env)
	}
	async start() {
		const client = new Client({
			user: cfg.user,
			host: cfg.host,
			database: cfg.database,
			password: cfg.password,
			port: Number(cfg.port),
		});
		await client.connect()
		return {
			end: async () => {
				await client.end()
			},
			query: async (sql: string, params: any[] = []) => {
				const result = await client.query(sql, params)
				return {
					rows: result.rows,
					rowCount: result.rowCount,
					command: result.command,
					oid: result.oid,
				}
			}
		}
	}
	connectionString() {
		return `postgres://${cfg.user}:${cfg.password}@${cfg.host}:${cfg.port}/${cfg.database}`;
	}
}