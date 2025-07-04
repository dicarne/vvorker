export interface DebugEndpointRequest {
    service: "oss" | "pgsql" | "kv" | "vars" | "mysql";
    binding: string;
    method: string;
    params: any;
}