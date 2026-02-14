#!/usr/bin/env node
import { Command } from "commander";
import { initCommand } from "./commands/init";
import { deployCommand } from "./commands/deploy";
import { configCommand } from "./commands/config";
import { typesCommand } from "./commands/types";
import { envCommand } from "./commands/env";
import { devCommand } from "./commands/dev";
import { versionCommand } from "./commands/version";
import { createCommand } from "./commands/create";
import { upgradeCommand } from "./commands/upgrade";
import { logsCommand } from "./commands/logs";
import { installCommand } from "./commands/install";

// 全局变量，用于存储工作目录
let globalWorkingDir: string | null = null;

/**
 * 设置工作目录（供外部使用）
 */
export function setWorkingDir(dir: string) {
  globalWorkingDir = dir;
}

/**
 * 获取工作目录
 */
export function getWorkingDir(): string {
  if (globalWorkingDir) {
    return globalWorkingDir;
  }
  return process.cwd();
}

const program = new Command();

// 全局选项：指定工作目录
program.option('-d, --directory <dir>', '指定工作目录');

// 存储工作目录选项的hook
program.hook('preAction', (thisCommand) => {
  const opts = thisCommand.opts();
  if (opts.directory) {
    globalWorkingDir = opts.directory;
  }
});

program.addCommand(initCommand);
program.addCommand(deployCommand);
program.addCommand(configCommand);
program.addCommand(typesCommand);
program.addCommand(envCommand);
program.addCommand(devCommand);
program.addCommand(versionCommand);
program.addCommand(createCommand);
program.addCommand(upgradeCommand);
program.addCommand(logsCommand);
program.addCommand(installCommand);

program.parse(process.argv);
