import { Hono } from "hono";
import { EnvBinding } from "./binding";

const app = new Hono<{ Bindings: EnvBinding }>();

app.get("*", async (c) => {
	try {
		const r = await c.env.ASSETS.fetch(c.req.url, c.req)
		const url = new URL(c.req.url);
		if (r.status === 404) {
			return c.env.ASSETS.fetch("https://" + url.host + "/index.html", c.req)
		}
		return r
	} catch (error) {
		c.status(404);
	}

});


export default app;
