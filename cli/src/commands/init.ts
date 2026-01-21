import { Command } from 'commander';
import inquirer from 'inquirer';
import * as fs from 'fs-extra';
import * as path from 'path';
import json5 from 'json5';
import { runCommand } from '../utils/system';
import pc from "picocolors"
import { config, getEnv } from '../utils/config';
import { execSync } from 'child_process';

async function createWorkerProject(projectName: string, jsonData: object, gitRepo: string) {
  try {
    execSync(`git clone ${gitRepo} ${projectName}`, { stdio: 'inherit' });
    // 删除 .git 目录
    await fs.remove(path.join(projectName, '.git'));
  } catch (error) {
    console.log(pc.red(`克隆 Git 仓库失败: ${error}`));
    throw error;
  }


  const env = getEnv();
  const jsonFilePath = `vvorker.${env}.json`;
  await fs.writeJson(path.join(projectName, jsonFilePath), jsonData, { spaces: 2 });

  const packageJsonPath = path.join(projectName, 'package.json');
  const packageJson = json5.parse(await fs.readFile(packageJsonPath, 'utf-8'));
  packageJson.name = projectName;
  await fs.writeJson(packageJsonPath, packageJson, { spaces: 2 });

  console.log(pc.green(`项目 ${projectName} 初始化完成`));

  try {
    await runCommand('pnpm', ['install'], projectName);
  } catch (error) {
    console.log(pc.red(`安装依赖失败，请手动安装`));
  }
  try {
    await runCommand('git', ['init'], projectName);
    await runCommand('vvcli', ['types'], projectName);
    await runCommand('pnpm', ['run', 'cf-typegen'], projectName);
  } catch (error) {
    console.log(pc.red(`生成类型提示失败，请手动运行 vvcli types 生成`));
  } finally {
    await runCommand('git', ['add', "*"], projectName);
    await runCommand('git', ['commit', "-m", "\"init: Create with vvcli.\""], projectName);
  }
}

async function createVueProject(projectName: string, jsonData: object, gitRepo: string) {
  try {
    execSync(`git clone ${gitRepo} ${projectName}`, { stdio: 'inherit' });
    // 删除 .git 目录
    await fs.remove(path.join(projectName, '.git'));
  } catch (error) {
    console.log(pc.red(`克隆 Git 仓库失败: ${error}`));
    throw error;
  }


  (jsonData as any)["assets"] = [
    {
      "directory": "./dist/client",
      "binding": "ASSETS"
    }
  ];
  const env = getEnv();
  const jsonFilePath = `vvorker.${env}.json`;
  await fs.writeJson(path.join(projectName, jsonFilePath), jsonData, { spaces: 2 });


  const packageJsonPath = path.join(projectName, 'package.json');
  const packageJson = json5.parse(await fs.readFile(packageJsonPath, 'utf-8'));
  packageJson.name = projectName;
  await fs.writeJson(packageJsonPath, packageJson, { spaces: 2 });

  console.log(pc.green(`项目 ${projectName} 初始化完成`));

  try {
    await runCommand('pnpm', ['install'], projectName);
  } catch (error) {
    console.log(pc.red(`安装依赖失败，请手动安装`));
  }
  try {
    await runCommand('git', ['init'], projectName);
    await runCommand('vvcli', ['types'], projectName);
    await runCommand('pnpm', ['run', 'cf-typegen'], projectName);
  } catch (error) {
    console.log(pc.red(`生成类型提示失败，请手动运行 vvcli types 生成`));
  } finally {
    await runCommand('git', ['add', "*"], projectName);
    await runCommand('git', ['commit', "-m", "\"init: Create with vvcli.\""], projectName);
  }
}

export const initCommand = new Command('init')
  .command('init <projectName>')
  .description('初始化VVorker项目')
  .option('--git-repo <url>', '使用自定义 Git 仓库作为模板')
  .action(async (projectName, options) => {
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
      "version": "0.0.0",
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
        await createWorkerProject(projectName, jsonData, "https://git.cloud.zhishudali.ink/template/vv-template-worker.git");
        break;
      }
      case "vue": {
        await createVueProject(projectName, jsonData, "https://git.cloud.zhishudali.ink/template/vv-template-vue.git");
        break;
      }
    }
  });
