import * as fs from 'fs-extra';

export interface EnvConfig {
  url?: string;
  token?: string;
}

export interface Config {
  current_env: string;
  env: { [key: string]: EnvConfig };
}

const userHome = require('os').homedir();
const vvcliDir = `${userHome}/.vvcli`;
fs.ensureDirSync(vvcliDir)

const configFilePath = `${vvcliDir}/config.json`;
if (!fs.existsSync(configFilePath)) {
  fs.writeFileSync(configFilePath, JSON.stringify({ env: {}, current_env: 'default' }));
}

export let config = JSON.parse(fs.readFileSync(configFilePath, 'utf-8')) as Config;

export function saveConfig() {
  fs.writeFileSync(configFilePath, JSON.stringify(config, null, 2));
}

export function getToken() {
  let env = config.current_env ?? "default";
  return config.env[env]?.token;
}

export function getUrl() {
  let env = config.current_env ?? "default";
  let url = config.env[env]?.url;
  if (url?.endsWith('/')) {
    url = url.slice(0, -1)
  }
  return url
}

export function ensureEnv(env: string) {
  if (!config.env) {
    config.env = {

    }
  }
  if (!config.env[env]) {
    config.env[env] = {
    }
  }
}

export function setUrl(url: string) {
  let env = config.current_env ?? "default";
  ensureEnv(env);
  config.env[env].url = url;
}
