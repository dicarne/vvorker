import { Command } from 'commander';
import { apiClient } from '../utils/api';
import { loadVVorkerConfig } from '../utils/vvorker-config';
import pc from "picocolors";
import axios from 'axios';

export const logsCommand = new Command('logs')
  .description('获取VVorker worker日志')
  .option('-f, --follow', '实时显示日志')
  .option('-p,--page <number>', '页码，从1开始', '1')
  .option('--page-size <number>', '每页记录数', '50')
  .action(async (options) => {
    try {
      const vvorkerJson = loadVVorkerConfig();
      const uid = vvorkerJson.project?.uid ?? vvorkerJson.uid;
      
      if (!uid) {
        console.error(pc.red('未找到 worker uid'));
        return;
      }

      if (options.follow) {
        await streamLogs(uid);
      } else {
        await fetchLogs(uid, parseInt(options.page), parseInt(options.pageSize));
      }
    } catch (error) {
      console.error(pc.red(`获取日志失败: ${error}`));
    }
  });

async function fetchLogs(uid: string, page: number, pageSize: number) {
  try {
    const response = await apiClient.post(`/api/worker/logs/${uid}`, {
      page: page,
      page_size: pageSize
    });

    if (response.data.code !== 0) {
      console.error(pc.red(`获取日志失败: ${response.data.message || '未知错误'}`));
      return;
    }

    const { total, logs } = response.data.data;

    if (!logs || logs.length === 0) {
      console.log(pc.yellow('暂无日志'));
      return;
    }

    // 显示日志总数
    console.log(pc.gray(`总计 ${total} 条日志，第 ${page} 页`));
    console.log(pc.gray('='.repeat(60)));

    // 显示日志，按时间正序排列（最新在底部）
    for (const log of logs.slice().reverse()) {
      const time = new Date(log.time).toLocaleString('zh-CN');
      const typeColor = log.type === 'error' ? pc.red : log.type === 'warn' ? pc.yellow : pc.white;
      console.log(typeColor(`[${time}] [${log.type.toUpperCase()}]`));
      console.log(log.output);
    }
  } catch (error) {
    if (axios.isAxiosError(error)) {
      if (error.response?.status === 404) {
        console.error(pc.red('日志接口不存在'));
      } else if (error.response?.status === 401) {
        console.error(pc.red('认证失败，请检查 token'));
      } else {
        console.error(pc.red(`请求失败: ${error.message}`));
      }
    } else {
      console.error(pc.red(`获取日志失败: ${error}`));
    }
  }
}

async function streamLogs(uid: string) {
  console.log(pc.gray('实时日志模式启动，按 Ctrl+C 停止...'));
  console.log(pc.gray('='.repeat(60)));

  let lastTime: Date | null = null;
  let isRunning = true;

  // 处理 Ctrl+C
  process.on('SIGINT', () => {
    console.log(pc.gray('\n停止实时日志...'));
    isRunning = false;
  });

  while (isRunning) {
    try {
      const response = await apiClient.post(`/api/worker/logs/${uid}`, {
        page: 1,
        page_size: 50
      });

      if (response.data.code !== 0) {
        console.error(pc.red(`获取日志失败: ${response.data.message || '未知错误'}`));
        await sleep(2000);
        continue;
      }

      const { logs } = response.data.data;

      if (!logs || logs.length === 0) {
        await sleep(2000);
        continue;
      }

      // 按时间排序
      const sortedLogs = logs.sort((a: any, b: any) =>
        new Date(a.time).getTime() - new Date(b.time).getTime()
      );

      // 过滤出新的日志
      const newLogs = lastTime
        ? sortedLogs.filter((log: any) => new Date(log.time) > lastTime!)
        : sortedLogs;

      if (newLogs.length > 0) {
        for (const log of newLogs) {
          const time = new Date(log.time).toLocaleString('zh-CN');
          const typeColor = log.type === 'error' ? pc.red : log.type === 'warn' ? pc.yellow : pc.white;
          console.log(typeColor(`[${time}] [${log.type.toUpperCase()}]`));
          console.log(log.output);
        }
        console.log(pc.gray('-'.repeat(60)));
        
        // 更新最后一条日志时间
        lastTime = new Date(sortedLogs[sortedLogs.length - 1].time);
      }

      // 等待 2 秒后再次查询
      await sleep(2000);
    } catch (error) {
      if (isRunning) {
        console.error(pc.red(`获取日志失败: ${error}`));
      }
      await sleep(2000);
    }
  }
}

function sleep(ms: number): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms));
}
