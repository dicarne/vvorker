import { Context, Env, Hono } from "hono";
import { DebugEndpointRequest, ServiceBinding } from "../types/debug-endpoint";
import { KVBinding } from "@dicarne/vvorker-kv";
import { PGSQLBinding } from "@dicarne/vvorker-pgsql";
import { OSSBinding } from "@dicarne/vvorker-oss";
import { isDev } from "../common/common";
import { MYSQLBinding } from "@dicarne/vvorker-mysql";

export function useDebugEndpoint(app: any) {
    if (!isDev()) return
    app.post("/__vvorker__debug", async (c: Context) => {
        const req = await c.req.json<DebugEndpointRequest>();
        switch (req.service) {
            case "oss":
                {
                    let oss = ((c.env as any)[req.binding] as OSSBinding)
                    if (!oss) {
                        return c.json({ error: "oss binding not found", req }, 404)
                    }
                    switch (req.method) {
                        case "listBuckets":
                            return c.json({ message: "oss", data: await oss.listBuckets() });
                        case "listObjects":
                            return c.json({ message: "oss", data: await oss.listObjects(req.params.bucket) });
                        case "downloadFile":
                            return c.json({ message: "oss", data: await oss.downloadFile(req.params.fileName) });
                        case "uploadFile":
                            return c.json({ message: "oss", data: await oss.uploadFile(req.params.data, req.params.fileName) });
                        case "uploadStreamFile":
                            return c.json({ message: "oss", data: await oss.uploadStreamFile(req.params.stream, req.params.fileName) });
                        case "downloadStreamFile":
                            return c.json({ message: "oss", data: await oss.downloadStreamFile(req.params.fileName) });
                        case "deleteObject":
                            return c.json({ message: "oss", data: await oss.deleteObject(req.params.fileName) });
                        default:
                            return c.json({ error: "method not found", req }, 404)
                    }
                }
            case "pgsql":
                {
                    let pgsql = ((c.env as any)[req.binding] as PGSQLBinding)
                    if (!pgsql) {
                        return c.json({ error: "pgsql binding not found", req }, 404)
                    }
                    let client = await pgsql.client()
                    switch (req.method) {
                        case "query":
                            return c.json({ message: "pgsql", data: await client.query(req.params.sql) });
                        case "connectionString":
                            return c.json({ message: "pgsql", data: await pgsql.connectionString() });
                        default:
                            return c.json({ error: "method not found", req }, 404)
                    }
                }
            case "mysql":
                {
                    let mysql = ((c.env as any)[req.binding] as MYSQLBinding)
                    if (!mysql) {
                        return c.json({ error: "mysql binding not found", req }, 404)
                    }
                    switch (req.method) {
                        case "connectionString":
                            return c.json({ message: "mysql", data: await mysql.connectionString() });
                        default:
                            return c.json({ error: "method not found", req }, 404)
                    }
                }
            case "kv":
                {
                    let kv = ((c.env as any)[req.binding] as KVBinding)
                    if (!kv) {
                        return c.json({ error: "kv binding not found", req }, 404)
                    }
                    let client = await kv.client()
                    switch (req.method) {
                        case "get":
                            return c.json({ message: "kv", data: await client.get(req.params.key) });
                        case "set":
                            return c.json({ message: "kv", data: await client.set(req.params.key, req.params.value) });
                        case "del":
                            return c.json({ message: "kv", data: await client.del(req.params.key) });
                        case "keys":
                            return c.json({ message: "kv", data: await client.keys(req.params.pattern) });
                        default:
                            return c.json({ error: "method not found", req }, 404)
                    }
                }
            case "vars":
                {
                    switch (req.method) {
                        case "get":
                            return c.json({ message: "vars", data: await c.env.vars });
                        default:
                            return c.json({ error: "method not found", req }, 404)
                    }
                }
            case "service":
                {
                    let service = ((c.env as any)[req.binding] as ServiceBinding)
                    if (!service) {
                        return c.json({ error: "service binding not found", req }, 404)
                    }
                    switch (req.method) {
                        case "fetch":
                            return c.json({ message: "service", data: await (await service.fetch(req.params.path, req.params.init)).json() });
                        default:
                            return c.json({ error: "method not found", req }, 404)
                    }
                }
        }
    })
}