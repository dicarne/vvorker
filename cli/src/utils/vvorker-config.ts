import * as fs from 'fs-extra';
import json5 from 'json5';
import { config } from './config';
import pc from "picocolors"

export function loadVVorkerConfig() {
  let current_env = config.current_env
  if (fs.existsSync(`${process.cwd()}/vvorker.${current_env}.json`)) {
    let vvorkerJson = json5.parse(fs.readFileSync(`${process.cwd()}/vvorker.${current_env}.json`, 'utf-8'))
    return vvorkerJson
  } else {
    console.log(pc.red(`未找到配置文件 vvorker.${current_env}.json`));
    throw new Error(`未找到配置文件 vvorker.${current_env}.json`);
  }
}

export function saveVVorkerConfig(vvorkerJson: any) {
  let current_env = config.current_env
  if (fs.existsSync(`${process.cwd()}/vvorker.${current_env}.json`)) {
    fs.writeFileSync(`${process.cwd()}/vvorker.${current_env}.json`, json5.stringify(vvorkerJson, null, 2))
  } else {
    fs.writeFileSync(`${process.cwd()}/vvorker.${current_env}.json`, json5.stringify(vvorkerJson, null, 2))
    return
  }
}
