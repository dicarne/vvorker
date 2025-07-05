# 节点配置


## Worker访问端口

### WORKER_PORT 【开放】

通过该端口，外部可以访问到worker服务。（注意无法访问控制台）
```
WORKER_PORT=8080
```

## 控制台端口

通过这些配置，可以访问到控制台。

### MASTER_ENDPOINT

```
MASTER_ENDPOINT=http://127.0.0.1:8888  
```
该地址为master节点的地址，子节点通过该地址与master节点进行通信

### API_PORT【开放】

```
API_PORT=8888
```
该端口为开放的控制台的端口，在该端口上显示控制台网页。

以上两个端口应当保持一致（除非被容器映射成不同的端口）。

## Worker网络转发端口

### TUNNEL_HOST

```
TUNNEL_HOST=127.0.0.1
```
该地址为master节点的地址，子节点通过该地址与master节点进行内部通信。

主节点保持为127.0.0.1，子节点则为master节点的ip。

### TUNNEL_ENTRY_PORT

```
TUNNEL_ENTRY_PORT=10080
```
主节点的内部端口，如非必要请勿配置，子节点无法访问，主节点用该端口向子节点发送请求。

### TUNNEL_API_PORT【开放】

```
TUNNEL_API_PORT=18080
```
主节点应当公开给子节点的端口，子节点通过该端口配置每个worker的网络。


## Worker配置

### WORKER_URL_SUFFIX

```
WORKER_URL_SUFFIX=.example.com
```
配置域名后缀，如`.example.com`，则worker的域名将为`WORKERNAME+WORKER_URL_SUFFIX`，如`worker.example.com`。

### SCHEME

```
SCHEME=http
```
http/https，从网页访问的协议，可能影响跨域。


### NODE_NAME

```
NODE_NAME=default
```
节点名称，主节点默认为`default`，子节点需要配置的与众不同。

### AGENT_SECRET

```
AGENT_SECRET=123123
```
子节点注册master时需要的密钥，不同节点的密钥应当保持一致。

### RUN_MODE

```
RUN_MODE=master
```
主节点为`master`，子节点为`agent`。

### WORKER_HOST_MODE

```
WORKER_HOST_MODE=host
```
host/path，host模式下，worker的url为`workername.vvorker.domain`，path模式下，worker的url为`vvorker.domain/workername`。

### WORKER_HOST_PATH

```
WORKER_HOST_PATH=something
```
在path模式下，worker的url为`vvorker.domain/WORKER_HOST_PATH/workername`。在某些反向代理的情况下可能有用。

### ADMIN_API_PROXY

```
ADMIN_API_PROXY=false
```
允许admin页面代理api请求，这可能会导致路径冲突，并且WORKER_HOST_MODE必须为path。这将允许一个端口同时提供admin页面和worker服务。不推荐使用。

### DB_TYPE

```
DB_TYPE=sqlite
```
数据库类型，sqlite/postgres/mysql。
这将选择平台使用的数据库类型，注意需要配置各个数据库的连接信息。

### DB_PATH【重要】
当使用sqlite时

```
DB_PATH=/app/data/db.sqlite
```
sqlite数据库的路径。需要针对操作系统配置合适的路径。

### DB_NAME

```
DB_NAME=vvorker
```
当使用postgres/mysql时，可以指定系统数据库的名称。对于多个节点使用同一个数据库服务时较为有用。

### WORKERD_DIR【重要】

```
WORKERD_DIR=/app/data
```
workerd的目录，用于存储各个worker的运行时文件。无需持久化，仅作为缓存使用。需要针对操作系统配置合适的路径。

### WORKERD_BIN_PATH【重要】

```
WORKERD_BIN_PATH=/bin/workerd
```
workerd的二进制文件路径。需要对应位置有workerd的二进制文件，如`.exe`，注意是二进制的，而不是脚本文件。使用npm安装workerd时需要格外注意。


### API_WEB_BASE_URL

```
API_WEB_BASE_URL=http://127.0.0.1:8080
```
用于指定控制台中快速跳转服务的地址前缀。在path模式下使用。

### COOKIE_NAME

```
COOKIE_NAME=vv-authorization
```
用于指定admin的cookie的名称。

### COOKIE_AGE

```
COOKIE_AGE=86400
```
cookie的过期时间，单位为秒。

### COOKIE_DOMAIN

```
COOKIE_DOMAIN=vvorker.local
```
cookie的域名。

### ENABLE_REGISTER

```
ENABLE_REGISTER=false
```
是否允许注册。

### JWT_SECRET

```
JWT_SECRET=123123
```
jwt的密钥。

## KV 配置

### KV_PROVIDER

```
KV_PROVIDER=redis
```
kv的提供者，redis/nutsdb。选择redis通常是理想的情况，但还需要安装redis；nutsdb则为本地存储，不需要进行安装。

### SERVER_REDIS_HOST

```
SERVER_REDIS_HOST=127.0.0.1
```
redis的地址。

### SERVER_REDIS_PORT

```
SERVER_REDIS_PORT=6379
```
redis的端口。

## OSS

### SERVER_MINIO_HOST

```
SERVER_MINIO_HOST=127.0.0.1
```
minio的地址。

### SERVER_MINIO_PORT

```
SERVER_MINIO_PORT=9000
```
minio的端口。

### SERVER_MINIO_REGION

```
MINIO_REGION=us-east-1
```
minio的区域。

### SERVER_MINIO_ACCESS

```
MINIO_ACCESS=minioadmin
```
minio的访问密钥。

### SERVER_MINIO_SECRET

```
MINIO_SECRET=minioadmin
```
minio的访问密钥。

### SERVER_MINIO_USE_SSL

```
MINIO_USE_SSL=false
```
是否使用https。

### MINIO_SINGLE_BUCKET_MODE

```
MINIO_SINGLE_BUCKET_MODE=false
```
是否使用单个bucket，所有应用都使用同一个bucket下的不同文件夹，注意，这将不进行权限管控。

### MINIO_SINGLE_BUCKET_NAME

```
MINIO_SINGLE_BUCKET_NAME=vvorker
```
如果使用单个bucket，bucket名称。

## PostgreSQL

### SERVER_POSTGRE_HOST

```
SERVER_POSTGRE_HOST=127.0.0.1
```
postgres的地址。

### SERVER_POSTGRE_PORT

```
SERVER_POSTGRE_PORT=5432
```
postgres的端口。

### SERVER_POSTGRE_USER

```
SERVER_POSTGRE_USER=postgres
```
postgres的用户。

### SERVER_POSTGRE_PASSWORD

```
SERVER_POSTGRE_PASSWORD=postgres
```
postgres的密码。

## MySQL

### SERVER_MYSQL_HOST

```
SERVER_MYSQL_HOST=127.0.0.1
```
mysql的地址。

### SERVER_MYSQL_PORT

```
SERVER_MYSQL_PORT=3306
```
mysql的端口。

### SERVER_MYSQL_USER

```
SERVER_MYSQL_USER=root
```
mysql的用户。

### SERVER_MYSQL_PASSWORD

```
SERVER_MYSQL_PASSWORD=root
```
mysql的密码。

### SERVER_MYSQL_ONE_DB_NAME

```
SERVER_MYSQL_ONE_DB_NAME=vvorker
```
当不为空时，所有mysql资源都将在同一个库中，并且不进行权限控制。

## 本地内部端口

注意这些端口不要重复即可，通常一台机器上不使用容器部署多个节点时，需要配置不同的端口。

### CLIENT_POSTGRE_PORT

```
CLIENT_POSTGRE_PORT=15432
```
postgres的端口。

### CLIENT_MYSQL_PORT

```
CLIENT_MYSQL_PORT=15433
```
mysql的端口。

### CLIENT_REDIS_PORT

```
CLIENT_REDIS_PORT=16379
```
redis的端口。

### LOCAL_TMP_POSTGRE_PORT

```
LOCAL_TMP_POSTGRE_PORT=13420
```
tmp postgres的端口。

### LOCAL_TMP_REDIS_PORT

```
LOCAL_TMP_REDIS_PORT=13421
```
tmp redis的端口。

### LOCAL_TMP_MYSQL_PORT

```
LOCAL_TMP_MYSQL_PORT=13422
```
tmp mysql的端口。

## SSO

### SSO_AUTH_URL

```
SSO_AUTH_URL=http://127.0.0.1:8888/sso/auth
```
sso认证地址。

### SSO_REDIRECT_URL

```
SSO_REDIRECT_URL=http://127.0.0.1:8888/sso/redirect
```
登录页。

### SSO_BASE_URL

```
SSO_BASE_URL=http://127.0.0.1:8888
```
sso基础地址。`xxxx.com/WORKER_HOST_PATH/SSO_BASE_URL/WORKER_NAME`
在服务器不配置WORKER_HOST_PATH时，但前端需要BASE_URL时可能有用。

### SSO_COOKIE_NAME

```
SSO_COOKIE_NAME=vv-sso
```
sso cookie名称。
