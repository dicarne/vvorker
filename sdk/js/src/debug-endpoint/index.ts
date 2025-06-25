import { Hono } from "hono";
import { DebugEndpointRequest } from "../types/debug-endpoint";

export const debugEndpoint = new Hono();

debugEndpoint.post("/__vvorker__debug", async (c) => {
    const req = await c.req.json<DebugEndpointRequest>();
    switch (req.service) {
        case "oss":
            return c.json({ message: "oss" });
        case "pgsql":
            return c.json({ message: "pgsql" });
        case "kv":
            return c.json({ message: "kv" });
    }
});
