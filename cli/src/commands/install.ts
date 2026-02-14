import { Command } from 'commander';
import { runCommand } from '../utils/system';
import { withWorkingDir } from '../utils/working-dir';
import { checkForUpdate } from '../utils/update-check';

export const installCommand = new Command('install')
    .description("安装依赖并更新vvorker-sdk到最新版本")
    .action(async () => {
        await checkForUpdate();
        
        await withWorkingDir(async () => {
            console.log('正在运行 pnpm install...');
            await runCommand('pnpm', ['install']);
            
            console.log('正在更新 @dicarne/vvorker-sdk 到最新版本...');
            await runCommand('pnpm', ['update', '@dicarne/vvorker-sdk', '--latest']);
            
            console.log('安装完成！');
        });
    });
