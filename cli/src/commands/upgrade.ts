import { Command } from 'commander';
import * as fs from 'fs-extra';
import { loadVVorkerConfig } from '../utils/vvorker-config';
import { apiClient } from '../utils/api';
import { runCommand } from '../utils/system';

export const upgradeCommand = new Command('upgrade')
    .description("用于升级cli")
    .action(async () => {
        await runCommand('pnpm', ['i', '@dicarne/vvcli', '-g']);

    });
