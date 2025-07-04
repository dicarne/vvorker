export interface KVBinding {
    client: () => Promise<KV>;
}

export interface KV {
    get(key: string): Promise<string>;
    set(key: string, value: string, ttl?: number): Promise<void>;
    del(key: string): Promise<void>;
    keys(pattern: string, offset: number, size: number): Promise<string[]>;

}