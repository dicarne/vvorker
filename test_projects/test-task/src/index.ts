
import { Hono } from "hono";
import { EnvBinding } from "./binding";
import { init, useDebugEndpoint, vvbind } from "@dicarne/vvorker-sdk";
import { env } from "cloudflare:workers";
init(env)

const app = new Hono<{ Bindings: EnvBinding }>();
useDebugEndpoint(app)

app.get("*", async (c) => {
  console.log("start")
  const task = vvbind(c).task("task")
  const traceId = await task.create()
  if (!traceId) {
    return c.json({ code: 500, msg: "failed to create task", data: null })
  }
  await task.log(traceId, "new task! 111")
  c.executionCtx.waitUntil(new Promise(resolve => setTimeout(async () => {
    for (let i = 0; i < 100; i++) {
      await task.log(traceId, `log ${i}`)
    }

    await task.complete(traceId)
    resolve(0)
  }, 5000)))
  return c.json({
    code: 200,
    msg: "ok",
    data: { trace_id: traceId }
  })
})

export default app;
