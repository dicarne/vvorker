export interface PGSQLBinding {
    connectionString: () =>Promise<string>;
    connectionInfo: () => Promise<{user: string, host: string, database: string, password: string, port: number}>;
    client: () => Promise<PGSQLClient>;
}

export interface PGSQLClient {
    query(sql: string): Promise<{
        rows: any[],
        rowCount: number
        command: string
        oid: number
    }>;
}