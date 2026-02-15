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
  client: () => Promise<TaskRpcTarget>;
  getTask: (trace_id: string) => Promise<TaskRpcTarget>;
}

export interface TaskRpcTarget {
  should_exit: () => Promise<boolean>;
  complete: () => Promise<void>;
  log: (text: string) => Promise<void>;
}
