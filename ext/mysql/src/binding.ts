export interface MYSQLBinding {
    connectionString: () => Promise<string>;
}
