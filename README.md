<p align="center">
    <img src="https://seamoon.oss-cn-hangzhou.aliyuncs.com/logo.png" width="360" alt="logo"/>
</p>
<h1 align="center">Sea Moon</h1>

<p align="center">
<img src="https://goreportcard.com/badge/github.com/DVKunion/SeaMoon" alt="go-report"/>
<img src="https://img.shields.io/github/languages/top/DVKunion/SeaMoon.svg?&color=blueviolet"
     alt="languages"/>
<img src="https://img.shields.io/badge/LICENSE-MIT-777777.svg" alt="license"/>
<img src="https://img.shields.io/github/downloads/dvkunion/seamoon/total?color=orange" alt="downloads"/>
<img src="https://img.shields.io/github/stars/DVKunion/SeaMoon.svg" alt="stars"/>
</p>

<p align="center">
    月海(Sea Moon) 是一款 FaaS/BaaS 实现的 Serverless 网络工具集，期望利用云原生的优势，实现更简单、更便宜的网络功能。
</p>
<p align="center">
    月海之名取自于苏轼的《西江月·顷在黄州》，寓意月海取自于传统工具，用之于云，最终达到隐匿于海的效果。
</p>

## ☁️ 什么是月海

> 🌕 月出于云却隐于海

月海(Sea Moon) 是一款 FaaS/BaaS 实现的 Serverless 网络工具集，期望利用云原生的优势，实现更简单、更便宜的网络工具。

月海之名取自于苏轼的《西江月·顷在黄州》，寓意月海取自于传统工具，用之于云，最终达到隐匿于海的效果。

月海基于 Serverless 的动态与无状态的特性，从网络层实现了一个基于 Serverless 的网络工具集，包括代理、转发、隧道等等常见网络功能；
同时在客户端集成了大量云厂商，实现快捷的一键式部署和跨厂商与平台操作。

想要了解更多，请移步 [官方手册](https://seamoon.dvkunion.cn)

觉得项目不错的话，[还请给一个star ✨](https://github.com/DVKunion/SeaMoon), 你的支持是更新的最大动力～

## 🌟 月海能做什么

Serverless 的动态实例不同的出口IP，从而获取到了干净(非威胁情报黑名单)、随机的外网IP代理、用后即销毁的无痕状态等。

**网络代理**

| 代理类型        | 技术文档                                                      | Seamoon 客户端支持 | 其他客户端支持 |
|-------------|-----------------------------------------------------------|:-------------:|:-------:|
| HTTP(S)     | [HTTP.md](https://seamoon.dvkunion.cn/tech/net/http/)     |       ✅       |    ✅    |
| Socks5      | [Socks5.md](https://seamoon.dvkunion.cn/tech/net/socks5/) |       ✅       |    ✅    |
| Socks4      | []()                                                      |       ❌       |    ✅    |
| Vmess       | []()                                                      |       ✅       |    ✅    |
| Vless       | []()                                                      |       ✅       |    ✅    |
| shadowsocks | []()                                                      |       ✅       |    ✅    | 
 |

**网络隧道**

| 隧道类型      | 技术文档 | 支持情况  |
|-----------|------|:-----:|
| websockst | []() |   ✅   |
| grpc      | []() |   ✅   |
| oss       | []() | 🐷调研中 |


**代理模式**
+ 正向代理
+ 反向代理
+ 端口转发

**其他**

+ 💻 多客户端支持，clash/shadowrocket 等。
+ 🧅 Tor 网络 .onion 支持. [如何开启 Tor 代理](https://seamoon.dvkunion.cn/guide/client/tor/)
+ ......

更多特性相关请移步: [技术文档](https://seamoon.dvkunion.cn/tech/feature/)

## 🧭 支持平台

| 平台名称     |            免费力度            | 是否支持  |   
|----------|:--------------------------:|:-----:|
| 阿里云      |           新用户三个月           |   ✅   |
| 腾讯云      |            🈚️             |   ✅   | 
| Sealos   |            五元余额            |   ✅   | 
| 华为云      |                            |   ✅   |  
| 百度云      |                            |   ✅   | 
| 🙅Render | ~~每月750小时免费 + 100G流量~~但是封号 |   ❌   |
| AWS      |                            | 🐷调研中 |  
| Google   |                            | 🐷调研中 | 

## 🕹开始使用

[继续阅读: 快速开始](https://seamoon.dvkunion.cn/guide/start)

## 💻 技术文档

[🧑‍💻 技术文档](https://seamoon.dvkunion.cn/tech/feature)

## 🛜 使用展示

**登陆认证**
![login](https://seamoon.oss-cn-hangzhou.aliyuncs.com/62564a7263484cddb622d27abf09e4ed.png)

**代理管控**
![proxy](https://seamoon.oss-cn-hangzhou.aliyuncs.com/a473e1b3a2cd45379737bba56bc9cb8b.png)

**函数管控**
![func](https://seamoon.oss-cn-hangzhou.aliyuncs.com/ac38d83adf69439baf694f6705b3f9f4.png)

**账户管控**
![account](https://seamoon.oss-cn-hangzhou.aliyuncs.com/ea911c9b2f3c4fb886f04f7043a6e5f9.png)

## ❗ 免责声明

本工具仅用于学习serverless以及云原生相关技术，请勿用于其他用途。

如您在使用本工具的过程中存在任何非法行为，您需自行承担相应后果，我们将不承担任何法律及连带责任。
