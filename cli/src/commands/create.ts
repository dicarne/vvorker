import { Command } from 'commander';
import inquirer from 'inquirer';
import { config, ensureEnv, saveConfig } from '../utils/config';
import pc from "picocolors";

export const createCommand = new Command('create')
  .description('创建新的环境')
  .action(async () => {
    console.log(pc.gray('请根据文档说明创建新的环境：\n') + pc.cyan('https://vvorker-docs.vvorker.zhishudali.ink/config/env.html'));
    console.log();
    const answers = await inquirer.prompt([
      {
        type: 'input',
        name: 'envName',
        message: '请输入环境名称：',
        validate: (input: string) => {
          if (!input || input.trim() === '') {
            return '环境名称不能为空';
          }
          if (config.env && config.env[input]) {
            return '该环境已存在，请使用其他名称';
          }
          return true;
        }
      },
      {
        type: 'input',
        name: 'url',
        message: '请输入环境 URL：',
        validate: (input: string) => {
          if (!input || input.trim() === '') {
            return 'URL 不能为空';
          }
          try {
            new URL(input);
            return true;
          } catch {
            return '请输入有效的 URL';
          }
        }
      },
      {
        type: 'input',
        name: 'token',
        message: '请输入 API 密钥：',
        validate: (input: string) => {
          if (!input || input.trim() === '') {
            return 'API 密钥不能为空';
          }
          return true;
        }
      }
    ]);

    const { envName, url, token } = answers;

    ensureEnv(envName);
    config.env[envName] = {
      url,
      token
    };

    // 如果这是第一个环境，自动设置为当前环境
    if (!config.current_env) {
      config.current_env = envName;
    }

    saveConfig();

    console.log(pc.green('环境创建成功！'));
    console.log(pc.green('环境名称：') + envName);
    console.log(pc.green('URL：') + url);
    const mid = Math.floor(token.length / 2);
    const showToken = token.slice(0, mid) + "*".repeat(8) + token.slice(mid + 8);
    console.log(pc.green('API 密钥：') + showToken);
    if (!config.current_env || config.current_env === envName) {
      console.log(pc.gray('已自动设置为当前环境'));
    }
  });
