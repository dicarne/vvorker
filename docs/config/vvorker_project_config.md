# 项目配置

在每个项目的根目录中都存在一个或多个`vvorker.{envname}.json`文件，用于存储项目的配置信息。
该文件至关重要，声明了资源的绑定、环境变量、其他配置等。

## 通用部分

```json
{
    "$schema": "vvorker-schema.json",   // 配置文件 schema，通过 vvcli types 自动生成
    "name": "huanbao",                  // 项目名称，请和 wrangler.jsonc 中的 name 保持一致
    "project": {
      "uid": "6eb40d7735f7499dae93d40a31912905",    // 项目UID，请和vvorker控制台中创建的项目UID保持一致
      "type": "vue"                                 // 项目类型，自动生成
    },
    "version": "0.0.2",                             // 项目版本号，请使用vvcli ver更新
    ...
}
```

## 绑定服务

用于绑定内部服务调用，绑定成功后，无需鉴权即可安全的调用内部服务。

services字段中填写想要绑定的服务名称。填写后，需要到对应服务控制台的“鉴权”页面，添加“内部访问”，并填写本服务的名称。

当填写完毕后，请使用`vvcli types`生成类型文件，以便在代码中使用类型提示。短横将被消除，每个单词将首字母大写。

```json
{
    "services": [
        "myservice", "another-service"  // 生成的名称： MyService, AnotherService
    ],
}
```

## 环境变量

环境变量可以用于在不同环境中使用不同的配置。在修改环境变量后，请使用`vvcli types`生成类型文件，以便在代码中使用类型提示。

环境变量的值可以为任何合法的JSON值，包括字符串、数字、布尔值、数组、对象等。

```json
{
    "envs": {
        "MY_ENV_VAR": "my-value",
        "MY_SECRET_VAR": {
            "type": "secret",
            "value": "my-secret-value"
        },
        "MODE": "development"  // 特殊环境变量值，用于区分开发环境和生产环境。当处于开发环境时，将暴露调试端点。
    }
}
```

## 秘密变量

秘密变量用于存储敏感信息，如数据库密码、API密钥等。

秘密变量请前往对应worker的控制台进行配置，在项目的配置文件中只需要填写变量名称和空字符串即可。

秘密变量必须为字符串，不支持其他类型。

## OSS

OSS用于存储静态资源，如图片、音频、视频等。

```json
{
    "oss": [{
        "resource_id": "aaaaaaaaaaaaaaaaaaaaa",         // 在vvorker控制台中创建的oss资源id
        "binding": "myoss"                              // 绑定的变量名称，用于在代码中访问
    }]
}
```

## KV
KV用于存储键值对，如配置、状态等。

```json
{
    "kv": [{
        "resource_id": "bbbbbbbbbbbbbbbbbbbbbbb",         // 在vvorker控制台中创建的kv资源id
        "binding": "mykv"                               // 绑定的变量名称，用于在代码中访问
    }]
}
```

## mysql/pgsql
mysql/pgsql用于存储关系型数据库，如MySQL、PostgreSQL等。

```json
{
    "mysql": [{
        "resource_id": "ccccccccccccccccccccccc",         // 在vvorker控制台中创建的mysql资源id
        "binding": "mymysql",                             // 绑定的变量名称，用于在代码中访问
        "migrate": "./server/db/drizzle"                  // 数据库迁移文件路径
    }]
}
```

```json
{
    "pgsql": [{
        "resource_id": "ddddddddddddddddddddddd",         // 在vvorker控制台中创建的pgsql资源id
        "binding": "mypgsql",                             // 绑定的变量名称，用于在代码中访问
        "migrate": "./server/db/drizzle"                  // 数据库迁移文件路径
    }]
}
```

## assets

assets用于存储静态资源，如图片、音频、视频等。通常用于展示网页。

```json
{
    "assets": [
        {
            "directory": "./dist/client",
            "binding": "ASSETS"
        }
    ]
}
```

## proxy

proxy用于代理服务器，如果你的内网中有一台通往互联网的服务器，可以使用proxy进行转发。

```json
{
    "proxy": [{
        "binding": "internet",
        "address": "IP:PORT",
        "type": "http"
    }]
}
```

## 其他配置

```json
{
    "compatibility_flags": [
        "nodejs_compat"             // 参考 Cloudflare Workers 兼容性标志
    ]
}
```