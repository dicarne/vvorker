import { Context, Hono } from "hono";
import { DebugEndpointRequest, ServiceBinding, TaskBinding } from "../types/debug-endpoint";
import { KVBinding } from "@dicarne/vvorker-kv";
import { PGSQLBinding } from "@dicarne/vvorker-pgsql";
import { OSSBinding } from "@dicarne/vvorker-oss";
import { isDev } from "../common/common";
import { MYSQLBinding } from "@dicarne/vvorker-mysql";
import { Base64 } from "js-base64";

/**
 * VITE_VVORKER_BASE_URL 设置为到服务的url
 * VITE_VVORKER_TOKEN 设置为服务的token
 * @param app0
 * @returns
 */
export function useDebugEndpoint(app0: any) {
  if (!isDev()) return;
  let app = new Hono();
  app.post("/", async (c: Context) => {
    const req = await c.req.json<DebugEndpointRequest>();
    switch (req.service) {
      case "oss": {
        let oss = (c.env as any)[req.binding] as OSSBinding;
        if (!oss) {
          return c.json({ error: "oss binding not found", req }, 404);
        }
        switch (req.method) {
          case "listBuckets":
            return c.json({ message: "oss", data: await oss.listBuckets() });
          case "listObjects":
            return c.json({
              message: "oss",
              data: await oss.listObjects(req.params.bucket),
            });
          case "downloadFile":
            return c.json({
              message: "oss",
              data: Base64.fromUint8Array(
                await oss.downloadFile(req.params.fileName)
              ),
            });
          case "uploadFile": {
            let base64 = req.params.data;
            const bytes = Base64.toUint8Array(base64);
            const file = await oss.uploadFile(bytes, req.params.fileName);
            return c.json({
              message: "oss",
              data: file,
            });
          }
          case "uploadStreamFile": {
            let base64 = req.params.stream;
            const bytes = Base64.toUint8Array(base64);
            let stream = new ReadableStream({
              start(controller) {
                controller.enqueue(bytes);
                controller.close();
              },
            });
            return c.json({
              message: "oss",
              data: await oss.uploadStreamFile(stream, req.params.fileName),
            });
          }
          case "downloadStreamFile":
            return c.json({
              message: "oss",
              data: await oss.downloadStreamFile(req.params.fileName),
            });
          case "deleteObject":
            return c.json({
              message: "oss",
              data: await oss.deleteObject(req.params.fileName),
            });
          default:
            return c.json({ error: "method not found", req }, 404);
        }
      }
      case "pgsql": {
        let pgsql = (c.env as any)[req.binding] as PGSQLBinding;
        if (!pgsql) {
          return c.json({ error: "pgsql binding not found", req }, 404);
        }
        let client = await pgsql.client();
        switch (req.method) {
          case "query":
            return c.json({
              message: "pgsql",
              data: await client.query(req.params.sql),
            });
          case "connectionString":
            return c.json({
              message: "pgsql",
              data: await pgsql.connectionString(),
            });
          case "connectionInfo":
            return c.json({
              message: "pgsql",
              data: await pgsql.connectionInfo(),
            });
          case "querysql":
            return c.json({
              message: "pgsql",
              data: await pgsql.query(
                req.params.sql,
                req.params.params,
                req.params.method
              ),
            });
          default:
            return c.json({ error: "method not found", req }, 404);
        }
      }
      case "mysql": {
        let mysql = (c.env as any)[req.binding] as MYSQLBinding;
        if (!mysql) {
          return c.json({ error: "mysql binding not found", req }, 404);
        }
        switch (req.method) {
          case "connectionString":
            return c.json({
              message: "mysql",
              data: await mysql.connectionString(),
            });
          case "connectionInfo":
            return c.json({
              message: "mysql",
              data: await mysql.connectionInfo(),
            });
          case "query":
            return c.json({
              message: "mysql",
              data: await mysql.query(
                req.params.sql,
                req.params.params,
                req.params.method
              ),
            });
          default:
            return c.json({ error: "method not found", req }, 404);
        }
      }
      case "kv": {
        let client = (c.env as any)[req.binding] as KVBinding;
        if (!client) {
          return c.json({ error: "kv binding not found", req }, 404);
        }
        switch (req.method) {
          case "get":
            return c.json({
              message: "kv",
              data: await client.get(req.params.key),
            });
          case "set":
            return c.json({
              message: "kv",
              data: await client.set(
                req.params.key,
                req.params.value,
                req.params.options
              ),
            });
          case "del":
            return c.json({
              message: "kv",
              data: await client.del(req.params.key),
            });
          case "keys":
            return c.json({
              message: "kv",
              data: await client.keys(
                req.params.pattern,
                req.params.offset,
                req.params.size
              ),
            });
          default:
            return c.json({ error: "method not found", req }, 404);
        }
      }
      case "vars": {
        switch (req.method) {
          case "get":
            return c.json({ message: "vars", data: await c.env.vars });
          default:
            return c.json({ error: "method not found", req }, 404);
        }
      }
      case "service": {
        let service = (c.env as any)[req.binding] as ServiceBinding;
        if (!service) {
          return c.json({ error: "service binding not found", req }, 404);
        }
        switch (req.method) {
          case "fetch":
            return service.fetch(req.params.url, req.params.init);
          default:
            return c.json({ error: "method not found", req }, 404);
        }
      }
      case "assets": {
        let assets = (c.env as any)[req.binding] as Fetcher;
        if (!assets) {
          return c.json({ error: "assets binding not found", req }, 404);
        }
        switch (req.method) {
          case "fetch":
            return assets.fetch(req.params.url, req.params.init);
          default:
            return c.json({ error: "method not found", req }, 404);
        }
      }
      case "proxy": {
        let proxy = (c.env as any)[req.binding] as ServiceBinding;
        if (!proxy) {
          return c.json({ error: "proxy binding not found", req }, 404);
        }
        switch (req.method) {
          case "fetch":
            return proxy.fetch(req.params.url, req.params.init);
          default:
            return c.json({ error: "method not found", req }, 404);
        }
      }
      case "task": {
        let task = (c.env as any)[req.binding] as TaskBinding;
        if (!task) {
          return c.json({ error: "task binding not found", req }, 404);
        }
        switch (req.method) {
          case "create": {
            const traceId = await task.create(req.params.name, req.params.trace_id);
            return c.json({
              code: 0,
              message: "task",
              data: { trace_id: traceId },
            });
          }
          case "should_exit": {
            const shouldExit = await task.should_exit(req.params.trace_id);
            return c.json({
              message: "task",
              data: shouldExit,
            });
          }
          case "complete": {
            await task.complete(req.params.trace_id);
            return c.json({ message: "task", data: null });
          }
          case "log": {
            await task.log(req.params.trace_id, req.params.text);
            return c.json({ message: "task", data: null });
          }
          default:
            return c.json({ error: "method not found", req }, 404);
        }
      }
    }
  });

  app0.route("/__vvorker__debug", app);
}
