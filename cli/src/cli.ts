#!/usr/bin/env node
import { Command } from 'commander';
import inquirer from 'inquirer';
import axios from 'axios';
import * as fs from 'fs-extra';
import * as path from 'path'
import { spawn } from 'node:child_process';

import json5 from 'json5'

interface EnvConfig {
  url?: string;
  token?: string;
}

interface Config {
  current_env: string;
  env: { [key: string]: EnvConfig };
}

function getToken() {
  let env = config.current_env ?? "default";
  return config.env[env]?.token;
}

function getUrl() {
  let env = config.current_env ?? "default";
  let url = config.env[env]?.url;
  if (url?.endsWith('/')) {
    url = url.slice(0, -1)
  }
  return url
}

function ensureEnv(env: string) {
  if (!config.env) {
    config.env = {

    }
  }
  if (!config.env[env]) {
    config.env[env] = {
    }
  }
}

function setUrl(url: string) {
  let env = config.current_env ?? "default";
  ensureEnv(env);
  config.env[env].url = url;
}

function loadVVorkerConfig() {
  let current_env = config.current_env
  if (fs.existsSync(`${process.cwd()}/vvorker.${current_env}.json`)) {
    let vvorkerJson = json5.parse(fs.readFileSync(`${process.cwd()}/vvorker.${current_env}.json`, 'utf-8'))
    return vvorkerJson
  } else {
    return json5.parse(fs.readFileSync(`${process.cwd()}/vvorker.json`, 'utf-8'))
  }
}

function saveVVorkerConfig(vvorkerJson: any) {
  let current_env = config.current_env
  if (fs.existsSync(`${process.cwd()}/vvorker.${current_env}.json`)) {
    fs.writeFileSync(`${process.cwd()}/vvorker.${current_env}.json`, json5.stringify(vvorkerJson, null, 2))
  } else {
    fs.writeFileSync(`${process.cwd()}/vvorker.json`, json5.stringify(vvorkerJson, null, 2))
  }
}

const program = new Command();

// 初始化项目命令
program
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
        const jsFilePath = `${projectName}/src/index.ts`;
        await fs.writeFile(jsFilePath, vueJSCode);
        await new Promise((resolve, reject) => {
          // npm create cloudflare@latest -- my-vue-app --framework=vue
          const child = spawn('pnpm', ['create', "cloudflare@latest", projectName, "--framework=vue"], { stdio: 'inherit', shell: true });
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
      }

      const jsonFilePath = `vvorker.json`;
      (jsonData as any)["assets"] = [
        {
          "directory": "./dist",
          "binding": "ASSETS"
        }
      ],
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

// 部署命令
program
  .command('deploy')
  .description('部署到vvorker')
  .action(async () => {
    if (!getUrl()) {
      console.error('请先配置vvorker平台的url');
      return;
    }
    // 读取当前目录下的 vvorker.json 文件
    const vvorkerJson = loadVVorkerConfig();
    const packageJson = await fs.readJson('package.json');
    let serviceName = vvorkerJson.name;
    if (!serviceName) {
      console.error('服务名称不能为空');
      return;
    }

    const uid = vvorkerJson.project?.uid ?? vvorkerJson.uid;
    if (!uid) {
      console.error('uid不能为空');
      return;
    }

    const token = getToken();
    if (!token) {
      console.error('token不能为空');
      console.error('请先配置token');
      console.error('执行命令：vvcli config set token <token>');
      return;
    }
    // url join
    if (getUrl()!.endsWith('/')) {
      setUrl(getUrl()!.slice(0, -1));
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

    let jsFilePath = "";
    if (vvorkerJson.assets && vvorkerJson.assets.length > 0) {


      let wwwAssetsPath = path.join(process.cwd(), vvorkerJson.assets[0].directory)
      // walk wwwAssetsPath，调用接口上传每一个文件
      const walk = async (dir: string) => {
        const files = fs.readdirSync(dir);
        for (let file of files) {
          const filePath = path.join(dir, file);
          const stat = fs.statSync(filePath);
          if (stat.isDirectory()) {
            walk(filePath);
          } else {
            const fileContent = fs.readFileSync(filePath);
            // const fileBase64 = Buffer.from(fileContent).toString('base64');
            const fileUrl = filePath.replace(wwwAssetsPath, '').replace(/\\/g, '/');
            console.log(fileUrl);

            // 创建 FormData 对象
            const formData = new FormData();
            // 将文件数据添加到 FormData 中
            formData.append('file', new Blob([fileContent]), fileUrl);

            let up1 = await axios.post(`${getUrl()}/api/file/upload`, formData, {
              headers: {
                'Authorization': `Bearer ${token}`,
                // 设置 Content-Type 为 multipart/form-data
                'Content-Type': 'multipart/form-data'
              },
            })

            if (up1.status !== 200) {
              throw new Error(`上传失败：${fileUrl} ${up1.status} ${up1.statusText}`);
            }

            let fileuid = up1.data.data.fileId;

            try {
              let resp = await axios.post(`${getUrl()}/api/ext/assets/create-assets`, {
                uid: fileuid,
                "worker_uid": uid,
                "path": fileUrl,
              }, {
                headers: {
                  'Authorization': `Bearer ${token}`,
                },
              })

              if (resp.status != 200) {
                console.log(`上传失败：${fileUrl} ${resp.status} ${resp.statusText}`);
                throw new Error(`上传失败：${fileUrl}`);
              }

            } catch (error) {
              console.log(`上传失败：${fileUrl} ${error}`);
            }


          }
        }
      }
      await walk(wwwAssetsPath);
    }

    if (vvorkerJson.project.type === "vue") {
      jsFilePath = `${process.cwd()}/dist/${packageJson.name.replaceAll("-", "_")}/index.js`;
    } else {
      jsFilePath = `${process.cwd()}/dist/index.js`;
    }

    // 往dist目录下写入vvorker.json
    const distFilePath = `${process.cwd()}/dist/vvorker.json`;
    await fs.writeJson(distFilePath, vvorkerJson, { spaces: 2 });

    // 读取js并转化成base64
    const jsContent = await fs.readFile(jsFilePath, 'utf-8');
    const jsBase64 = Buffer.from(jsContent).toString('base64');




    let resp = await axios.post(`${getUrl()}/api/worker/v2/get-worker`, {
      uid: uid,
    }, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })
    let prev = resp.data.data[0]

    //读取 /dist/vvorker.json
    const distVvorkerJson = await fs.readJson(`${process.cwd()}/dist/vvorker.json`);
    if (distVvorkerJson.pgsql && distVvorkerJson.pgsql.length > 0) {
      for (let pgsql of distVvorkerJson.pgsql) {
        let rid = pgsql.resource_id;
        if (pgsql.migrate) {
          let migrateFiles = fs.readdirSync(pgsql.migrate);
          let allFile: {
            name: string,
            content: string
          }[] = []
          for (let migrateFile of migrateFiles) {
            if (!migrateFile.endsWith('.sql')) {
              continue;
            }
            let migrateFilePath = path.join(pgsql.migrate, migrateFile);
            let migrateFileContent = fs.readFileSync(migrateFilePath, 'utf-8');
            allFile.push({
              name: migrateFile,
              content: migrateFileContent
            })
          }

          let resp = await axios.post(`${getUrl()}/api/ext/pgsql/migrate`, {
            resource_id: rid,
            files: allFile,
          }, {
            headers: {
              'Authorization': `Bearer ${token}`,
            },
          })
          if (resp.data.code !== 0) {
            console.log(`迁移失败：${pgsql.migrate}`);
            throw new Error(`迁移失败：${pgsql.migrate}`);
          }
          console.log(`迁移成功：${pgsql.migrate}`);
        }
      }
    }

    //
    prev.Code = jsBase64;
    prev.Template = JSON.stringify(distVvorkerJson);
    prev.HostName = undefined;
    prev.ExternalPath = undefined;
    prev.TunnelID = undefined;

    resp = await axios.post(`${getUrl()}/api/worker/v2/update-worker`, prev, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
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
const config = JSON.parse(fs.readFileSync(configFilePath, 'utf-8')) as Config;


// 配置命令：通用配置命令，支持配置不同的键值对
const cmd_config = program
  .command('config')
  .description('配置vvorker平台的参数')

cmd_config.command('set <key> <value>')
  .description('设置配置项')
  .action((key, value) => {
    // 根据传入的 key 设置对应的配置项
    if (key === 'url') {
      setUrl(value);
    } else if (key === 'env') {
      config.current_env = value;
    } else if (key === 'token') {
      ensureEnv(config.current_env ?? "default");
      config.env[config.current_env ?? "default"].token = value;
    }
    fs.writeFileSync(configFilePath, JSON.stringify(config));
    console.log('配置成功');
  })

program.command('types')
  .description("用于自动生成配置文件对应的TypeScript类型")
  .action(async () => {
    const vvv = loadVVorkerConfig();
    let resp = await axios.post(`${getUrl()}/api/ext/types`, vvv, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getToken()}`,
      }
    })
    if (resp.status !== 200) {
      throw new Error(`获取类型失败：${resp.status} ${resp.statusText}`);
    }
    let w = resp.data.data.types;
    if (vvv.project.type === "vue") {
      fs.writeFileSync(`${process.cwd()}/server/binding.ts`, w);
    } else {
      fs.writeFileSync(`${process.cwd()}/src/binding.ts`, w);
    }
    fs.writeFileSync(`${process.cwd()}/vvorker-schema.json`, resp.data.data.schema)
    console.log('类型生成成功');
  })

program.command("env")
  .description("切换当前可用环境")
  .action(async (env) => {
    console.log(`当前环境：${config.current_env}`);
    let all_env = []
    for (const key in config.env) {
      if (Object.prototype.hasOwnProperty.call(config.env, key)) {
        const element = config.env[key];
        all_env.push(key)
      }
    }

    const { envname } = await inquirer.prompt([{
      type: 'list',
      name: 'envname',
      message: '切换环境',
      choices: all_env.map(s => {
        return {
          name: s,
          value: s,
        }
      }),
      default: config.current_env
    },])
    config.current_env = envname;
    fs.writeFileSync(configFilePath, JSON.stringify(config));
  })

program.parse(process.argv);