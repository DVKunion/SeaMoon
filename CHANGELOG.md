# CHANGELOG

## SeaMoon 1.1.3

### ❤️ What's New

* 📝 docs: 增加手册页面sitemap站点地图(#7)(#8)
* ✨ feat(server): 修改了阿里云默认的部署资源类型(vcpu 0.1/mem 128)，来降低普通用户使用的价格消费 (#10)
* ✨ feat(server): 增加了sealos部署方案，用更加便宜的价格使用seamoon (#11)
* ✨ feat(server): 增加了docker server， 现在可以通过docker来启动服务端  (#12)
* 🔧 fix(config): 用更友好的方式来使用config，不再单一的通过域名特征来判断服务端地址类型。(#13)

**Full Changelog**: https://github.com/DVKunion/SeaMoon/compare/1.1.2...1.1.3

* 41c5ce8 feat(docker): add docker server (#12)
* 1414293 feat: low cpu && mem cost (#10)
* 99c98fd fix(client): use more friendly config (#13)

## SeaMoon 1.1.2

### ❤️ What's New

* 🔧 fix(websocket): 修正了protocol error detect 时仍挂起gorouting导致卡死的问题 (#6)
* ✨ feat(dockerfile): 增加了docker client， 现在可以通过docker来启动客户端  (#6)

**Full Changelog**: https://github.com/DVKunion/SeaMoon/compare/1.1.1...1.1.2

## SeaMoon 1.1.1

### ❤️ What's New

* 🔧 fix(websocket): 修正了 websocket 在超出 32768 slice导致的 panic。 (#4)
* 🔧 fix(websocket): 修整了 websocket 在 close 时写入 message 导致的 panic。 (#4)
* 🔧 fix(websocket): 忽略了大量 websocket 链接导致的 1006 abnormal close 报错。 (#4)
* 🔧 fix(s.yaml): 修整了 serverless-devs 工具编排文件，目前可以通过 serverless-devs 工具`s deploy`一件部署至阿里云。 (#4)
* 🔧 fix(ci): 修整了 go-releaser ci 配置 (#3)
* 🔧 fix(docs): 更新了 README.md 较为过时的使用手册。

### 🌈 Small Talk

> Hi, 各位，SeaMoon成功挤入2023Kcon兵器谱，使得整个项目获得了一批关注；在此感谢大家对SeaMoon项目的浓厚兴趣与支持，谢谢各位🙏    
> 由于工作原因，以及个人的一些想法枯竭，项目于去年创建，直到现在目前也仅支持了阿里云一个demo QAQ，因此整体给人并不是一个较为完善的使用效果。1.1.1 版本后，我会尽量保证一些活跃性质的更新，以及一些比较有意思的想法demo迭代。  
> 也欢迎对serverless感兴趣的安全小伙伴留言来交个朋友～

**Full Changelog**: https://github.com/DVKunion/SeaMoon/compare/1.1.0...1.1.1

* bc209a9 doc: update README.md
* a2e7360 fix: go-releaser ci config (#3)
* 8f51e63 fix: readme.md
* fe658ff fix: some websocket error optimization (#4)
* c316527 hotfix: some docs and code format

## 1.1.0 (2022-09-27)

### Bug Fixes

* optimize connection ([70dfc5a](https://github.com/DVKunion/SeaMoon/commit/70dfc5ad4d25fd5b529097183c873d87ec37f126))
* optimize connection ([2b416c0](https://github.com/DVKunion/SeaMoon/commit/2b416c0b106ad0a6a21aa3da838cf311061e9ef8))

## 1.0.0 (2022-09-09)

### Features

* **ci:** add build
  client ([215400c](https://github.com/DVKunion/SeaMoon/commit/215400cb7a3ae6c3f5f12df6828c8735156b810b))
* **pkg/socks5:** socks5 proxy beta
  version ([20d586c](https://github.com/DVKunion/SeaMoon/commit/20d586ce1ac36f143c1e340aa3bf9132e35af230))
* **pkg/http:** http proxy beta
  version ([3b41846](https://github.com/DVKunion/SeaMoon/commit/3b41846f75fe6d9510a9d040d76f97b35ce8c494))



