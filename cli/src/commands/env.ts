import { Command } from 'commander';
import inquirer from 'inquirer';
import { config, saveConfig } from '../utils/config';
import pc from "picocolors";

export const envCommand = new Command('env')
  .description("切换当前可用环境")
  .action(async (env) => {
    // 检查是否有环境
    const envKeys = Object.keys(config.env || {});
    if (envKeys.length === 0) {
      console.log(pc.yellow('当前没有环境'));
      console.log(pc.gray('请先使用 ') + pc.cyan('vvcli create') + pc.gray(' 创建环境'));
      return;
    }

    console.log(`当前环境：${config.current_env || pc.gray('未选择')}`);
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
    saveConfig();
  });
