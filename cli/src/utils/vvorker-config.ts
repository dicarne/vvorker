import * as fs from 'fs-extra';
import json5 from 'json5';
import { config } from './config';

export function loadVVorkerConfig() {
  let current_env = config.current_env
  if (fs.existsSync(`${process.cwd()}/vvorker.${current_env}.json`)) {
    let vvorkerJson = json5.parse(fs.readFileSync(`${process.cwd()}/vvorker.${current_env}.json`, 'utf-8'))
    return vvorkerJson
  } else {
    return json5.parse(fs.readFileSync(`${process.cwd()}/vvorker.json`, 'utf-8'))
  }
}

export function saveVVorkerConfig(vvorkerJson: any) {
  let current_env = config.current_env
  if (fs.existsSync(`${process.cwd()}/vvorker.${current_env}.json`)) {
    fs.writeFileSync(`${process.cwd()}/vvorker.${current_env}.json`, json5.stringify(vvorkerJson, null, 2))
  } else {
    fs.writeFileSync(`${process.cwd()}/vvorker.json`, json5.stringify(vvorkerJson, null, 2))
  }
}
