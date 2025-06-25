export interface DebugEndpointRequest {
    service: "oss" | "pgsql" | "kv" | "vars";
    binding: string;
    method: string;
    params: any;
}