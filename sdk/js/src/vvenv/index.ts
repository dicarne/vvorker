import { KVBinding } from "@dicarne/vvorker-kv";
import { OSSBinding } from "@dicarne/vvorker-oss";
import { PGSQLBinding } from "@dicarne/vvorker-pgsql";
import { MYSQLBinding } from "@dicarne/vvorker-mysql";
import { config, isLocalDev } from "../common/common";
import { ServiceBinding, TaskBinding } from "../types/debug-endpoint";
import { Base64 } from "js-base64";

function vvoss(key: string, binding: OSSBinding): OSSBinding {
  if (isLocalDev()) {
    return {
      listBuckets: async () => {
        const r = await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "oss",
            binding: key,
            method: "listBuckets",
            params: {},
          }),
        });
        return ((await r.json()) as any).data;
      },
      listObjects: async (bucket: string) => {
        const r = await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "oss",
            binding: key,
            method: "listObjects",
            params: {
              bucket,
            },
          }),
        });
        return ((await r.json()) as any).data;
      },
      downloadFile: async (fileName: string) => {
        const r = await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "oss",
            binding: key,
            method: "downloadFile",
            params: {
              fileName,
            },
          }),
        });
        return Base64.toUint8Array(((await r.json()) as any).data);
      },
      uploadFile: async (data: Uint8Array, fileName: string) => {
        const base64 = Base64.fromUint8Array(data);
        const r = await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "oss",
            binding: key,
            method: "uploadFile",
            params: {
              data: base64,
              fileName,
            },
          }),
        });
        return ((await r.json()) as any).data;
      },
      uploadStreamFile: async (
        stream: ReadableStream<Uint8Array>,
        fileName: string
      ) => {
        const r = await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "oss",
            binding: key,
            method: "uploadStreamFile",
            params: {
              stream: Base64.fromUint8Array(
                (await stream.getReader().read()).value ?? new Uint8Array()
              ),
              fileName,
            },
          }),
        });
        return ((await r.json()) as any).data;
      },
      downloadStreamFile: async (fileName: string) => {
        const r = await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "oss",
            binding: key,
            method: "downloadStreamFile",
            params: {
              fileName,
            },
          }),
        });
        return ((await r.json()) as any).data;
      },
      deleteObject: async (fileName: string) => {
        const r = await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "oss",
            binding: key,
            method: "deleteObject",
            params: {
              fileName,
            },
          }),
        });
        return ((await r.json()) as any).data;
      },
    };
  } else {
    return binding;
  }
}

function vvpgsql(key: string, binding: PGSQLBinding): PGSQLBinding {
  if (isLocalDev()) {
    return {
      client: () =>
        Promise.resolve({
          query: async (sql: string) => {
            const r = await fetch(`${config().url}/__vvorker__debug`, {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${config().token}`,
              },
              body: JSON.stringify({
                service: "pgsql",
                binding: key,
                method: "query",
                params: {
                  sql,
                },
              }),
            });
            return ((await r.json()) as any).data;
          },
        }),
      connectionString: async () => {
        const r = await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "pgsql",
            binding: key,
            method: "connectionString",
            params: {},
          }),
        });
        return ((await r.json()) as any).data;
      },
      connectionInfo: async () => {
        const r = await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "pgsql",
            binding: key,
            method: "connectionInfo",
            params: {},
          }),
        });
        return ((await r.json()) as any).data;
      },
      query: async (sql: string, params: any, method: string) => {
        const r = await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "pgsql",
            binding: key,
            method: "querysql",
            params: {
              sql,
              params,
              method,
            },
          }),
        });
        return ((await r.json()) as any).data;
      },
    };
  } else {
    return binding;
  }
}

function vvmysql(key: string, binding: MYSQLBinding): MYSQLBinding {
  if (isLocalDev()) {
    return {
      connectionString: async () => {
        const r = await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "mysql",
            binding: key,
            method: "connectionString",
            params: {},
          }),
        });
        return ((await r.json()) as any).data as string;
      },
      connectionInfo: async () => {
        const r = await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "mysql",
            binding: key,
            method: "connectionInfo",
            params: {},
          }),
        });
        let data = ((await r.json()) as any).data;
        return data;
      },
      query: async (sql: string, params: any, method: string) => {
        const r = await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "mysql",
            binding: key,
            method: "query",
            params: {
              sql,
              params,
              method,
            },
          }),
        });
        return ((await r.json()) as any).data;
      },
    };
  } else {
    return binding;
  }
}

function vvkv(binding_key: string, binding: KVBinding): KVBinding {
  if (isLocalDev()) {
    return {
      get: async (key: string) => {
        const r = await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "kv",
            binding: binding_key,
            method: "get",
            params: {
              key,
            },
          }),
        });
        return ((await r.json()) as any).data;
      },
      set: async (
        key: string,
        value: string,
        options?:
          | {
              EX?: number;
              NX?: boolean;
              XX?: boolean;
            }
          | number
      ) => {
        const r = await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "kv",
            binding: binding_key,
            method: "set",
            params: {
              key,
              value,
              options,
            },
          }),
        });
        return ((await r.json()) as any).data;
      },
      del: async (key: string) => {
        const r = await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "kv",
            binding: binding_key,
            method: "del",
            params: {
              key,
            },
          }),
        });
        return ((await r.json()) as any).data;
      },
      keys: async (pattern: string) => {
        const r = await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "kv",
            binding: binding_key,
            method: "keys",
            params: {
              pattern,
            },
          }),
        });
        return ((await r.json()) as any).data;
      },
    };
  } else {
    return binding;
  }
}

async function vars<T extends { vars: any }>(binding: any): Promise<T["vars"]> {
  if (isLocalDev()) {
    let r = await fetch(`${config().url}/__vvorker__debug`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${config().token}`,
      },
      body: JSON.stringify({
        service: "vars",
        binding: "",
        method: "get",
        params: {},
      }),
    });
    return ((await r.json()) as any).data;
  }
  return binding.vars;
}

/**
 * env: VITE_APP_UID 用于指定模仿单点登录的uid
 * @param key
 * @param binding
 * @returns
 */
function service(key: string, binding: ServiceBinding) {
  if (isLocalDev()) {
    return {
      fetch: async (path: string, init?: RequestInit) => {
        if (!init) {
          init = {};
        }
        if (!init.headers) {
          init.headers = {};
        }
        init.headers = {
          ...init.headers,
          "vvorker-worker-uid": ((import.meta as any).env.VITE_APP_UID ??
            "") as string,
        };
        let r = await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "service",
            binding: key,
            method: "fetch",
            params: {
              url:
                "http://vvorker.local" +
                (path.startsWith("/") ? path : "/" + path),
              init: init,
            },
          }),
        });
        return r;
      },
      /**
       * 简化调用方法，直接以json形式调用
       * @param path
       * @param data
       * @returns 响应体json，而不是Response
       */
      call: async (path: string, data?: any) => {
        let r = await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "service",
            binding: key,
            method: "fetch",
            params: {
              url:
                "http://vvorker.local" +
                (path.startsWith("/") ? path : "/" + path),
              init: {
                method: "POST",
                headers: {
                  "Content-Type": "application/json",
                  "vvorker-worker-uid": ((import.meta as any).env
                    .VITE_APP_UID ?? "") as string,
                },
                body: JSON.stringify(data),
              },
            },
          }),
        });
        return rpcResultWrap(r);
      },
    };
  }
  return {
    fetch: async (path: string, init?: RequestInit) =>
      binding.fetch(
        "http://vvorker.local" + (path.startsWith("/") ? path : "/" + path),
        init
      ),
    /**
     * 简化调用方法，直接以json形式调用
     * @param path
     * @param data
     * @returns 响应体json，而不是Response
     */
    call: async (path: string, data?: any) => {
      let r = await binding.fetch(
        "http://vvorker.local" + (path.startsWith("/") ? path : "/" + path),
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(data),
        }
      );
      return rpcResultWrap(r);
    },
  };
}

async function rpcResultWrap(r: Response) {
  if (r.status === 200) {
    const j: any = await r.json();
    if (j.code === 0) {
      return j;
    }
    throw new Error(`调用失败：${j.code} ${j.message} ${JSON.stringify(j)}`);
  }
  throw new Error(`调用失败：${r.status} ${r.statusText} ${await r.text()}`);
}

function proxy(key: string, binding: ServiceBinding, remote?: boolean) {
  if (isLocalDev()) {
    if (remote) {
      return {
        fetch: async (path: string, init?: RequestInit) => {
          if (!init) {
            init = {};
          }
          if (!init.headers) {
            init.headers = {};
          }
          init.headers = {
            ...init.headers,
            "vvorker-worker-uid": ((import.meta as any).env.VITE_APP_UID ??
              "") as string,
          };
          let r = await fetch(`${config().url}/__vvorker__debug`, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${config().token}`,
            },
            body: JSON.stringify({
              service: "proxy",
              binding: key,
              method: "fetch",
              params: {
                url:
                  "http://vvorker.local" +
                  (path.startsWith("/") ? path : "/" + path),
                init: init,
              },
            }),
          });
          return r;
        },
      };
    } else {
      return {
        fetch: (url: string, init?: RequestInit) => fetch(url, init),
      };
    }
  }

  return binding;
}

function assets(key: string, binding: Fetcher) {
  if (isLocalDev()) {
    return {
      fetch: async (path: string, init?: RequestInit) => {
        if (!init) {
          init = {};
        }
        if (!init.headers) {
          init.headers = {};
        }
        init.headers = {
          ...init.headers,
          "vvorker-worker-uid": ((import.meta as any).env.VITE_APP_UID ??
            "") as string,
        };
        let r = await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "assets",
            binding: key,
            method: "fetch",
            params: {
              url:
                "http://vvorker.local" +
                (path.startsWith("/") ? path : "/" + path),
              init: init,
            },
          }),
        });
        return r;
      },
    };
  }
  return {
    fetch: async (path: string, init?: RequestInit) =>
      binding.fetch(
        "http://vvorker.local" + (path.startsWith("/") ? path : "/" + path),
        init
      ),
  };
}

function vvtask(bindingKey: string, binding: TaskBinding): TaskBinding {
  if (isLocalDev()) {
    return {
      create: async (trace_id?: string) => {
        const r = await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "task",
            binding: bindingKey,
            method: "create",
            params: { trace_id },
          }),
        });
        const result = (await r.json()) as any;
        return result.data?.trace_id;
      },
      should_exit: async (trace_id: string) => {
        const r = await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "task",
            binding: bindingKey,
            method: "should_exit",
            params: { trace_id },
          }),
        });
        const result = (await r.json()) as any;
        return result.data === true;
      },
      complete: async (trace_id: string) => {
        await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "task",
            binding: bindingKey,
            method: "complete",
            params: { trace_id },
          }),
        });
      },
      log: async (trace_id: string, text: string) => {
        await fetch(`${config().url}/__vvorker__debug`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${config().token}`,
          },
          body: JSON.stringify({
            service: "task",
            binding: bindingKey,
            method: "log",
            params: { trace_id, text },
          }),
        });
      },
    };
  }
  return binding;
}

/**
 * 用于转换环境变量和绑定，在开发时（env.vars.MODE="development"）将通过代理和节点进行交互，从而获取节点的绑定和变量。
 * 在生产时，将直接返回绑定和变量。
 */
export function vvbind<T extends { env: { vars: any; [key: string]: any } }>(
  c: T
) {
  return {
    oss: (key: keyof T["env"]) => vvoss(key as string, c.env[key as string]),
    pgsql: (key: keyof T["env"]) =>
      vvpgsql(key as string, c.env[key as string]),
    mysql: (key: keyof T["env"]) =>
      vvmysql(key as string, c.env[key as string]),
    kv: (key: keyof T["env"]) => vvkv(key as string, c.env[key as string]),
    proxy: (key: keyof T["env"]) => proxy(key as string, c.env[key as string]),
    vars: () => vars<{ vars: T["env"]["vars"] }>(c.env),
    service: (name: keyof T["env"]) =>
      service(name as string, c.env[name as string]),
    assets: (key: keyof T["env"]) =>
      assets(key as string, c.env[key as string]),
    task: (key: keyof T["env"]) => vvtask(key as string, c.env[key as string]),
  };
}
