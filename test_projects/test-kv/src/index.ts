
import { Hono } from "hono";
import { EnvBinding } from "./binding";
import { init, useDebugEndpoint, vvbind } from "@dicarne/vvorker-sdk";
import { env } from "cloudflare:workers";
init(env)

const app = new Hono<{ Bindings: EnvBinding }>();
useDebugEndpoint(app)

app.get("/", async (c) => {
  const kv = vvbind(c).kv("kv")
  await kv.set("hi", "hello")
  return c.text((await kv.get("hi")) || "[NULL]")
})

export default app;
