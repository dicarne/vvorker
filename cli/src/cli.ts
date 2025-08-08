#!/usr/bin/env node
import { Command } from 'commander';
import { initCommand } from './commands/init';
import { deployCommand } from './commands/deploy';
import { configCommand } from './commands/config';
import { typesCommand } from './commands/types';
import { envCommand } from './commands/env';
import { devCommand } from './commands/dev';

const program = new Command();

program.addCommand(initCommand);
program.addCommand(deployCommand);
program.addCommand(configCommand);
program.addCommand(typesCommand);
program.addCommand(envCommand);
program.addCommand(devCommand);

program.parse(process.argv);