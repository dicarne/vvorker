export interface DebugEndpointRequest {
    service: "oss" | "pgsql" | "kv" | "vars" | "mysql" | "service" | "proxy" | "assets";
    binding: string;
    method: string;
    params: any;
}

export interface ServiceBinding {
    fetch: (url: string, init?: RequestInit) => Promise<Response>
}