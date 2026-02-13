import { Command } from 'commander';
import { config, ensureEnv, saveConfig, setUrl } from '../utils/config';
import pc from "picocolors"

export const configCommand = new Command('config')
  .description('配置vvorker平台的参数');

configCommand.command('set <key> <value>')
  .description('设置配置项')
  .action((key, value) => {
    if (key === 'url') {
      setUrl(value);
    } else if (key === 'env') {
      if (!config.env || !config.env[value]) {
        console.error(pc.red(`环境 "${value}" 不存在`));
        console.error(pc.gray('请先使用 ') + pc.cyan('vvcli create') + pc.gray(' 创建环境'));
        return;
      }
      config.current_env = value;
    } else if (key === 'token') {
      if (!config.current_env) {
        console.error(pc.red('当前没有选择环境'));
        console.error(pc.gray('请先使用 ') + pc.cyan('vvcli create') + pc.gray(' 创建环境'));
        return;
      }
      ensureEnv(config.current_env);
      config.env[config.current_env].token = value;
    }
    saveConfig();
    console.log('配置成功');
  });

configCommand.command('show')
  .description('显示当前配置')
  .action(() => {
    if (!config.current_env) {
      console.log(pc.red('当前没有选择环境'));
      console.log(pc.gray('请先使用 ') + pc.cyan('vvcli create') + pc.gray(' 创建环境'));
      return;
    }
    console.log(pc.green('Env：') + "\t" + config.current_env);
    const envConfig = config.env[config.current_env];
    if (!envConfig) {
      console.log(pc.red('环境配置不存在'));
      return;
    }
    const token = envConfig.token;
    if (!token) {
      console.log(pc.red('Token：') + "\t" + "未配置");
    } else {
      const mid = Math.floor(token.length / 2);
      const showToken = token.slice(0, mid) + "*".repeat(8) + token.slice(mid + 8);
      console.log(pc.green('Token：') + "\t" + showToken);
    }
    console.log(pc.green('URL：') + "\t" + envConfig.url);

  });