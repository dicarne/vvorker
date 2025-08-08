import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import { cloudflare } from "@cloudflare/vite-plugin"

// https://vite.dev/config/
export default defineConfig({
	plugins: [
		cloudflare(),
	],
	base: "./",
	resolve: {
		alias: {
			'@': fileURLToPath(new URL('./src', import.meta.url))
		},
	},
})
