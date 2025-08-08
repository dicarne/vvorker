import { Command } from 'commander';
import { runCommand } from '../utils/system';

export const devCommand = new Command('dev')
    .description("用于开发")
    .action(async () => {
        runCommand('pnpm', ['run', 'dev']);
    });
