<p align="center">
    <img src="https://cdn.dvkunion.cn/SeaMoon/logo.png" width="360" alt="logo"/>
</p>
<h1 align="center">Sea Moon</h1>

<p align="center">

<img src="https://img.shields.io/github/stars/DVKunion/SeaMoon.svg"  alt="stars"/>
<img src="https://img.shields.io/github/languages/top/DVKunion/SeaMoon.svg?&color=red" alt="languages">
<img src="https://img.shields.io/github/license/DVKunion/SeaMoon.svg"  alt="license"/>
</p>

<p align="center">
    月海(Sea Moon) 是一款 FaaS/BaaS 实现的 Serverless 云渗透工具集，致力于开启云原生的渗透模式。  
</p>
<p align="center">
    月海之名取自于苏轼的《西江月·顷在黄州》，寓意月海取自于传统安全工具，用之于云，最终达到隐匿于海的效果。
</p>
<p align="center"><b>
目前工具正处于开发中，欢迎各位提交</b> <a href="https://github.com/DVKunion/SeaMoon/issues" ><b>Issue</b></a> |  <a href="https://github.com/DVKunion/SeaMoon/pulls"><b>Pr</b></a>
</p>

<br />

## ☁️ 什么是月海

> 🌕 月出于云却隐于海

月海(Sea Moon) 是一款使用 FaaS/BaaS 实现的的 Serverless 渗透工具，致力于开启云原生的渗透模式。

说人话的总结，月海其实就是一款利用云函数来隐匿攻击行踪以及分布式处理扫描任务的集成工具，

## 🌟 月海能做什么

### 网络层

网络层支持是月海的基础功能，也是云函数最基本的优势和特性。 支持一级代理 / 链式代理，

| 代理类型      | 原理文档                                                                          | 服务端支持 | 客户端支持 |
|-----------|-------------------------------------------------------------------------------|:-----:|:-----:|
| HTTP      | [HTTP.md](https://github.com/DVKunion/SeaMoon/blob/main/docs/net/HTTP.md)     |   ✅   |   ✅   |
| HTTPS     | [HTTP.md](https://github.com/DVKunion/SeaMoon/blob/main/docs/net/HTTP.md)     |   ✅   |   ✅   |
| Socks5    | [Socks5.md](https://github.com/DVKunion/SeaMoon/blob/main/docs/net/SOCKS5.md) |   ✅   |   ✅   |
| SS/SSR    | [Socks5.md](https://github.com/DVKunion/SeaMoon/blob/main/docs/net/SOCKS5.md) | 🐷待开发 | 🐷待开发 |
| VMess     | [Socks5.md](https://github.com/DVKunion/SeaMoon/blob/main/docs/net/SOCKS5.md) | 🐷待开发 | 🐷待开发 |
| websocket |                                                                               | 🐷待开发 | 🐷待开发 |
| 链式代理      |                                                                               | ❌暂无计划 | ❌暂无计划 |

### 应用层

月海的应用层能力是基于网络层的思考基础上，实现的真正上层渗透业务，例如：端口探测，网络反馈等

| 能力名称         | 原理文档                                                                              | 服务端支持 | 客户端支持 |
|--------------|-----------------------------------------------------------------------------------|:-----:|:-----:|
| 动态WebShell隐匿 | [WebShell.md](https://github.com/DVKunion/SeaMoon/blob/main/docs/app/WEBSHELL.md) | 🐷待开发 | 🐷待开发 |
| 分布式扫描        |                                                                                   | 🐷待开发 | 🐷待开发 |
| 反弹Shell代理    |                                                                                   | 🐷待开发 | 🐷待开发 |
| CI容器云扫描利用    |                                                                                   | 🐷待开发 | 🐷待开发 |

### 其他特性

+ 身份认证加强保密性: 🐶开发中
+ 探活机制/心跳检测: 🐷待开发
+ 多云平台/区域环境部署后随机选择机制: 🐷待开发
+ 精美的客户端web控制台: ? MayBe

## 🕹 ️开始使用

[三步部署月海(Sea Moon)到阿里云](https://github.com/DVKunion/SeaMoon/blob/main/docs/DEPLOY.md)

[开启月海客户端使用](https://github.com/DVKunion/SeaMoon/blob/main/docs/START.md)

## ❗ 免责声明

本工具仅用于学习serverless以及云原生相关技术，请勿用于其他用途。

如您在使用本工具的过程中存在任何非法行为，您需自行承担相应后果，我们将不承担任何法律及连带责任。

## 参考文献与项目

感谢各位前辈师傅们的分享与沉淀。

**文档类**

+ [浅谈云函数的利用面](https://xz.aliyun.com/t/9502)
+ [白嫖CDN，打造封不尽IP的代理池](https://freewechat.com/a/MzI0MDI5MTQ3OQ==/2247484068/1)
+ [Serverless 应用开发指南](https://serverless.ink/)
+ [HTTP被动扫描代理的那些事](https://www.freebuf.com/articles/web/212382.html)
+ [Subsocks: 用GO实现一个Socks5安全代理](https://luyuhuang.tech/2020/12/02/subsocks.html)

**项目类**

+ [SFCProxy](https://github.com/shimmeris/SCFProxy)
+ [go-socks5](https://github.com/armon/go-socks5)
+ [subsocks](https://github.com/luyuhuang/subsocks)
+ [gost](https://github.com/ginuerzh/gost)
+ [InCloud](https://github.com/inbug-team/InCloud)
+ [sfc-proxy](https://github.com/Sakurasan/scf-proxy)
+ [Serverless-transitcode](https://github.com/copriwolf/serverless-transitcode)
+ [protoplex](https://github.com/SapphicCode/protoplex)