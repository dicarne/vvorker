# SSO配置

## 准备服务

至少需要几个服务提供支持。

1. 登录服务（通常提供鉴权的二维码、密码框，并在登陆后跳转到指定服务页面）
2. 用户服务（通常提供具体的认证服务，如用户系统、token检验等）

### 用户服务

用户服务需要实现鉴权接口：

#### 请求值：

POST 方法，默认接受header参数。

| 参数名             | 参数类型 | 参数说明                 |
| ------------------ | -------- | ---------------------- |
| vv-sso-data        | string   | 认证信息（由路由配置）   |
| vv-sso-worker-uid  | string   | worker UID             |
| vv-sso-worker-name | string   | worker名称             |
| vv-sso-request-url | string   | 请求url                |
| vv-sso-channel     | string?  | 渠道                   |

#### 返回值：
BODY：
```typescript
interface ret {
    user_id: string
    token: string
    real_name: string
    // 扩展数据
    ext?: string
    // 重定向到其他url，用于覆盖默认配置
    redirect?: string  
    // 是否设置cookie，默认不设置
    // 只有在其他字段正确返回时才应该设置本字段，注意此时 http code 应为 200
    set_cookie?: boolean    
    // 阻止默认重定向逻辑
    // 因为默认会将服务名作为query的name
    // 某些对接第三方sso系统时需要自定义的重定向路径
    prevent_default_redirect?: boolean
}
```

HTTP CODE:
- 200: 认证成功
- 401: 未认证，需要网关进行进一步处理
- 其他: 认证失败

## 简单模式

简单模式指的是只需要重定向到某个url，并由该url进行认证，并再重定向到对应服务。

### 配置说明
编辑环境变量
```
SSO_AUTH_URL=http://xxxxxxxx    # 该url配置用户服务中的鉴权接口，将通过cookie进行鉴权
SSO_REDIRECT_URL=http://xxxxxxxxx   # 该url配置登陆服务中的登录页，通常显示二维码或者登录框
```

## 中继模式

在简单模式的基础上，主要是鉴权接口进行了更多配置。

目标是处理`http://services/xxx?token=xxxxx`的形式。

鉴权接口除了拿到cookie外，还将拿到url。若cookie有效，则直接通过；
如果cookie无效，则通过token进行认证，若成功则返回用户信息，并设置cookie。

## 多渠道模式

通过nginx对来自不同渠道的地址设置channel，或通过完整url进行识别。

鉴权失败时，可以根据渠道返回不同的重定向地址以覆盖默认地址。

鉴权成功时，也可根据渠道决定是否设置cookie。