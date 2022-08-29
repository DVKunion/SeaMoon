<p align="center">
    <img src="https://cdn.dvkunion.cn/github/logo.png" width="360" alt="logo"/>
</p>
<h1 align="center">Sea Moon</h1>

---

<p align="center">

<img src="https://img.shields.io/github/stars/DVKunion/SeaMoon.svg"  alt="stars"/>
<img src="https://img.shields.io/github/languages/top/DVKunion/SeaMoon.svg?&color=red" alt="languages">
<img src="https://img.shields.io/github/license/DVKunion/SeaMoon.svg"  alt="license"/>
</p>


<p align="center">
    月海(Sea Moon) 是一款使用 FaaS/BaaS 实现的的 Serverless 云渗透工具集，致力于开启云原生的渗透模式。  
</p>
<p align="center">
    月海之名取自于苏轼的《西江月·顷在黄州》，寓意云安全工具取自于传统安全工具，用之于云，最终达到隐匿于海的效果。
</p>

--- 

## ☁️ 什么是月海

> 🌕 月出于云却隐于海


月海(Sea Moon) 是一款使用 FaaS/BaaS 实现的的 Serverless 渗透工具，致力于开启云原生的渗透模式。

说人话的总结，月海其实就是一款利用云函数来隐匿攻击行踪以及分布式处理扫描任务的集成工具，

## 🌟 月海能做什么

### 网络层

网络层支持是月海的基础功能，也是云函数最基本的优势和特性。 支持一级代理 / 链式代理，

+ HTTP/HTTPS 代理
+ C2 代理
+ DNS 隧道代理
+ Socks5 代理

### 应用层

月海的应用层能力是基于网络层的思考基础上，实现的真正上层渗透业务，例如：端口探测，网络反馈等

+ 分布式扫描: 如:Nmap端口探测。解耦过于复杂的网络探测/扫描类任务
+ XSS小工具，属于您自己的专属xss payload平台
+ DNSLog 小工具，属于您自己的 DNSLog 平台
+ Shell代理

## 🕹 ️开始使用

## ❗ 免责声明

## 参考文献与项目

感谢各位前辈师傅们的分享与沉淀。

**文档类**

+ [浅谈云函数的利用面](https://xz.aliyun.com/t/9502)
+ [白嫖CDN，打造封不尽IP的代理池](https://freewechat.com/a/MzI0MDI5MTQ3OQ==/2247484068/1)

**项目类**

+ [SFCProxy](https://github.com/shimmeris/SCFProxy)
+ [GOProxy](https://github.com/snail007/goproxy)
+ [sfc-proxy](https://github.com/Sakurasan/scf-proxy)
+ [protoplex](https://github.com/SapphicCode/protoplex)