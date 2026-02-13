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
  fs.writeFileSync(configFilePath, JSON.stringify({ env: {}, current_env: '' }));
}

export let config = JSON.parse(fs.readFileSync(configFilePath, 'utf-8')) as Config;

export function saveConfig() {
  fs.writeFileSync(configFilePath, JSON.stringify(config, null, 2));
}

export function getToken() {
  if (!config.current_env) {
    return undefined;
  }
  return config.env[config.current_env]?.token;
}

export function getUrl() {
  if (!config.current_env) {
    return undefined;
  }
  let url = config.env[config.current_env]?.url;
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
  if (!config.current_env) {
    throw new Error("当前没有选择环境，请先使用 vvcli create 创建环境");
  }
  ensureEnv(config.current_env);
  config.env[config.current_env].url = url;
}

export function getEnv() {
  if (!config.current_env) {
    throw new Error("当前没有选择环境，请先使用 vvcli create 创建环境");
  }
  return config.current_env
}