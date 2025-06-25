export interface PGSQLBinding {
    connectionString: () =>Promise<string>;
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