export interface DebugEndpointRequest {
  service:
    | "oss"
    | "pgsql"
    | "kv"
    | "vars"
    | "mysql"
    | "service"
    | "proxy"
    | "assets"
    | "task";
  binding: string;
  method: string;
  params: any;
}

export interface ServiceBinding {
  fetch: (url: string, init?: RequestInit) => Promise<Response>;
}

export interface TaskBinding {
  create: (trace_id?: string) => Promise<string | undefined>;
  should_exit: (trace_id: string) => Promise<boolean>;
  complete: (trace_id: string) => Promise<void>;
  log: (trace_id: string, text: string) => Promise<void>;
}
