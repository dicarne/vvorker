import * as path from "path";
import { getWorkingDir as getCliWorkingDir } from '../cli';

/**
 * 获取当前工作目录的绝对路径
 * 从cli.ts的全局变量中获取
 */
export function getWorkingDir(): string {
  return getCliWorkingDir();
}

/**
 * 获取相对于工作目录的绝对路径
 * @param relativePath 相对路径
 */
export function getAbsolutePath(relativePath: string): string {
  return path.resolve(getWorkingDir(), relativePath);
}

/**
 * 切换到工作目录执行函数
 * @param fn 要执行的函数
 */
export async function withWorkingDir<T>(fn: () => Promise<T>): Promise<T> {
  const originalCwd = process.cwd();
  const workingDir = getWorkingDir();
  try {
    process.chdir(workingDir);
    return await fn();
  } finally {
    // 恢复原始工作目录
    process.chdir(originalCwd);
  }
}

/**
 * 同步版本：切换到工作目录执行函数
 * @param fn 要执行的函数
 */
export function withWorkingDirSync<T>(fn: () => T): T {
  const originalCwd = process.cwd();
  const workingDir = getWorkingDir();

  try {
    process.chdir(workingDir);
    return fn();
  } finally {
    // 恢复原始工作目录
    process.chdir(originalCwd);
  }
}
