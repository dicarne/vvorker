import axios from 'axios';
import { config, saveConfig } from './config';
import pc from 'picocolors';
import { runCommand } from './system';

const PACKAGE_NAME = '@dicarne/vvcli';

/**
 * 获取今天的日期字符串 (YYYY-MM-DD)
 */
function getTodayString(): string {
  const now = new Date();
  return now.toISOString().split('T')[0];
}

/**
 * 检查是否需要检查更新
 */
function shouldCheckUpdate(): boolean {
  // 如果禁用更新检查，则跳过
  if (config.disable_update === true) {
    return false;
  }

  // 如果没有上次检查时间，需要检查
  if (!config.last_update_time) {
    return true;
  }

  // 如果上次检查不是今天，需要检查
  const today = getTodayString();
  return config.last_update_time !== today;
}

/**
 * 更新 CLI 到最新版本
 */
export async function upgradeCLI(): Promise<void> {
  console.log(pc.yellow('正在更新 vvcli...'));
  try {
    await runCommand('pnpm', ['update', PACKAGE_NAME, '-g']);
    console.log(pc.green('✓ vvcli 更新完成！'));
  } catch (error) {
    console.log(pc.red('✗ vvcli 更新失败，请手动运行: pnpm update @dicarne/vvcli -g'));
  }
}

/**
 * 检查 CLI 更新
 * 每天只检查一次，如果发现新版本则提示用户
 */
export async function checkForUpdate(): Promise<void> {
  if (!shouldCheckUpdate()) {
    return;
  }

  // 记录今天的检查时间
  config.last_update_time = getTodayString();
  saveConfig();
  await runCommand('pnpm', ['update', '@dicarne/vvcli', '-g']);
}
