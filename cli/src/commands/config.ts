import { Command } from 'commander';
import { config, ensureEnv, saveConfig, setUrl } from '../utils/config';

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
