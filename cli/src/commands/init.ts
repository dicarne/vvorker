import { Command } from 'commander';
import inquirer from 'inquirer';
import * as fs from 'fs-extra';
import * as path from 'path';
import json5 from 'json5';
import { runCommand } from '../utils/system';
import pc from "picocolors"
import { config, getEnv } from '../utils/config';

async function createWorkerProject(projectName: string, jsonData: object) {
  await fs.ensureDir(projectName);

  const gitignoreContent = `node_modules\n`;

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
    `;

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
    "dependencies": {
      "@dicarne/vvorker-sdk": "^0.1.10",
      "hono": "^4.7.11"
    },
    "devDependencies": {
      "@cloudflare/vitest-pool-workers": "^0.8.19",
      "typescript": "^5.5.2",
      "vitest": "~3.0.7",
      "wrangler": "^4.15.2"
    },
  };

  const wranglerConfig = {
    "$schema": "node_modules/wrangler/config-schema.json",
    "name": projectName,
    "main": "src/index.ts",
    "compatibility_date": "2025-05-20",
    "observability": {
      "enabled": true
    }
  };

  const projectPath = projectName;
  await fs.writeJson(path.join(projectPath, `vvorker.${getEnv()}.json`), jsonData, { spaces: 2 });
  await fs.writeFile(path.join(projectPath, '.gitignore'), gitignoreContent);
  await fs.ensureDir(path.join(projectPath, 'src'));
  await fs.writeFile(path.join(projectPath, 'src', 'index.ts'), defaultJSCode);
  await fs.writeJson(path.join(projectPath, 'package.json'), packageJson, { spaces: 2 });
  await fs.writeJson(path.join(projectPath, 'wrangler.jsonc'), wranglerConfig);

  console.log(`项目 ${projectName} 初始化完成`);
  console.log(`请执行以下命令开始开发：`);
  console.log(`  cd ${projectName}`);
  console.log(`  pnpm install`);
}

async function createVueProject(projectName: string, jsonData: object) {
  const vueJSCode = `
import { Hono } from "hono";
import { EnvBinding } from "./binding";

const app = new Hono<{ Bindings: EnvBinding }>();

app.notFound(async (c) => {
	try {
		const r = await c.env.ASSETS.fetch(c.req.url, c.req)
		const url = new URL(c.req.url);
		if (r.status === 404) {
			return c.env.ASSETS.fetch("https://" + url.host + "/index.html", c.req)
		}
		return r
	} catch (error) {
		return c.text("404 Not Found", 404)
	}

});

export default app;
    `;

  await runCommand('pnpm', ['create', "cloudflare@latest", projectName, "--framework=vue"]);


  const jsFilePath = path.join(projectName, 'server', 'index.ts');
  await fs.writeFile(jsFilePath, vueJSCode);

  (jsonData as any)["assets"] = [
    {
      "directory": "./dist/client",
      "binding": "ASSETS"
    }
  ];
  const env = getEnv();
  const jsonFilePath = `vvorker.${env}.json`;
  await fs.writeJson(path.join(projectName, jsonFilePath), jsonData, { spaces: 2 });

  console.log(pc.green(`项目 ${projectName} 初始化完成`));

  try {
    await runCommand('pnpm', ['install'], projectName);
    await runCommand('pnpm', ['install', "hono", "@dicarne/vvorker-sdk"], projectName);
  } catch (error) {
    console.log(pc.red(`安装依赖失败，请手动安装`));
  }
  try {
    await runCommand('vvcli', ['types'], projectName);
  } catch (error) {
    console.log(pc.red(`生成类型提示失败，请手动运行 vvcli types 生成`));
  }
}

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
      "mysql": [],
      "kv": [],
      "assets": [],
    };

    switch (projtype) {
      case "worker": {
        await createWorkerProject(projectName, jsonData);
        break;
      }
      case "vue": {
        await createVueProject(projectName, jsonData);
        break;
      }
    }
  });
