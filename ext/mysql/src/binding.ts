export interface MYSQLBinding {
    connectionString: () => Promise<string>;
    connectionInfo: () => Promise<{user: string, host: string, database: string, password: string, port: number}>;
}
