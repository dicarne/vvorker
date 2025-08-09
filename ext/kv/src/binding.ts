export interface KVBinding {
    get(key: string): Promise<string | null>;
    set(key: string, value: string, options?: {
        EX?: number,
        NX?: boolean,
        XX?: boolean
    } | number): Promise<number>;
    del(key: string): Promise<void>;
    keys(pattern: string, offset: number, size: number): Promise<string[]>;
}