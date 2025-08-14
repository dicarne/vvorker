import { Command } from 'commander';
import inquirer from 'inquirer';
import * as fs from 'fs-extra';
import * as path from 'path';
import json5 from 'json5';
import { runCommand } from '../utils/system';
import pc from "picocolors"
import { config, getEnv } from '../utils/config';

async function createWorkerProject(projectName: string, jsonData: object) {
  const vueJSCode = `
import { Hono } from "hono";
import { EnvBinding } from "./binding";
import { init, useDebugEndpoint } from "@dicarne/vvorker-sdk";
import { env } from "cloudflare:workers";
init(env)

const app = new Hono<{ Bindings: EnvBinding }>();
useDebugEndpoint(app)

app.onError(async (err, c) => {
  console.error(err)
  return c.json({
    code: 500,
    msg: err.message,
    data: null
  }, 500)
})
  
export default app;
`;

  await runCommand('pnpm', ['create', "cloudflare@latest", projectName, "--template=cloudflare/templates/hello-world-do-template", "--git", "--no-deploy", "--lang=ts"]);


  const jsFilePath = path.join(projectName, 'src', 'index.ts');
  await fs.writeFile(jsFilePath, vueJSCode);

  const env = getEnv();
  const jsonFilePath = `vvorker.${env}.json`;
  await fs.writeJson(path.join(projectName, jsonFilePath), jsonData, { spaces: 2 });

  const wranglerJsonPath = path.join(projectName, 'wrangler.json');
  const wranglerJson = json5.parse(await fs.readFile(wranglerJsonPath, 'utf-8'));
  wranglerJson.compatibility_flags = ["nodejs_compat"];
  wranglerJson.durable_objects = undefined;
  wranglerJson.migrations = undefined;
  await fs.writeJson(wranglerJsonPath, wranglerJson, { spaces: 2 });

  const packageJsonPath = path.join(projectName, 'package.json');
  const packageJson = json5.parse(await fs.readFile(packageJsonPath, 'utf-8'));
  packageJson.type = "module";
  packageJson.name = projectName;
  packageJson.scripts.dev = "vite";
  packageJson.scripts.build = "vite build";
  packageJson.scripts.start = undefined;
  packageJson.scripts.deploy = undefined;
  await fs.writeJson(packageJsonPath, packageJson, { spaces: 2 });

  const viteConfigText = `import { fileURLToPath, URL } from 'node:url'
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
`

  await fs.writeFile(path.join(projectName, 'vite.config.ts'), viteConfigText);

  console.log(pc.green(`项目 ${projectName} 初始化完成`));

  try {
    await runCommand('pnpm', ['install'], projectName);
    await runCommand('pnpm', ['install', "@types/node", "vite", "@cloudflare/vite-plugin", "-D"], projectName);
    await runCommand('pnpm', ['install', "hono", "@dicarne/vvorker-sdk", "@hono/zod-validator", "zod"], projectName);
  } catch (error) {
    console.log(pc.red(`安装依赖失败，请手动安装`));
  }
  try {
    await runCommand('vvcli', ['types'], projectName);
    await runCommand('pnpm', ['run', 'cf-typegen'], projectName);
  } catch (error) {
    console.log(pc.red(`生成类型提示失败，请手动运行 vvcli types 生成`));
  }
}

async function createVueProject(projectName: string, jsonData: object) {
  const vueJSCode = `
import { Hono } from "hono";
import { EnvBinding } from "./binding";
import { init, useDebugEndpoint } from "@dicarne/vvorker-sdk";
import { env } from "cloudflare:workers";
init(env)

const app = new Hono<{ Bindings: EnvBinding }>();
useDebugEndpoint(app)
app.onError(async (err, c) => {
  console.error(err)
  return c.json({
    code: 500,
    msg: err.message,
    data: null
  }, 500)
})

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

  await runCommand('pnpm', ['create', "cloudflare@latest", projectName, "--framework=vue", "--git", "--no-deploy", "--lang=ts"]);


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

  const wranglerJsonPath = path.join(projectName, 'wrangler.json');
  const wranglerJson = json5.parse(await fs.readFile(wranglerJsonPath, 'utf-8'));
  wranglerJson.compatibility_flags = ["nodejs_compat"];
  wranglerJson.durable_objects = undefined;
  wranglerJson.migrations = undefined;
  await fs.writeJson(wranglerJsonPath, wranglerJson, { spaces: 2 });

  const packageJsonPath = path.join(projectName, 'package.json');
  const packageJson = json5.parse(await fs.readFile(packageJsonPath, 'utf-8'));
  packageJson.scripts.deploy = undefined;
  await fs.writeJson(packageJsonPath, packageJson, { spaces: 2 });

  console.log(pc.green(`项目 ${projectName} 初始化完成`));

  try {
    await runCommand('pnpm', ['install'], projectName);
    await runCommand('pnpm', ['install', "@types/node", "-D"], projectName);
    await runCommand('pnpm', ['install', "hono", "@dicarne/vvorker-sdk", "@hono/zod-validator", "zod"], projectName);
  } catch (error) {
    console.log(pc.red(`安装依赖失败，请手动安装`));
  }
  try {
    await runCommand('vvcli', ['types'], projectName);
    await runCommand('pnpm', ['run', 'cf-typegen'], projectName);
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
      "services": [],
      "vars": {},
      "ai": [],
      "oss": [],
      "pgsql": [],
      "mysql": [],
      "kv": [],
      "assets": [],
      "compatibility_flags": [
        "nodejs_compat"
      ],
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
