# SOCKS5

## 原理

基础理论: [浅谈云函数的利用面](https://xz.aliyun.com/t/9502)

在云函数(FC)的限制下，大佬提出了一种通过vps建立起socks5隧道的模式，从模式上来看，更像是一种反向连接。

但是这种模式，需要一台VPS。对于穷逼的脚本小子的我，实在是不够优雅。

FC的不成熟的确限制了大部分的玩法，比如触发器种类，比如协议，比如端口限制等等。

在这种大环境下，我们无力去变更云函数的生态(其实也有可能云函数并没有为我们这种使用方式进行设计)，只能自寻出路。

想要优雅的正向连接，只能在HTTP上做文章。

突然联想到早些年，做安全服务时拿到了WebShell后如何进行内网渗透？这就想起了一个利器工具，也是我们今天的主角：

[reGeorg](https://github.com/sensepost/reGeorg)

这个工具提供了各种语言的脚本，能够通过HTTP隧道的方式，结合本地客户端，建立socks连接代理。

他的原理其实是依赖于，socks属于会话层，是对TCP/IP协议的封装，而在应用层的HTTP协议也是同样属于对TCP/IP协议的封装。

通俗来说，socks就是爸爸，而HTTP只是他众多的子类而已，相互之间的转化是存在某种方式的。

举个例子，如python中的urllib库，底层就是使用sockets实现的HTTP。

因此，我们云函数socks代理的模型就可以画出来了:

用户 -> socks -> client -> 转化为HTTP -> FC云函数 -> 解析HTTP -> 发送socks

用户 <- 转化为socks <- client <- 转化为HTTP <- FC云函数 <- socks数据

我们的client开启一个socks的监听，然后将监听到的数据转化为http请求发给fc处理， fc根据http提供的数据发起socks连接，获取数据。之后fc函数再将数据通过http
返回byte字节码，client端接收到响应，再根据协议降级为socks。

理论存在，实践开始。 根据原理分析，我们要做的事情就比较明显了：

+ 在云函数部署好一个接受HTTP响应，并转化为socks连接的服务
+ 在本地启动client端，监听一个socks端口，将该端口的数据按照协议转化为HTTP请求发送给云函数

我们参考reGeorg的重构版[Neo-reGeorg](https://github.com/L-codes/Neo-reGeorg), 写出我们的客户端和服务端。

