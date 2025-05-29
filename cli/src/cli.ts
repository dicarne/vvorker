#!/usr/bin/env node
import { Command } from 'commander';
import inquirer from 'inquirer';
import axios from 'axios';
import * as fs from 'fs-extra';
import * as path from 'path'
import { spawn } from 'node:child_process';

const program = new Command();

// 初始化项目命令
program
  .command('init <projectName>')
  .description('初始化VVorker项目')
  .action(async (projectName) => {
    // 交互式输入uid
    const { uid } = await inquirer.prompt([{
      type: 'input',
      name: 'uid',
      message: '请输入vvorker平台worker的uid:',
    }])

    // 创建项目目录
    await fs.ensureDir(projectName);

    // 定义要写入的 JSON 数据
    const jsonData = {
      "name": "worker",
      uid,
      "version": "1.0.0",
      "extensions": [],
      "services": [],
      "vars": {},
      "ai": [],
      "oss": [],
      "pgsql": [],
      "kv": [],
    };
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
  });

// 部署命令
program
  .command('deploy')
  .description('部署到vvorker')
  .action(async () => {
    if (!config.url) {
      console.error('请先配置vvorker平台的url');
      return;
    }
    // 读取当前目录下的 vvorker.json 文件
    const vvorkerJson = await fs.readJson('vvorker.json');
    let serviceName = vvorkerJson.name;
    if (!serviceName) {
      console.error('服务名称不能为空');
      return;
    }

    const uid = vvorkerJson.uid;
    if (!uid) {
      console.error('uid不能为空');
      return;
    }

    const token = config.token;
    if (!token) {
      console.error('token不能为空');
      console.error('请先配置token');
      console.error('执行命令：vvcli config set token <token>');
      return;
    }

    // 运行 pnpm run build 命令
    // 如果存在pnpm-lock.yaml，则使用pnpm
    if (fs.existsSync('pnpm-lock.yaml')) {
      // 使用 spawn 执行 pnpm run build 并将输出重定向到主进程
      await new Promise((resolve, reject) => {
        const child = spawn('pnpm', ['run', 'build'], { stdio: 'inherit', shell: true });
        child.on('close', (code) => {
          if (code === 0) {
            resolve(code);
          } else {
            reject(new Error(`pnpm run build exited with code ${code}`));
          }
        });
        child.on('error', (error) => {
          reject(error);
        });
      });
    } else {
      // 使用 spawn 执行 npm run build 并将输出重定向到主进程
      await new Promise((resolve, reject) => {
        const child = spawn('npm', ['run', 'build'], { stdio: 'inherit', shell: true });
        child.on('close', (code) => {
          if (code === 0) {
            resolve(code);
          } else {
            reject(new Error(`npm run build exited with code ${code}`));
          }
        });
        child.on('error', (error) => {
          reject(error);
        });
      });
    }

    // 往dist目录下写入vvorker.json
    const distFilePath = `${process.cwd()}/dist/vvorker.json`;
    await fs.writeJson(distFilePath, vvorkerJson, { spaces: 2 });

    // 读取js并转化成base64
    const jsFilePath = `${process.cwd()}/dist/index.js`;
    const jsContent = await fs.readFile(jsFilePath, 'utf-8');
    const jsBase64 = Buffer.from(jsContent).toString('base64');

    // url join
    if (config.url.endsWith('/')) {
      config.url = config.url.slice(0, -1);
    }


    let resp = await axios.get(`${config.url}/api/worker/${uid}`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      }
    })
    let prev = resp.data.data[0]

    //读取 /dist/vvorker.json
    const distVvorkerJson = await fs.readJson(`${process.cwd()}/dist/vvorker.json`);
    prev.Code = jsBase64;
    prev.Template = JSON.stringify(distVvorkerJson);
    prev.HostName = undefined;
    prev.ExternalPath = undefined;
    prev.TunnelID = undefined;

    resp = await axios.post(`${config.url}/api/worker/${uid}`, prev, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      data: JSON.stringify(prev)
    })
    if (resp.data.code === 0) {
      console.log('部署成功');
    }

  });

// 在用户目录下创建 .vvcli 目录
const userHome = require('os').homedir();
const vvcliDir = `${userHome}/.vvcli`;
fs.ensureDirSync(vvcliDir)

// 检查是否存在配置文件，如果不存在则创建
const configFilePath = `${vvcliDir}/config.json`;
if (!fs.existsSync(configFilePath)) {
  fs.writeFileSync(configFilePath, JSON.stringify({}));
}

// 读取配置文件
const config = JSON.parse(fs.readFileSync(configFilePath, 'utf-8'));


// 配置命令：通用配置命令，支持配置不同的键值对
const cmd_config = program
  .command('config')
  .description('配置vvorker平台的参数')

cmd_config.command('set <key> <value>')
  .description('设置配置项')
  .action((key, value) => {
    // 根据传入的 key 设置对应的配置项
    config[key] = value;
    fs.writeFileSync(configFilePath, JSON.stringify(config));
    console.log('配置成功');
  })

program.parse(process.argv);