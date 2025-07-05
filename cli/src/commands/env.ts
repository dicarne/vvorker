import { Command } from 'commander';
import inquirer from 'inquirer';
import { config, saveConfig } from '../utils/config';

export const envCommand = new Command('env')
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
    saveConfig();
  });
