import { Command } from 'commander';
import inquirer from 'inquirer';
import * as fs from 'fs-extra';
import * as path from 'path';
import json5 from 'json5';
import { runCommand } from '../utils/system';

export const initCommand = new Command('init')
  .command('init <projectName>')
  .description('初始化VVorker项目')
  .action(async (projectName) => {
    // 交互式输入uid
    const { uid, projtype } = await inquirer.prompt([{
      type: 'input',
      name: 'uid',
      message: '请输入vvorker平台worker的uid:',
    }, {
      type: 'list',
      name: 'projtype',
      message: '请选择工程类型',
      choices: [
        { name: "纯Worker工程", value: "worker", description: "纯Worker工程" },
        { name: "Vue工程", value: "vue", description: "Vue工程，前后端分离" },
      ],
      default: 'worker'
    },])

    const jsonData = {
      "$schema": "vvorker-schema.json",
      "name": projectName,
      "project": {
        "uid": uid,
        "type": projtype
      },
      "version": "1.0.0",
      "extensions": [],
      "services": [],
      "vars": {},
      "ai": [],
      "oss": [],
      "pgsql": [],
      "kv": [],
    };

    // 如果当前文件夹下不存在package.json 
    if (!fs.existsSync('package.json')) {

      if (projtype === "worker") {
        // 创建项目目录
        await fs.ensureDir(projectName);

        // 定义要写入的 JSON 数据

        const gitignoreContent = `node_modules\n`

        const defaultJSCode = `
export default {
  async fetch(req, env) {
    try {
		let resp = new Response("worker: " + req.url + " is online! -- " + new Date())
		return resp
	} catch(e) {
		return new Response(e.stack, { status: 500 })
	}
  }
};
    `

        const packageJson = {
          "name": projectName,
          "version": "1.0.0",
          "description": "",
          "private": true,
          "scripts": {
            "deploy": "vvcli deploy",
            "dev": "wrangler dev",
            "start": "wrangler dev",
            "test": "vitest",
            "cf-typegen": "wrangler types",
            "build": "wrangler deploy --dry-run --outdir dist"
          },
          "devDependencies": {
            "@cloudflare/vitest-pool-workers": "^0.8.19",
            "typescript": "^5.5.2",
            "vitest": "~3.0.7",
            "wrangler": "^4.15.2"
          },
        }

        const wranglerConfig =
        {
          "$schema": "node_modules/wrangler/config-schema.json",
          "name": projectName,
          "main": "src/index.ts",
          "compatibility_date": "2025-05-20",
          "observability": {
            "enabled": true
          }
        }

        const jsonFilePath = `${projectName}/vvorker.json`;
        await fs.writeJson(jsonFilePath, jsonData, { spaces: 2 });

        const gitignoreFilePath = `${projectName}/.gitignore`;
        await fs.writeFile(gitignoreFilePath, gitignoreContent);

        const jsFilePath = `${projectName}/src/index.ts`;
        await fs.ensureDir(`${projectName}/src`);
        await fs.writeFile(jsFilePath, defaultJSCode);

        const packageJsonFilePath = `${projectName}/package.json`;
        await fs.writeJson(packageJsonFilePath, packageJson, { spaces: 2 });

        const wranglerConfigFilePath = `${projectName}/wrangler.jsonc`;
        await fs.writeJson(wranglerConfigFilePath, wranglerConfig);

        console.log(`项目 ${projectName} 初始化完成`);
        console.log(`请执行以下命令开始开发：`);
        console.log(`  cd ${projectName}`);
        console.log(`  pnpm install`);
      } else if (projtype === "vue") {
        const vueJSCode = `
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
        `

        await runCommand('pnpm', ['create', "cloudflare@latest", projectName, "--framework=vue"]);

        const jsFilePath = `${projectName}/server/index.ts`;
        await fs.writeFile(jsFilePath, vueJSCode);

        (jsonData as any)["assets"] = [
          {
            "directory": "./dist/client",
            "binding": "ASSETS"
          }
        ]
      }

      const jsonFilePath = `vvorker.json`;

      await fs.writeJson(path.join(projectName, jsonFilePath), jsonData, { spaces: 2 });
      console.log(`项目 ${projectName} 初始化完成`);
      console.log(`运行 vvcli types 生成相关类型提示`);
    } else {
      if (fs.existsSync('wrangler.jsonc')) {
        let wrtxt = await fs.readFile('wrangler.jsonc', 'utf-8');
        let jsonb = await json5.parse(wrtxt);
        wrtxt = wrtxt.replace(`"name": "${jsonb.name}"`, `"name": "${projectName}"`)
        await fs.writeFile('wrangler.jsonc', wrtxt);
      } else {

      }
      // 读取package.json，并将name改为proj name
      const packageJson = await fs.readJson('package.json');
      packageJson.name = projectName;
      await fs.writeJson('package.json', packageJson, { spaces: 2 });
      const jsonFilePath = `vvorker.json`;
      await fs.writeJson(jsonFilePath, jsonData, { spaces: 2 });
      console.log(`项目 ${projectName} 初始化完成`);
      console.log(`运行 vvcli types 生成相关类型提示`);
    }
  });
