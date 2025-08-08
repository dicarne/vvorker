
  import { Hono } from "hono";
  import { EnvBinding } from "./binding";
  import { init, useDebugEndpoint } from "@dicarne/vvorker-sdk";
  import { env } from "cloudflare:workers";
  init(env)
  
  const app = new Hono<{ Bindings: EnvBinding }>();
  useDebugEndpoint(app)
  
  app.get("/", (c) => {
    return c.text("Hello World!!!!!!!2222!!")
  })
  
  export default app;  
    