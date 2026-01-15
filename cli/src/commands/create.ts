import { Command } from 'commander';
import inquirer from 'inquirer';
import { config, ensureEnv, saveConfig } from '../utils/config';
import pc from "picocolors";

export const createCommand = new Command('create')
  .description('创建新的环境')
  .action(async () => {
    const answers = await inquirer.prompt([
      {
        type: 'input',
        name: 'envName',
        message: '请输入环境名称：',
        validate: (input: string) => {
          if (!input || input.trim() === '') {
            return '环境名称不能为空';
          }
          if (config.env[input]) {
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
        message: '请输入 Token：',
        validate: (input: string) => {
          if (!input || input.trim() === '') {
            return 'Token 不能为空';
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
    saveConfig();

    console.log(pc.green('环境创建成功！'));
    console.log(pc.green('环境名称：') + envName);
    console.log(pc.green('URL：') + url);
    const mid = Math.floor(token.length / 2);
    const showToken = token.slice(0, mid) + "*".repeat(8) + token.slice(mid + 8);
    console.log(pc.green('Token：') + showToken);
  });
