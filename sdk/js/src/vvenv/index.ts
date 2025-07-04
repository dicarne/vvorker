import { KVBinding, KV } from "@dicarne/vvorker-kv";
import { OSSBinding } from "@dicarne/vvorker-oss";
import { PGSQLBinding } from "@dicarne/vvorker-pgsql";
import { MYSQLBinding } from "@dicarne/vvorker-mysql";
import { config, isDev } from "../common/common";


function vvoss(key: string, binding: OSSBinding): OSSBinding {
    if (isDev()) {
        return {
            listBuckets: async () => {
                const r = await fetch(`${config().url}/__vvorker__debug`, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": `Bearer ${config().token}`
                    },
                    body: JSON.stringify({
                        service: "oss",
                        binding: key,
                        method: "listBuckets",
                        params: {}
                    })
                })
                return (await r.json()).data
            },
            listObjects: async (bucket: string) => {
                const r = await fetch(`${config().url}/__vvorker__debug`, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": `Bearer ${config().token}`
                    },
                    body: JSON.stringify({
                        service: "oss",
                        binding: key,
                        method: "listObjects",
                        params: {
                            bucket
                        }
                    })
                })
                return (await r.json()).data
            },
            downloadFile: async (fileName: string) => {
                const r = await fetch(`${config().url}/__vvorker__debug`, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": `Bearer ${config().token}`
                    },
                    body: JSON.stringify({
                        service: "oss",
                        binding: key,
                        method: "downloadFile",
                        params: {
                            fileName
                        }
                    })
                })
                return (await r.json()).data
            },
            uploadFile: async (data: Uint8Array, fileName: string) => {
                const r = await fetch(`${config().url}/__vvorker__debug`, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": `Bearer ${config().token}`
                    },
                    body: JSON.stringify({
                        service: "oss",
                        binding: key,
                        method: "uploadFile",
                        params: {
                            data,
                            fileName
                        }
                    })
                })
                return (await r.json()).data
            },
            uploadStreamFile: async (stream: ReadableStream<Uint8Array>, fileName: string) => {
                const r = await fetch(`${config().url}/__vvorker__debug`, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": `Bearer ${config().token}`
                    },
                    body: JSON.stringify({
                        service: "oss",
                        binding: key,
                        method: "uploadStreamFile",
                        params: {
                            stream,
                            fileName
                        }
                    })
                })
                return (await r.json()).data
            },
            downloadStreamFile: async (fileName: string) => {
                const r = await fetch(`${config().url}/__vvorker__debug`, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": `Bearer ${config().token}`
                    },
                    body: JSON.stringify({
                        service: "oss",
                        binding: key,
                        method: "downloadStreamFile",
                        params: {
                            fileName
                        }
                    })
                })
                return (await r.json()).data
            },
            deleteObject: async (fileName: string) => {
                const r = await fetch(`${config().url}/__vvorker__debug`, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": `Bearer ${config().token}`
                    },
                    body: JSON.stringify({
                        service: "oss",
                        binding: key,
                        method: "deleteObject",
                        params: {
                            fileName
                        }
                    })
                })
                return (await r.json()).data
            },
        }
    } else {
        return binding
    }
}

function vvpgsql(key: string, binding: PGSQLBinding): PGSQLBinding {
    if (isDev()) {
        return {
            client: () => Promise.resolve({
                query: async (sql: string) => {
                    const r = await fetch(`${config().url}/__vvorker__debug`, {
                        method: "POST",
                        headers: {
                            "Content-Type": "application/json",
                            "Authorization": `Bearer ${config().token}`
                        },
                        body: JSON.stringify({
                            service: "pgsql",
                            binding: key,
                            method: "query",
                            params: {
                                sql
                            }
                        })
                    })
                    return (await r.json()).data
                },
            }),
            connectionString: async () => {
                const r = await fetch(`${config().url}/__vvorker__debug`, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": `Bearer ${config().token}`
                    },
                    body: JSON.stringify({
                        service: "pgsql",
                        binding: key,
                        method: "connectionString",
                        params: {}
                    })
                })
                return (await r.json()).data
            },
        }
    } else {
        return binding
    }
}

function vvmysql(key: string, binding: MYSQLBinding): MYSQLBinding {
    if (isDev()) {
        return {
            connectionString: async () => {
                const r = await fetch(`${config().url}/__vvorker__debug`, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": `Bearer ${config().token}`
                    },
                    body: JSON.stringify({
                        service: "mysql",
                        binding: key,
                        method: "connectionString",
                        params: {}
                    })
                })
                return (await r.json()).data
            },
        }
    } else {
        return binding
    }
}

function vvkv(binding_key: string, binding: KVBinding): KVBinding {
    console.log(`${config().url}/__vvorker__debug`)
    if (isDev()) {
        return {
            client: () => Promise.resolve({
                get: async (key: string) => {
                    const r = await fetch(`${config().url}/__vvorker__debug`, {
                        method: "POST",
                        headers: {
                            "Content-Type": "application/json",
                            "Authorization": `Bearer ${config().token}`
                        },
                        body: JSON.stringify({
                            service: "kv",
                            binding: binding_key,
                            method: "get",
                            params: {
                                key
                            }
                        })
                    })
                    return (await r.json()).data
                },
                set: async (key: string, value: string) => {
                    const r = await fetch(`${config().url}/__vvorker__debug`, {
                        method: "POST",
                        headers: {
                            "Content-Type": "application/json",
                            "Authorization": `Bearer ${config().token}`
                        },
                        body: JSON.stringify({
                            service: "kv",
                            binding: binding_key,
                            method: "set",
                            params: {
                                key,
                                value
                            }
                        })
                    })
                    return (await r.json()).data
                },
                del: async (key: string) => {
                    const r = await fetch(`${config().url}/__vvorker__debug`, {
                        method: "POST",
                        headers: {
                            "Content-Type": "application/json",
                            "Authorization": `Bearer ${config().token}`
                        },
                        body: JSON.stringify({
                            service: "kv",
                            binding: binding_key,
                            method: "del",
                            params: {
                                key
                            }
                        })
                    })
                    return (await r.json()).data
                },
                keys: async (pattern: string) => {
                    const r = await fetch(`${config().url}/__vvorker__debug`, {
                        method: "POST",
                        headers: {
                            "Content-Type": "application/json",
                            "Authorization": `Bearer ${config().token}`
                        },
                        body: JSON.stringify({
                            service: "kv",
                            binding: binding_key,
                            method: "keys",
                            params: {
                                pattern
                            }
                        })
                    })
                    return (await r.json()).data
                },
            })
        }
    } else {
        return binding
    }
}

async function vars<T extends { vars: any }>(binding: any): Promise<T['vars']> {
    if (isDev()) {
        let r = await fetch(`${config().url}/__vvorker__debug`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${config().token}`
            },
            body: JSON.stringify({
                service: "vars",
                binding: "",
                method: "get",
                params: {}
            })
        })
        return (await r.json()).data
    }
    return binding.vars
}

/**
 * 用于转换环境变量和绑定，在开发时（env.vars.MODE="development"）将通过代理和节点进行交互，从而获取节点的绑定和变量。
 * 在生产时，将直接返回绑定和变量。
 */
export function vvbind<T extends { env: { vars: any, [key: string]: any } }>(c: T) {
    return {
        oss: (key: string) => vvoss(key, c.env[key]),
        pgsql: (key: string) => vvpgsql(key, c.env[key]),
        mysql: (key: string) => vvmysql(key, c.env[key]),
        kv: (key: string) => vvkv(key, c.env[key]),
        vars: () => vars<{ vars: T['env']['vars'] }>(c.env)
    }
}
