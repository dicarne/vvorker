export interface DebugEndpointRequest {
    service: "oss" | "pgsql" | "kv";
    binding: string;
    method: string;
    params: any;
}