# 客户端使用手册

> 在阅读客户端使用手册之前，请确保您已经阅读并部署好了[云端服务](https://github.com/DVKunion/SeaMoon/blob/main/docs/DEPLOY.md)

## 客户端使用

> Github Action自动打包还在开发中，后续会陆续支持各种平台环境的Client，点击下载即可用。

客户端启动:  
`go mod tidy`

**http代理**  
`go run cmd/client.go -m http -l :9000 -p http://YOUR_FC_SERVER -v`

**socks5代理**
`go run cmd/client.go -m socks5 -l :9000 -p ws://YOUR_FC_SERVER -v`

证书信任:
客户端运行后，会自动在运行目录下生成证书文件。  
以mac为例，双击ca.crt，信任证书即可(原理同burp证书信任)

各参数详情:

| 参数名称          | 参数描述                                                      |  默认值  |
|---------------|-----------------------------------------------------------|:-----:|
| proxy         | 客户端运行模式: 代理模式                                             |   无   |
| -m / --mod    | 代理模式 :http/socks5                                         | http  |
| -l / --laddr  | 本地代理地址: 127.0.0.1:9000                                    | :9000 | 
| -p / --paddr  | 云端代理地址: http://xxxxxxx.xxxx.xxxx ｜ ws://xxxxxxx.xxxx.xxxx |   无   |
| -v /--verbose | 是否展示代理日志详情                                                | false |

