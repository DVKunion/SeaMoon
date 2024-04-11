module.exports = {
    title: "SeaMoon",
    theme: 'vdoing',
    description: "月海(Sea Moon) 是一款 FaaS/BaaS 实现的 Serverless 网络工具集，期望利用云原生的优势，实现更简单、更便宜的网络功能。",

    head: [ // 注入到页面<head> 中的标签，格式[tagName, { attrName: attrValue }, innerHTML?]
        ['link', {rel: 'icon', href: '/img/favicon.ico'}], //favicons，资源放在public文件夹
        ['meta', {name: 'keywords', content: 'serverless,proxy,pentest,seamoon'}],
        ['meta', {name: 'baidu-site-verification', content: 'codeva-vXPumeNBPL'}],
        ['script', {}, `
        var _hmt = _hmt || [];
        (function() {
          var hm = document.createElement("script");
          hm.src = "https://hm.baidu.com/hm.js?7dac4248d29ddaacd4b3c0b71d9b2015";
          var s = document.getElementsByTagName("script")[0]; 
          s.parentNode.insertBefore(hm, s);
        })();
        </script>        
        `],
    ],

    themeConfig: {
        defaultMode: "dark",
        bodyBgImgOpacity: 1.0,
        nav: [
            {text: '首页', link: '/'},
            {
                text: '使用手册', link: '/guide/introduce/'
            },
            {
                text: '技术博客', link: '/tech/'
            }
        ],
        archive: false,
        category: false,
        tag: false,
        updateBar: { // 最近更新栏
            showToArticle: false, // 显示到文章页底部，默认true
        },
        sidebar: 'structuring', //  'structuring' | { mode: 'structuring', collapsable: Boolean} | 'auto' | 自定义
        sidebarOpen: true,
        searchMaxSuggestions: 10,
        repo: 'Dvkunion/SeaMoon',
        author: { // 文章默认的作者信息，可在md文件中单独配置此信息 String | {name: String, href: String}
            name: 'DVKunion', // 必需
            href: 'https://github.com/DVKunion' // 可选的
        },
        social: { // 社交图标，显示于博主信息栏和页脚栏
            icons: [
                {
                    iconClass: 'icon-youjian',
                    title: '发邮件',
                    link: 'mailto:dvkunion@gmail.com'
                },
                {
                    iconClass: 'icon-mao',
                    title: '放一只猫',
                    link: 'https://github.com/DVKunion/SeaMoon'
                },
                {
                    iconClass: 'icon-github',
                    title: 'GitHub',
                    link: 'https://github.com/Dvkunion'
                },
                {
                    iconClass: 'icon-weixin',
                    title: 'weixin',
                    link: ''
                },
            ]
        },
        footer: { // 页脚信息
            createYear: 2022, // 博客创建年份
            copyrightInfo: 'DVKunion | MIT License', // 博客版权信息，支持a标签
        },
        extendFrontmatter: {
            article: false
        }
    },
}