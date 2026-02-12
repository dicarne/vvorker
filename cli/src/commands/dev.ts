import { Command } from 'commander';
import { runCommand } from '../utils/system';
import { withWorkingDir } from '../utils/working-dir';

export const devCommand = new Command('dev')
    .description("用于开发")
    .action(async () => {
        await withWorkingDir(async () => {
            runCommand('pnpm', ['run', 'dev']);
        });
    });
