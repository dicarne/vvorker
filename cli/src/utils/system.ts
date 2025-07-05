import { spawn } from 'node:child_process';
import path from 'node:path';

export function runCommand(command: string, args: string[], cwd?: string): Promise<number> {
  if (cwd) {
    cwd = path.join(process.cwd(), cwd)
  } else {
    cwd = process.cwd()
  }
  return new Promise((resolve, reject) => {
    const child = spawn(command, args, { stdio: 'inherit', shell: true, cwd: cwd });
    child.on('close', (code) => {
      if (code === 0) {
        resolve(code);
      } else {
        reject(new Error(`${command} ${args.join(' ')} exited with code ${code}`));
      }
    });
    child.on('error', (error) => {
      reject(error);
    });
  });
}
