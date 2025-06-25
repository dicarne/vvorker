import { KVBinding, KV } from "@dicarne/vvorker-kv";
import { OSSBinding } from "@dicarne/vvorker-oss";
import { PGSQLBinding } from "@dicarne/vvorker-pgsql";
import { Context } from "hono";

function isDev() {
    return true
}

function vvoss(key: string, binding: OSSBinding): OSSBinding {
    if (isDev()) {
        return {
            listBuckets: () => Promise.resolve([]),
            listObjects: (bucket: string) => Promise.resolve([]),
            downloadFile: (fileName: string) => Promise.resolve(new Uint8Array()),
            uploadFile: (data: Uint8Array, fileName: string) => Promise.resolve(),
            uploadStreamFile: (stream: ReadableStream<Uint8Array>, fileName: string) => Promise.resolve(),
            downloadStreamFile: (fileName: string) => Promise.resolve(new ReadableStream<Uint8Array>()),
            deleteObject: (fileName: string) => Promise.resolve(),
        }
    } else {
        return binding
    }
}

function vvpgsql(key: string, binding: PGSQLBinding): PGSQLBinding {
    if (isDev()) {
        return {
            client: () => Promise.resolve({
                query: () => Promise.resolve({
                    rows: [],
                    rowCount: 0,
                    command: "",
                    oid: 0,
                }),
            }),
            connectionString: () => Promise.resolve(""),
        }
    } else {
        return binding
    }
}

function vvkv(key: string, binding: KVBinding): KVBinding {
    if (isDev()) {
        return {
            client: () => Promise.resolve({
                get: () => Promise.resolve(""),
                set: () => Promise.resolve(),
                del: () => Promise.resolve(),
                keys: () => Promise.resolve([]),
            })
        }
    } else {
        return binding
    }
}


export function vvbind(c: Context) {
    return {
        oss: (key: string) => vvoss(key, c.env[key]),
        pgsql: (key: string) => vvpgsql(key, c.env[key]),
        kv: (key: string) => vvkv(key, c.env[key]),
    }
}
