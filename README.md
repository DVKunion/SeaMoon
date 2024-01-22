<p align="center">
    <img src="https://seamoon.oss-cn-hangzhou.aliyuncs.com/logo.png" width="360" alt="logo"/>
</p>
<h1 align="center">Sea Moon</h1>

<p align="center">
<img src="https://goreportcard.com/badge/github.com/DVKunion/SeaMoon" />
<img src="https://img.shields.io/github/stars/DVKunion/SeaMoon.svg"  alt="stars"/>
<img src="https://img.shields.io/github/downloads/dvkunion/seamoon/total?color=orange" alt="downloads" />
<img src="https://img.shields.io/github/languages/top/DVKunion/SeaMoon.svg?&color=blueviolet" alt="languages">
<img src="https://img.shields.io/badge/LICENSE-MIT-777777.svg"  alt="license"/>
</p>

<p align="center">
    月海(Sea Moon) 是一款 FaaS/BaaS 实现的 Serverless 网络工具集，期望利用云原生的优势，实现更简单、更便宜的网络功能。
</p>
<p align="center">
    月海之名取自于苏轼的《西江月·顷在黄州》，寓意月海取自于传统工具，用之于云，最终达到隐匿于海的效果。
</p>

## ☁️ 什么是月海

> 🌕 月出于云却隐于海

月海(Sea Moon) 是一款 FaaS/BaaS 实现的 Serverless 网络工具集，期望利用云原生的优势，实现更简单、更便宜的网络功能。

月海之名取自于苏轼的《西江月·顷在黄州》，寓意月海取自于传统工具，用之于云，最终达到隐匿于海的效果。

月海基于 Serverless 的动态特性、分别从网络层、应用层实现对应的能力，并集成了SDK来实现快捷的部署/跨厂商操作。

想要了解更多，请移步 [官方手册](https://seamoon.dvkunion.cn)

觉得项目不错的话，[还请给一个star ✨](https://github.com/DVKunion/SeaMoon), 你的支持是更新的最大动力～

## 🔔 V2.0 大更新

最近在测试 2.0 大版本的更新发布，但是 2.0 更新跨度太大，需要逐步测试完善后再一同和入主线分支，来保证主线代码正常工作。

2.0 版本的提前试用再 2.0-dev 分支，增加和变动了很多功能如：

+ 客户端改版升级
+ 支持自动部署服务端
+ 认证相关功能增加
+ 可视化流量统计等等

预计近期完成测试后就会合并进来, 这里放几张预告图:

**登陆认证**
![login](https://seamoon.oss-cn-hangzhou.aliyuncs.com/62564a7263484cddb622d27abf09e4ed.png)

**代理管控**
![proxy](https://seamoon.oss-cn-hangzhou.aliyuncs.com/a473e1b3a2cd45379737bba56bc9cb8b.png)

**函数管控**
![func](https://seamoon.oss-cn-hangzhou.aliyuncs.com/ac38d83adf69439baf694f6705b3f9f4.png)

**账户管控**
![account](https://seamoon.oss-cn-hangzhou.aliyuncs.com/ea911c9b2f3c4fb886f04f7043a6e5f9.png)


## 🕹 快速开始

[⚡️ 快速开始](https://seamoon.dvkunion.cn/guide/start)

[✨ 新增Render部署, ~~完全免费的套餐且功能不减！！！~~ 免费但是封号的SB ](https://seamoon.dvkunion.cn/guide/deploy/render)

[✨ 新增Sealos部署, 价格更低更实惠！](https://seamoon.dvkunion.cn/guide/deploy/sealos)

[🔥 新增 seamoon 支持 tor 网络访问 onion 域名](https://seamoon.dvkunion.cn/guide/client/tor/)

## 🧭 支持平台

| 平台名称     |            免费力度            | 是否支持  |                          部署文档                          | 
|----------|:--------------------------:|:-----:|:------------------------------------------------------:|
| 阿里云      |           新用户三个月           |   ✅   | [部署](https://seamoon.dvkunion.cn/guide/deploy/aliyun)  |
| 腾讯云      |            🈚️             |   ✅   | [部署](https://seamoon.dvkunion.cn/guide/deploy/tencent) |
| 🔥Sealos |         五元余额，流量不计费         |   ✅   | [部署](https://seamoon.dvkunion.cn/guide/deploy/sealos)  |
| 🙅Render | ~~每月750小时免费 + 100G流量~~但是封号 |   ✅   | [部署](https://seamoon.dvkunion.cn/guide/deploy/render)  |
| 华为云      |                            |   ✅   |                         [部署]()                         |
| 百度云      |                            |   ✅   |                         [部署]()                         |
| AWS      |                            | 🐷调研中 |                         [部署]()                         |
| Google   |                            | 🐷调研中 |                         [部署]()                         |

## 🔥 使用展示

![client](https://seamoon.oss-cn-hangzhou.aliyuncs.com/0dd37f5600364e59a9457e38eaf77b1f.png)

## 🔨 功能简要

**网络隧道**

| 隧道类型      | 技术文档 | 支持情况 |
|-----------|------|:----:|
| websockst | []() |  ✅   |
| grpc      | []() |  ✅   |

**网络代理**

| 代理类型        | 技术文档                                                      | Seamoon 客户端支持 | 其他客户端支持 |
|-------------|-----------------------------------------------------------|:-------------:|:-------:|
| HTTP(S)     | [HTTP.md](https://seamoon.dvkunion.cn/tech/net/http/)     |       ✅       |    ✅    |
| Socks5      | [Socks5.md](https://seamoon.dvkunion.cn/tech/net/socks5/) |       ✅       |    ✅    |
| Socks4      | []()                                                      |       ❌       |    ✅    |
| shadowsocks | []()                                                      |       ❌       |    ✅    |

**其他**

+ 🧅 Tor 网络 .onion 支持. [如何开启 Tor 代理](https://seamoon.dvkunion.cn/guide/client/tor/)
+ ......

## 💻 技术文档

[🧑‍💻 技术文档](https://seamoon.dvkunion.cn/tech/feature)

## ❗ 免责声明

本工具仅用于学习serverless以及云原生相关技术，请勿用于其他用途。

如您在使用本工具的过程中存在任何非法行为，您需自行承担相应后果，我们将不承担任何法律及连带责任。

## ☎️ 联系方式

Serverless in Cybersecurity 讨论微信群:

也可以在 项目 Discuss 处发起讨论。

## ✨ CTStack

<img src="https://ctstack-oss.oss-cn-beijing.aliyuncs.com/CT%20Stack-2.png" width="30%" />

SeaMoon 现已加入 [CTStack](https://stack.chaitin.com/tool/detail?id=186) 社区
