import { Command } from 'commander';
import { config, ensureEnv, saveConfig, setUrl } from '../utils/config';
import pc from "picocolors"

export const configCommand = new Command('config')
  .description('配置vvorker平台的参数');

configCommand.command('set <key> <value>')
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
    saveConfig();
    console.log('配置成功');
  });

configCommand.command('show')
  .description('显示当前配置')
  .action(() => {
    console.log(pc.green('Env：') + "\t" + config.current_env);
    const token = config.env[config.current_env ?? "default"].token;
    if (!token) {
      console.log(pc.red('Token：') + "\t" + "未配置");
    } else {
      const mid = Math.floor(token.length / 2);
      const showToken = token.slice(0, mid) + "*".repeat(8) + token.slice(mid + 8);
      console.log(pc.green('Token：') + "\t" + showToken);
    }
    console.log(pc.green('URL：') + "\t" + config.env[config.current_env ?? "default"].url);

  });