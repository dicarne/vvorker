// filepath: src/index.ts
import { Client } from "pg";

import { WorkerEntrypoint, env } from 'cloudflare:workers'

function config() {
    // 从环境变量中获取配置信息
    let cfg = {
        "user": env.USER,
        "host": env.HOST,
        "port": env.PORT,
        "password": env.PASSWORD,
        "database": env.DATABASE,
    }

    // 遍历配置对象，检查每个属性是否为空
    for (const [key, value] of Object.entries(cfg)) {
        if (!value) {
            throw new Error(`Environment variable ${key.toUpperCase()} is missing or empty`);
        }
    }

    return cfg
}

export default class PGSQL extends WorkerEntrypoint {

	constructor(ctx: any, env: any) {
		super(ctx, env)
	}
	async call() {
		let cfg = config()
		// Create a new client instance for each request.
		const client = new Client({
			user: cfg.user,
			host: cfg.host,
			database: cfg.database,
			password: cfg.password,
			port: Number(cfg.port),
		});

		// Connect to the database
		await client.connect();
		console.log("Connected to PostgreSQL database");

		// Perform a simple query
		const result = await client.query("SELECT * FROM pg_tables");
		
		return {
			rows: result.rows,
			rowCount: result.rowCount,
			command: result.command,
			oid: result.oid,
		}

	}
}