export interface KVBinding {
    client: () => Promise<KV>;
}

export interface KV {
    get(key: string): Promise<string>;
    set(key: string, value: string): Promise<void>;
    del(key: string): Promise<void>;
    keys(pattern: string): Promise<string[]>;

}