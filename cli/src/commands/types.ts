import { Command } from 'commander';
import * as fs from 'fs-extra';
import { loadVVorkerConfig } from '../utils/vvorker-config';
import { apiClient } from '../utils/api';
import { withWorkingDir } from '../utils/working-dir';

export const typesCommand = new Command('types')
  .description("用于自动生成配置文件对应的TypeScript类型")
  .action(async () => {
    await withWorkingDir(async () => {
      const vvv = loadVVorkerConfig();
      let resp = await apiClient.post(`/api/ext/types`, vvv)
      if (resp.status !== 200) {
        throw new Error(`获取类型失败：${resp.status} ${resp.statusText}`);
      }
      let w = resp.data.data.types;
      if (vvv.project.type === "vue") {
        fs.writeFileSync(`${process.cwd()}/server/binding.ts`, w);
      } else {
        fs.writeFileSync(`${process.cwd()}/src/binding.ts`, w);
      }
      fs.writeFileSync(`${process.cwd()}/vvorker-schema.json`, resp.data.data.schema)
      console.log('类型生成成功');
    });
  });
