# VVORKER CLI

## 安装

```
pnpm install -g @dicarne/vvcli
```

## 指令

### init

初始化项目，生成`vvorker.json`。
```
vvcli init <your project name>
```
并且同时需要输入vvorker平台上的项目uid以进行绑定，可以创建一个空项目并将uid复制过来。

可以选择api模式还是vue模式，api模式没有前端网页；vue模式为前后端分离项目，但不需要单独部署。

vue模式需要配置`vvorker.json`中的`assets`字段。
```json
// vvorker.json
{
    "assets": [
        {
        "directory": "./dist",
        "binding": "ASSETS"
        }
    ],
}
```

### types

根据`vvorker.json`中的绑定信息生成类型文件。
```
vvcli types
```

### deploy

发布到节点。
```
vvcli deploy
```
该命令将编译代码、网页并上传。

### env
交互式方式切换环境，不同的环境，vvorker环境的url、token均不同。
```
vvcli env
```

### config

配置当前环境。

#### set

设置当前环境。
```
vvcli config set env <env name>
```

设置当前环境的url
```
vvcli config set url <url>
```

设置当前环境的token
```
vvcli config set token <token>
```

