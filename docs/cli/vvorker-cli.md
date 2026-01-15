# VVORKER CLI

## 安装

```bash
pnpm install -g @dicarne/vvcli
```

## init

初始化项目，生成`vvorker.json`。

```bash
vvcli init <your project name>
```

并且同时需要输入 vvorker 平台上的项目 uid 以进行绑定，可以创建一个空项目并将 uid 复制过来。

可以选择 api 模式还是 vue 模式，api 模式没有前端网页；vue 模式为前后端分离项目，但不需要单独部署。

vue 模式需要配置`vvorker.json`中的`assets`字段。

```json
// vvorker.json
{
  "assets": [
    {
      "directory": "./dist",
      "binding": "ASSETS"
    }
  ]
}
```

## types

根据`vvorker.json`中的绑定信息生成类型文件。

```bash
vvcli types
```

## deploy

发布到节点。

```bash
vvcli deploy
```

该命令将编译代码、网页并上传。

## create

交互式方式创建新环境。

```bash
vvcli create
```

填写环境名、url、token。

url指的是平台的根地址url。

## env

交互式方式切换环境，不同的环境，vvorker 环境的 url、token 均不同。

```bash
vvcli env
```

## dev

等价于`pnpm run dev`，便于统一命令行。

```bash
vvcli dev
```

## config

配置当前环境。

### set

设置当前环境。

```bash
vvcli config set env <env name>
```

设置当前环境的 url

```bash
vvcli config set url <url>
```

设置当前环境的 token

```bash
vvcli config set token <token>
```
