import { env } from "cloudflare:workers";

export default function () {
	return {
		async invoke(url: string, init: any) {
			return "hello text";
		}
	}

}

