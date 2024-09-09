package external

// ClashConfig 不同 client 结构序列化时使用
type ClashConfig struct {
	MixedPort          int                `yaml:"mixed-port,omitempty"`
	AllowLan           bool               `yaml:"allow-lan,omitempty"`
	Mode               string             `yaml:"mode,omitempty"`
	LogLevel           string             `yaml:"log-level,omitempty"`
	ExternalController string             `yaml:"external-controller,omitempty"`
	Secret             string             `yaml:"secret,omitempty"`
	DNS                ClashDNS           `yaml:"dns,omitempty"`
	Proxies            []ClashProxies     `yaml:"proxies,omitempty"`
	ProxyGroups        []ClashProxyGroups `yaml:"proxy-groups,omitempty"`
	Rules              []string           `yaml:"rules,omitempty"`
}

type ClashFallbackFilter struct {
	Geoip  bool     `yaml:"geoip,omitempty"`
	Ipcidr []string `yaml:"ipcidr,omitempty"`
	Domain []string `yaml:"domain,omitempty"`
}

type ClashDNS struct {
	Enable         bool                `yaml:"enable,omitempty"`
	Ipv6           bool                `yaml:"ipv6,omitempty"`
	Listen         string              `yaml:"listen,omitempty"`
	EnhancedMode   string              `yaml:"enhanced-mode,omitempty"`
	FakeIPFilter   []string            `yaml:"fake-ip-filter,omitempty"`
	Nameserver     []string            `yaml:"nameserver,omitempty"`
	Fallback       []string            `yaml:"fallback,omitempty"`
	FallbackFilter ClashFallbackFilter `yaml:"fallback-filter,omitempty"`
}

type ClashProxies struct {
	Name           string `yaml:"name"`
	Type           string `yaml:"type"`
	Server         string `yaml:"server"`
	Port           int    `yaml:"port"`
	UUID           string `yaml:"uuid"`
	AlterId        int    `yaml:"alterId"`
	TLS            bool   `yaml:"tls"`
	SkipCertVerify bool   `yaml:"skip-cert-verify,omitempty"`
	NetWork        string `yaml:"network,omitempty"`
	WsOpts         struct {
		Path string `yaml:"path,omitempty"`
	} `yaml:"ws-opts"`
	Cipher        string `yaml:"cipher,omitempty"`
	Password      string `yaml:"password,omitempty"`
	Protocol      string `yaml:"protocol,omitempty"`
	ProtocolParam string `yaml:"protocol-param,omitempty"`
	Obfs          string `yaml:"obfs,omitempty"`
	ObfsParam     string `yaml:"obfs-param,omitempty"`
	UDP           bool   `yaml:"udp,omitempty"`
}

type ClashProxyGroups struct {
	Name     string   `yaml:"name,omitempty"`
	Type     string   `yaml:"type,omitempty"`
	Proxies  []string `yaml:"proxies,omitempty"`
	URL      string   `yaml:"url,omitempty,omitempty"`
	Interval string   `yaml:"interval,omitempty,omitempty"`
}

var BindingRules = []string{
	// (Global_TV)
	// > ABC
	"DOMAIN-SUFFIX,edgedatg.com,Proxies",
	"DOMAIN-SUFFIX,go.com,Proxies",

	// > AbemaTV
	"DOMAIN,linear-abematv.akamaized.net,Proxies",
	"DOMAIN-SUFFIX,abema.io,Proxies",
	"DOMAIN-SUFFIX,abema.tv,Proxies",
	"DOMAIN-SUFFIX,akamaized.net,Proxies",
	"DOMAIN-SUFFIX,ameba.jp,Proxies",
	"DOMAIN-SUFFIX,hayabusa.io,Proxies",

	// > Amazon Prime Video
	"DOMAIN-SUFFIX,aiv-cdn.net,Proxies",
	"DOMAIN-SUFFIX,amazonaws.com,Proxies",
	"DOMAIN-SUFFIX,amazonvideo.com,Proxies",
	"DOMAIN-SUFFIX,llnwd.net,Proxies",

	// > Bahamut
	"DOMAIN-SUFFIX,bahamut.com.tw,Proxies",
	"DOMAIN-SUFFIX,gamer.com.tw,Proxies",
	"DOMAIN-SUFFIX,hinet.net,Proxies",

	// > BBC
	"DOMAIN-KEYWORD,bbcfmt,Proxies",
	"DOMAIN-KEYWORD,co.uk,Proxies",
	"DOMAIN-KEYWORD,uk-live,Proxies",
	"DOMAIN-SUFFIX,bbc.co,Proxies",
	"DOMAIN-SUFFIX,bbc.co.uk,Proxies",
	"DOMAIN-SUFFIX,bbc.com,Proxies",
	"DOMAIN-SUFFIX,bbci.co,Proxies",
	"DOMAIN-SUFFIX,bbci.co.uk,Proxies",

	// > CHOCO TV
	"DOMAIN-SUFFIX,chocotv.com.tw,Proxies",

	// > Epicgames
	"DOMAIN-KEYWORD,epicgames,Proxies",
	"DOMAIN-SUFFIX,helpshift.com,Proxies",

	// > Fox+
	"DOMAIN-KEYWORD,foxplus,Proxies",
	"DOMAIN-SUFFIX,config.fox.com,Proxies",
	"DOMAIN-SUFFIX,emome.net,Proxies",
	"DOMAIN-SUFFIX,fox.com,Proxies",
	"DOMAIN-SUFFIX,foxdcg.com,Proxies",
	"DOMAIN-SUFFIX,foxnow.com,Proxies",
	"DOMAIN-SUFFIX,foxplus.com,Proxies",
	"DOMAIN-SUFFIX,foxplay.com,Proxies",
	"DOMAIN-SUFFIX,ipinfo.io,Proxies",
	"DOMAIN-SUFFIX,mstage.io,Proxies",
	"DOMAIN-SUFFIX,now.com,Proxies",
	"DOMAIN-SUFFIX,theplatform.com,Proxies",
	"DOMAIN-SUFFIX,urlload.net,Proxies",

	// > HBO && HBO Go
	"DOMAIN-SUFFIX,execute-api.ap-southeast-1.amazonaws.com,Proxies",
	"DOMAIN-SUFFIX,hbo.com,Proxies",
	"DOMAIN-SUFFIX,hboasia.com,Proxies",
	"DOMAIN-SUFFIX,hbogo.com,Proxies",
	"DOMAIN-SUFFIX,hbogoasia.hk,Proxies",

	// > Hulu
	"DOMAIN-SUFFIX,happyon.jp,Proxies",
	"DOMAIN-SUFFIX,hulu.com,Proxies",
	"DOMAIN-SUFFIX,huluim.com,Proxies",
	"DOMAIN-SUFFIX,hulustream.com,Proxies",

	// > Imkan
	"DOMAIN-SUFFIX,imkan.tv,Proxies",

	// > JOOX
	"DOMAIN-SUFFIX,joox.com,Proxies",

	// > MytvSUPER
	"DOMAIN-KEYWORD,nowtv100,Proxies",
	"DOMAIN-KEYWORD,rthklive,Proxies",
	"DOMAIN-SUFFIX,mytvsuper.com,Proxies",
	"DOMAIN-SUFFIX,tvb.com,Proxies",

	// > Netflix
	"DOMAIN-SUFFIX,netflix.com,Proxies",
	"DOMAIN-SUFFIX,netflix.net,Proxies",
	"DOMAIN-SUFFIX,nflxext.com,Proxies",
	"DOMAIN-SUFFIX,nflximg.com,Proxies",
	"DOMAIN-SUFFIX,nflximg.net,Proxies",
	"DOMAIN-SUFFIX,nflxso.net,Proxies",
	"DOMAIN-SUFFIX,nflxvideo.net,Proxies",

	// > Pandora
	"DOMAIN-SUFFIX,pandora.com,Proxies",

	// > Sky GO
	"DOMAIN-SUFFIX,sky.com,Proxies",
	"DOMAIN-SUFFIX,skygo.co.nz,Proxies",

	// > Spotify
	"DOMAIN-KEYWORD,spotify,Proxies",
	"DOMAIN-SUFFIX,scdn.co,Proxies",
	"DOMAIN-SUFFIX,spoti.fi,Proxies",

	// > viuTV
	"DOMAIN-SUFFIX,viu.tv,Proxies",

	// > Youtube
	"DOMAIN-KEYWORD,youtube,Proxies",
	"DOMAIN-SUFFIX,googlevideo.com,Proxies",
	"DOMAIN-SUFFIX,gvt2.com,Proxies",
	"DOMAIN-SUFFIX,youtu.be,Proxies",

	// (Asian_TV)
	// > Bilibili
	"DOMAIN-KEYWORD,bilibili,Direct",
	"DOMAIN-SUFFIX,acg.tv,Direct",
	"DOMAIN-SUFFIX,acgvideo.com,Direct",
	"DOMAIN-SUFFIX,b23.tv,Direct",
	"DOMAIN-SUFFIX,biliapi.com,Direct",
	"DOMAIN-SUFFIX,biliapi.net,Direct",
	"DOMAIN-SUFFIX,bilibili.com,Direct",
	"DOMAIN-SUFFIX,biligame.com,Direct",
	"DOMAIN-SUFFIX,biligame.net,Direct",
	"DOMAIN-SUFFIX,hdslb.com,Direct",
	"DOMAIN-SUFFIX,im9.com,Direct",

	// > IQIYI
	"DOMAIN-KEYWORD,qiyi,Direct",
	"DOMAIN-SUFFIX,qy.net,Direct",

	// > letv
	"DOMAIN-SUFFIX,api.mob.app.letv.com,Direct",

	// > NeteaseMusic
	"DOMAIN-SUFFIX,163yun.com,Direct",
	"DOMAIN-SUFFIX,music.126.net,Direct",
	"DOMAIN-SUFFIX,music.163.com,Direct",

	// > Tencent Video
	"DOMAIN-SUFFIX,vv.video.qq.com,Direct",

	// China Area Network
	// > Microsoft
	"DOMAIN-SUFFIX,microsoft.com,Direct",
	"DOMAIN-SUFFIX,windows.net,Direct",
	"DOMAIN-SUFFIX,sfx.ms,Direct",
	"DOMAIN-SUFFIX,sharepoint.com,Direct",
	"DOMAIN-KEYWORD,officecdn,Direct",
	// > Blizzard
	"DOMAIN-SUFFIX,blizzard.com,Direct",
	"DOMAIN-SUFFIX,battle.net,Direct",
	"DOMAIN,blzddist1-a.akamaihd.net,Direct",
	// > Tencent
	// "USER-AGENT,MicroMessenger%20Client,Direct",
	// "USER-AGENT,WeChat*,Direct",
	"DOMAIN-SUFFIX,qq.com,Direct",
	"DOMAIN-SUFFIX,qpic.cn,Direct",
	"DOMAIN-SUFFIX,tencent.com,Direct",
	// > Alibaba
	"DOMAIN-SUFFIX,alibaba.com,Direct",
	"DOMAIN-SUFFIX,alicdn.com,Direct",
	"DOMAIN-SUFFIX,amap.com,Direct",
	"DOMAIN-SUFFIX,dingtalk.com,Direct",
	"DOMAIN-SUFFIX,taobao.com,Direct",
	"DOMAIN-SUFFIX,tmall.com,Direct",
	"DOMAIN-SUFFIX,ykimg.com,Direct",
	"DOMAIN-SUFFIX,youku.com,Direct",
	"DOMAIN-SUFFIX,xiami.com,Direct",
	"DOMAIN-SUFFIX,xiami.net,Direct",
	// > NetEase
	"DOMAIN-SUFFIX,163.com,Direct",
	"DOMAIN-SUFFIX,126.net,Direct",
	"DOMAIN-SUFFIX,163yun.com,Direct",
	// > Sohu
	"DOMAIN-SUFFIX,sohu.com.cn,Direct",
	"DOMAIN-SUFFIX,itc.cn,Direct",
	"DOMAIN-SUFFIX,sohu.com,Direct",
	"DOMAIN-SUFFIX,v-56.com,Direct",
	// > Sina
	"DOMAIN-SUFFIX,weibo.com,Direct",
	"DOMAIN-SUFFIX,weibo.cn,Direct",
	// > JD
	"DOMAIN-SUFFIX,jd.com,Direct",
	"DOMAIN-SUFFIX,jd.hk,Direct",
	"DOMAIN-SUFFIX,360buyimg.com,Direct",
	// > MI
	"DOMAIN-SUFFIX,duokan.com,Direct",
	"DOMAIN-SUFFIX,mi-img.com,Direct",
	"DOMAIN-SUFFIX,mifile.cn,Direct",
	"DOMAIN-SUFFIX,xiaomi.com,Direct",
	// > bilibili
	"DOMAIN-SUFFIX,acgvideo.com,Direct",
	"DOMAIN-SUFFIX,bilibili.com,Direct",
	"DOMAIN-SUFFIX,hdslb.com,Direct",
	// > iQiyi
	"DOMAIN-SUFFIX,iqiyi.com,Direct",
	"DOMAIN-SUFFIX,iqiyipic.com,Direct",
	"DOMAIN-SUFFIX,71.am.com,Direct",
	// > HunanTV
	"DOMAIN-SUFFIX,hitv.com,Direct",
	"DOMAIN-SUFFIX,mgtv.com,Direct",
	// > Meitu
	"DOMAIN-SUFFIX,meitu.com,Direct",
	"DOMAIN-SUFFIX,meitudata.com,Direct",
	"DOMAIN-SUFFIX,meipai.com,Direct",
	// > YYeTs
	"DOMAIN-SUFFIX,zmzapi.com,Direct",
	"DOMAIN-SUFFIX,zimuzu.tv,Direct",
	"DOMAIN-SUFFIX,zmzfile.com,Direct",
	"DOMAIN-SUFFIX,zmzapi.net,Direct",
	// > 蛋蛋赞
	"DOMAIN-SUFFIX,baduziyuan.com,Direct",
	"DOMAIN-SUFFIX,com-hs-hkdy.com,Direct",
	"DOMAIN-SUFFIX,czybjz.com,Direct",
	"DOMAIN-SUFFIX,dandanzan.com,Direct",
	"DOMAIN-SUFFIX,fjhps.com,Direct",
	"DOMAIN-SUFFIX,kuyunbo.club,Direct",
	// > Baidu
	"DOMAIN-SUFFIX,baidu.com,Direct",
	"DOMAIN-SUFFIX,baidubcr.com,Direct",
	"DOMAIN-SUFFIX,bdstatic.com,Direct",
	// > ChinaNet
	"DOMAIN-SUFFIX,189.cn,Direct",
	"DOMAIN-SUFFIX,21cn.com,Direct",
	// > ByteDance
	"DOMAIN-SUFFIX,bytecdn.cn,Direct",
	"DOMAIN-SUFFIX,pstatp.com,Direct",
	"DOMAIN-SUFFIX,snssdk.com,Direct",
	"DOMAIN-SUFFIX,toutiao.com,Direct",
	// > Content Delivery Network
	// > Akamai
	"DOMAIN-SUFFIX,akadns.net,Direct",
	// - DOMAIN-SUFFIX,akamai.net,Direct",
	// - DOMAIN-SUFFIX,akamaiedge.net,Direct",
	// - DOMAIN-SUFFIX,akamaihd.net,Direct",
	// - DOMAIN-SUFFIX,akamaistream.net,Direct",
	// - DOMAIN-SUFFIX,akamaized.net,Direct",
	// > ChinaNetCenter
	"DOMAIN-SUFFIX,chinanetcenter.com,Direct",
	"DOMAIN-SUFFIX,wangsu.com,Direct",
	// > IP Query
	"DOMAIN-SUFFIX,ipip.net,Direct",
	"DOMAIN-SUFFIX,ip.cn,Direct",
	"DOMAIN-SUFFIX,ip.la,Direct",
	"DOMAIN-SUFFIX,ip-cdn.com,Direct",
	"DOMAIN-SUFFIX,ipv6-test.com,Direct",
	"DOMAIN-SUFFIX,test-ipv6.com,Direct",
	"DOMAIN-SUFFIX,whatismyip.com,Direct",
	"DOMAIN,ip.bjango.com,Direct",
	// > Other
	"DOMAIN-SUFFIX,40017.cn,Direct",
	"DOMAIN-SUFFIX,broadcasthe.net,Direct",
	"DOMAIN-SUFFIX,cailianpress.com,Direct",
	"DOMAIN-SUFFIX,chdbits.co,Direct",
	"DOMAIN-SUFFIX,chushou.tv,Direct",
	"DOMAIN-SUFFIX,cmbchina.com,Direct",
	"DOMAIN-SUFFIX,cmbimg.com,Direct",
	"DOMAIN-SUFFIX,cmct.tv,Direct",
	"DOMAIN-SUFFIX,cmvideo.cn,Direct",
	"DOMAIN-SUFFIX,cnlang.org,Direct",
	"DOMAIN-SUFFIX,doubanio.com,Direct",
	"DOMAIN-SUFFIX,douyu.com,Direct",
	"DOMAIN-SUFFIX,douyucdn.cn,Direct",
	"DOMAIN-SUFFIX,dxycdn.com,Direct",
	"DOMAIN-SUFFIX,hicloud.com,Direct",
	"DOMAIN-SUFFIX,hdchina.org,Direct",
	"DOMAIN-SUFFIX,hdcmct.org,Direct",
	"DOMAIN-SUFFIX,ithome.com,Direct",
	"DOMAIN-SUFFIX,kkmh.com,Direct",
	"DOMAIN-SUFFIX,ksosoft.com,Direct",
	"DOMAIN-SUFFIX,maoyun.tv,Direct",
	"DOMAIN-SUFFIX,meituan.net,Direct",
	"DOMAIN-SUFFIX,mobike.com,Direct",
	"DOMAIN-SUFFIX,mubu.com,Direct",
	"DOMAIN-SUFFIX,myzaker.com,Direct",
	"DOMAIN-SUFFIX,ourbits.club,Direct",
	"DOMAIN-SUFFIX,passthepopcorn.me,Direct",
	"DOMAIN-SUFFIX,paypal.com,Direct",
	"DOMAIN-SUFFIX,paypalobjects.com,Direct",
	"DOMAIN-SUFFIX,privatehd.to,Direct",
	"DOMAIN-SUFFIX,redacted.ch,Direct",
	"DOMAIN-SUFFIX,ruguoapp.com,Direct",
	"DOMAIN-SUFFIX,smzdm.com,Direct",
	"DOMAIN-SUFFIX,sogou.com,Direct",
	"DOMAIN-SUFFIX,teamviewer.com,Direct",
	"DOMAIN-SUFFIX,totheglory.im,Direct",
	"DOMAIN-SUFFIX,udacity.com,Direct",
	"DOMAIN-SUFFIX,xmcdn.com,Direct",
	"DOMAIN-SUFFIX,yangkeduo.com,Direct",
	"DOMAIN-SUFFIX,zhihu.com,Direct",
	"DOMAIN-SUFFIX,zhimg.com,Direct",
	// "USER-AGENT,NeteaseMusic*,Direct",
	// "USER-AGENT,%E7%BD%91%E6%98%93%E4%BA%91%E9%9F%B3%E4%B9%90*,Direct",

	// (DNS Cache Pollution Protection)
	// > Google
	"DOMAIN-SUFFIX,appspot.com,Proxies",
	"DOMAIN-SUFFIX,blogger.com,Proxies",
	"DOMAIN-SUFFIX,getoutline.org,Proxies",
	"DOMAIN-SUFFIX,gvt0.com,Proxies",
	"DOMAIN-SUFFIX,gvt1.com,Proxies",
	"DOMAIN-SUFFIX,gvt3.com,Proxies",
	"DOMAIN-SUFFIX,xn--ngstr-lra8j.com,Proxies",
	"DOMAIN-KEYWORD,google,Proxies",
	"DOMAIN-KEYWORD,blogspot,Proxies",
	// > Facebook
	"DOMAIN-SUFFIX,cdninstagram.com,Proxies",
	"DOMAIN-SUFFIX,fb.com,Proxies",
	"DOMAIN-SUFFIX,fb.me,Proxies",
	"DOMAIN-SUFFIX,fbaddins.com,Proxies",
	"DOMAIN-SUFFIX,fbcdn.net,Proxies",
	"DOMAIN-SUFFIX,fbsbx.com,Proxies",
	"DOMAIN-SUFFIX,fbworkmail.com,Proxies",
	"DOMAIN-SUFFIX,instagram.com,Proxies",
	"DOMAIN-SUFFIX,m.me,Proxies",
	"DOMAIN-SUFFIX,messenger.com,Proxies",
	"DOMAIN-SUFFIX,oculus.com,Proxies",
	"DOMAIN-SUFFIX,oculuscdn.com,Proxies",
	"DOMAIN-SUFFIX,rocksdb.org,Proxies",
	"DOMAIN-SUFFIX,whatsapp.com,Proxies",
	"DOMAIN-SUFFIX,whatsapp.net,Proxies",
	"DOMAIN-KEYWORD,facebook,Proxies",
	// > Twitter
	"DOMAIN-SUFFIX,pscp.tv,Proxies",
	"DOMAIN-SUFFIX,periscope.tv,Proxies",
	"DOMAIN-SUFFIX,t.co,Proxies",
	"DOMAIN-SUFFIX,twimg.co,Proxies",
	"DOMAIN-SUFFIX,twimg.com,Proxies",
	"DOMAIN-SUFFIX,twitpic.com,Proxies",
	"DOMAIN-SUFFIX,vine.co,Proxies",
	"DOMAIN-KEYWORD,twitter,Proxies",
	// > Telegram
	"DOMAIN-SUFFIX,t.me,Proxies",
	"DOMAIN-SUFFIX,tdesktop.com,Proxies",
	"DOMAIN-SUFFIX,telegra.ph,Proxies",
	"DOMAIN-SUFFIX,telegram.me,Proxies",
	"DOMAIN-SUFFIX,telegram.org,Proxies",
	// > Line
	"DOMAIN-SUFFIX,line.me,Proxies",
	"DOMAIN-SUFFIX,line-apps.com,Proxies",
	"DOMAIN-SUFFIX,line-scdn.net,Proxies",
	"DOMAIN-SUFFIX,naver.jp,Proxies",
	// > Other
	"DOMAIN-SUFFIX,4shared.com,Proxies",
	"DOMAIN-SUFFIX,881903.com,Proxies",
	"DOMAIN-SUFFIX,abc.net.au,Proxies",
	"DOMAIN-SUFFIX,abebooks.com,Proxies",
	"DOMAIN-SUFFIX,amazon.co.jp,Proxies",
	"DOMAIN-SUFFIX,apigee.com,Proxies",
	"DOMAIN-SUFFIX,apk-dl.com,Proxies",
	"DOMAIN-SUFFIX,apkmirror.com,Proxies",
	"DOMAIN-SUFFIX,apkmonk.com,Proxies",
	"DOMAIN-SUFFIX,apkpure.com,Proxies",
	"DOMAIN-SUFFIX,aptoide.com,Proxies",
	"DOMAIN-SUFFIX,archive.is,Proxies",
	"DOMAIN-SUFFIX,archive.org,Proxies",
	"DOMAIN-SUFFIX,arte.tv,Proxies",
	"DOMAIN-SUFFIX,ask.com,Proxies",
	"DOMAIN-SUFFIX,avgle.com,Proxies",
	"DOMAIN-SUFFIX,badoo.com,Proxies",
	"DOMAIN-SUFFIX,bandwagonhost.com,Proxies",
	"DOMAIN-SUFFIX,bbc.com,Proxies",
	"DOMAIN-SUFFIX,behance.net,Proxies",
	"DOMAIN-SUFFIX,bibox.com,Proxies",
	"DOMAIN-SUFFIX,biggo.com.tw,Proxies",
	"DOMAIN-SUFFIX,binance.com,Proxies",
	"DOMAIN-SUFFIX,bitcointalk.org,Proxies",
	"DOMAIN-SUFFIX,bitfinex.com,Proxies",
	"DOMAIN-SUFFIX,bitmex.com,Proxies",
	"DOMAIN-SUFFIX,bit-z.com,Proxies",
	"DOMAIN-SUFFIX,bloglovin.com,Proxies",
	"DOMAIN-SUFFIX,bloomberg.cn,Proxies",
	"DOMAIN-SUFFIX,bloomberg.com,Proxies",
	"DOMAIN-SUFFIX,book.com.tw,Proxies",
	"DOMAIN-SUFFIX,booklive.jp,Proxies",
	"DOMAIN-SUFFIX,books.com.tw,Proxies",
	"DOMAIN-SUFFIX,box.com,Proxies",
	"DOMAIN-SUFFIX,brookings.edu,Proxies",
	"DOMAIN-SUFFIX,businessinsider.com,Proxies",
	"DOMAIN-SUFFIX,bwh1.net,Proxies",
	"DOMAIN-SUFFIX,castbox.fm,Proxies",
	"DOMAIN-SUFFIX,cbc.ca,Proxies",
	"DOMAIN-SUFFIX,cdw.com,Proxies",
	"DOMAIN-SUFFIX,change.org,Proxies",
	"DOMAIN-SUFFIX,ck101.com,Proxies",
	"DOMAIN-SUFFIX,clarionproject.org,Proxies",
	"DOMAIN-SUFFIX,clyp.it,Proxies",
	"DOMAIN-SUFFIX,cna.com.tw,Proxies",
	"DOMAIN-SUFFIX,comparitech.com,Proxies",
	"DOMAIN-SUFFIX,conoha.jp,Proxies",
	"DOMAIN-SUFFIX,crucial.com,Proxies",
	"DOMAIN-SUFFIX,cts.com.tw,Proxies",
	"DOMAIN-SUFFIX,cw.com.tw,Proxies",
	"DOMAIN-SUFFIX,cyberctm.com,Proxies",
	"DOMAIN-SUFFIX,dailymotion.com,Proxies",
	"DOMAIN-SUFFIX,dailyview.tw,Proxies",
	"DOMAIN-SUFFIX,daum.net,Proxies",
	"DOMAIN-SUFFIX,daumcdn.net,Proxies",
	"DOMAIN-SUFFIX,dcard.tw,Proxies",
	"DOMAIN-SUFFIX,deepdiscount.com,Proxies",
	"DOMAIN-SUFFIX,deezer.com,Proxies",
	"DOMAIN-SUFFIX,depositphotos.com,Proxies",
	"DOMAIN-SUFFIX,disconnect.me,Proxies",
	"DOMAIN-SUFFIX,discordapp.com,Proxies",
	"DOMAIN-SUFFIX,discordapp.net,Proxies",
	"DOMAIN-SUFFIX,disqus.com,Proxies",
	"DOMAIN-SUFFIX,dns2go.com,Proxies",
	"DOMAIN-SUFFIX,dropbox.com,Proxies",
	"DOMAIN-SUFFIX,dropboxusercontent.com,Proxies",
	"DOMAIN-SUFFIX,duckduckgo.com,Proxies",
	"DOMAIN-SUFFIX,dw.com,Proxies",
	"DOMAIN-SUFFIX,dynu.com,Proxies",
	"DOMAIN-SUFFIX,earthcam.com,Proxies",
	"DOMAIN-SUFFIX,ebookservice.tw,Proxies",
	"DOMAIN-SUFFIX,economist.com,Proxies",
	"DOMAIN-SUFFIX,edgecastcdn.net,Proxies",
	"DOMAIN-SUFFIX,edu,Proxies",
	"DOMAIN-SUFFIX,elpais.com,Proxies",
	"DOMAIN-SUFFIX,enanyang.my,Proxies",
	"DOMAIN-SUFFIX,euronews.com,Proxies",
	"DOMAIN-SUFFIX,feedly.com,Proxies",
	"DOMAIN-SUFFIX,files.wordpress.com,Proxies",
	"DOMAIN-SUFFIX,flickr.com,Proxies",
	"DOMAIN-SUFFIX,flitto.com,Proxies",
	"DOMAIN-SUFFIX,foreignpolicy.com,Proxies",
	"DOMAIN-SUFFIX,friday.tw,Proxies",
	"DOMAIN-SUFFIX,gate.io,Proxies",
	"DOMAIN-SUFFIX,getlantern.org,Proxies",
	"DOMAIN-SUFFIX,getsync.com,Proxies",
	"DOMAIN-SUFFIX,goo.ne.jp,Proxies",
	"DOMAIN-SUFFIX,goodreads.com,Proxies",
	"DOMAIN-SUFFIX,gov.tw,Proxies",
	"DOMAIN-SUFFIX,gumroad.com,Proxies",
	"DOMAIN-SUFFIX,hbg.com,Proxies",
	"DOMAIN-SUFFIX,hightail.com,Proxies",
	"DOMAIN-SUFFIX,hk01.com,Proxies",
	"DOMAIN-SUFFIX,hkbf.org,Proxies",
	"DOMAIN-SUFFIX,hkbookcity.com,Proxies",
	"DOMAIN-SUFFIX,hkej.com,Proxies",
	"DOMAIN-SUFFIX,hket.com,Proxies",
	"DOMAIN-SUFFIX,hkgolden.com,Proxies",
	"DOMAIN-SUFFIX,hootsuite.com,Proxies",
	"DOMAIN-SUFFIX,hudson.org,Proxies",
	"DOMAIN-SUFFIX,huobi.pro,Proxies",
	"DOMAIN-SUFFIX,initiummall.com,Proxies",
	"DOMAIN-SUFFIX,ipfs.io,Proxies",
	"DOMAIN-SUFFIX,issuu.com,Proxies",
	"DOMAIN-SUFFIX,japantimes.co.jp,Proxies",
	"DOMAIN-SUFFIX,jiji.com,Proxies",
	"DOMAIN-SUFFIX,jinx.com,Proxies",
	"DOMAIN-SUFFIX,jkforum.net,Proxies",
	"DOMAIN-SUFFIX,joinmastodon.org,Proxies",
	"DOMAIN-SUFFIX,kakao.com,Proxies",
	"DOMAIN-SUFFIX,lihkg.com,Proxies",
	"DOMAIN-SUFFIX,live.com,Proxies",
	"DOMAIN-SUFFIX,mail.ru,Proxies",
	"DOMAIN-SUFFIX,matters.news,Proxies",
	"DOMAIN-SUFFIX,medium.com,Proxies",
	"DOMAIN-SUFFIX,mega.nz,Proxies",
	"DOMAIN-SUFFIX,mil,Proxies",
	"DOMAIN-SUFFIX,mobile01.com,Proxies",
	"DOMAIN-SUFFIX,naver.com,Proxies",
	"DOMAIN-SUFFIX,nikkei.com,Proxies",
	"DOMAIN-SUFFIX,nofile.io,Proxies",
	"DOMAIN-SUFFIX,now.com,Proxies",
	"DOMAIN-SUFFIX,nyt.com,Proxies",
	"DOMAIN-SUFFIX,nytchina.com,Proxies",
	"DOMAIN-SUFFIX,nytcn.me,Proxies",
	"DOMAIN-SUFFIX,nytco.com,Proxies",
	"DOMAIN-SUFFIX,nytimes.com,Proxies",
	"DOMAIN-SUFFIX,nytimg.com,Proxies",
	"DOMAIN-SUFFIX,nytlog.com,Proxies",
	"DOMAIN-SUFFIX,nytstyle.com,Proxies",
	"DOMAIN-SUFFIX,ok.ru,Proxies",
	"DOMAIN-SUFFIX,okex.com,Proxies",
	"DOMAIN-SUFFIX,pcloud.com,Proxies",
	"DOMAIN-SUFFIX,pinimg.com,Proxies",
	"DOMAIN-SUFFIX,pixiv.net,Proxies",
	"DOMAIN-SUFFIX,pornhub.com,Proxies",
	"DOMAIN-SUFFIX,pureapk.com,Proxies",
	"DOMAIN-SUFFIX,quora.com,Proxies",
	"DOMAIN-SUFFIX,quoracdn.net,Proxies",
	"DOMAIN-SUFFIX,rakuten.co.jp,Proxies",
	"DOMAIN-SUFFIX,reddit.com,Proxies",
	"DOMAIN-SUFFIX,redditmedia.com,Proxies",
	"DOMAIN-SUFFIX,resilio.com,Proxies",
	"DOMAIN-SUFFIX,reuters.com,Proxies",
	"DOMAIN-SUFFIX,scmp.com,Proxies",
	"DOMAIN-SUFFIX,scribd.com,Proxies",
	"DOMAIN-SUFFIX,seatguru.com,Proxies",
	"DOMAIN-SUFFIX,shadowsocks.org,Proxies",
	"DOMAIN-SUFFIX,slideshare.net,Proxies",
	"DOMAIN-SUFFIX,soundcloud.com,Proxies",
	"DOMAIN-SUFFIX,startpage.com,Proxies",
	"DOMAIN-SUFFIX,steemit.com,Proxies",
	"DOMAIN-SUFFIX,t66y.com,Proxies",
	"DOMAIN-SUFFIX,teco-hk.org,Proxies",
	"DOMAIN-SUFFIX,teco-mo.org,Proxies",
	"DOMAIN-SUFFIX,teddysun.com,Proxies",
	"DOMAIN-SUFFIX,theinitium.com,Proxies",
	"DOMAIN-SUFFIX,tineye.com,Proxies",
	"DOMAIN-SUFFIX,torproject.org,Proxies",
	"DOMAIN-SUFFIX,tumblr.com,Proxies",
	"DOMAIN-SUFFIX,turbobit.net,Proxies",
	"DOMAIN-SUFFIX,twitch.tv,Proxies",
	"DOMAIN-SUFFIX,udn.com,Proxies",
	"DOMAIN-SUFFIX,unseen.is,Proxies",
	"DOMAIN-SUFFIX,upmedia.mg,Proxies",
	"DOMAIN-SUFFIX,uptodown.com,Proxies",
	"DOMAIN-SUFFIX,ustream.tv,Proxies",
	"DOMAIN-SUFFIX,uwants.com,Proxies",
	"DOMAIN-SUFFIX,v2ray.com,Proxies",
	"DOMAIN-SUFFIX,viber.com,Proxies",
	"DOMAIN-SUFFIX,videopress.com,Proxies",
	"DOMAIN-SUFFIX,vimeo.com,Proxies",
	"DOMAIN-SUFFIX,voxer.com,Proxies",
	"DOMAIN-SUFFIX,vzw.com,Proxies",
	"DOMAIN-SUFFIX,w3schools.com,Proxies",
	"DOMAIN-SUFFIX,wattpad.com,Proxies",
	"DOMAIN-SUFFIX,whoer.net,Proxies",
	"DOMAIN-SUFFIX,wikimapia.org,Proxies",
	"DOMAIN-SUFFIX,wikipedia.org,Proxies",
	"DOMAIN-SUFFIX,wire.com,Proxies",
	"DOMAIN-SUFFIX,worldcat.org,Proxies",
	"DOMAIN-SUFFIX,wsj.com,Proxies",
	"DOMAIN-SUFFIX,wsj.net,Proxies",
	"DOMAIN-SUFFIX,xboxlive.com,Proxies",
	"DOMAIN-SUFFIX,xvideos.com,Proxies",
	"DOMAIN-SUFFIX,yahoo.com,Proxies",
	"DOMAIN-SUFFIX,yesasia.com,Proxies",
	"DOMAIN-SUFFIX,yes-news.com,Proxies",
	"DOMAIN-SUFFIX,yomiuri.co.jp,Proxies",
	"DOMAIN-SUFFIX,you-get.org,Proxies",
	"DOMAIN-SUFFIX,zb.com,Proxies",
	"DOMAIN-SUFFIX,zello.com,Proxies",
	"DOMAIN-SUFFIX,zeronet.io,Proxies",
	"DOMAIN,cdn-images.mailchimp.com,Proxies",
	"DOMAIN,id.heroku.com,Proxies",
	"DOMAIN-KEYWORD,github,Proxies",
	"DOMAIN-KEYWORD,jav,Proxies",
	"DOMAIN-KEYWORD,pinterest,Proxies",
	"DOMAIN-KEYWORD,porn,Proxies",
	"DOMAIN-KEYWORD,wikileaks,Proxies",

	// (Region-Restricted Access Denied)
	"DOMAIN-SUFFIX,apartmentratings.com,Proxies",
	"DOMAIN-SUFFIX,apartments.com,Proxies",
	"DOMAIN-SUFFIX,bankmobilevibe.com,Proxies",
	"DOMAIN-SUFFIX,bing.com,Proxies",
	"DOMAIN-SUFFIX,booktopia.com.au,Proxies",
	"DOMAIN-SUFFIX,centauro.com.br,Proxies",
	"DOMAIN-SUFFIX,clearsurance.com,Proxies",
	"DOMAIN-SUFFIX,costco.com,Proxies",
	"DOMAIN-SUFFIX,crackle.com,Proxies",
	"DOMAIN-SUFFIX,depositphotos.cn,Proxies",
	"DOMAIN-SUFFIX,dish.com,Proxies",
	"DOMAIN-SUFFIX,dmm.co.jp,Proxies",
	"DOMAIN-SUFFIX,dmm.com,Proxies",
	"DOMAIN-SUFFIX,dnvod.tv,Proxies",
	"DOMAIN-SUFFIX,esurance.com,Proxies",
	"DOMAIN-SUFFIX,extmatrix.com,Proxies",
	"DOMAIN-SUFFIX,fastpic.ru,Proxies",
	"DOMAIN-SUFFIX,flipboard.com,Proxies",
	"DOMAIN-SUFFIX,fnac.be,Proxies",
	"DOMAIN-SUFFIX,fnac.com,Proxies",
	"DOMAIN-SUFFIX,funkyimg.com,Proxies",
	"DOMAIN-SUFFIX,fxnetworks.com,Proxies",
	"DOMAIN-SUFFIX,gettyimages.com,Proxies",
	"DOMAIN-SUFFIX,jcpenney.com,Proxies",
	"DOMAIN-SUFFIX,kknews.cc,Proxies",
	"DOMAIN-SUFFIX,nationwide.com,Proxies",
	"DOMAIN-SUFFIX,nbc.com,Proxies",
	"DOMAIN-SUFFIX,nordstrom.com,Proxies",
	"DOMAIN-SUFFIX,nordstromimage.com,Proxies",
	"DOMAIN-SUFFIX,nordstromrack.com,Proxies",
	"DOMAIN-SUFFIX,read01.com,Proxies",
	"DOMAIN-SUFFIX,superpages.com,Proxies",
	"DOMAIN-SUFFIX,target.com,Proxies",
	"DOMAIN-SUFFIX,thinkgeek.com,Proxies",
	"DOMAIN-SUFFIX,tracfone.com,Proxies",
	"DOMAIN-SUFFIX,uploader.jp,Proxies",
	"DOMAIN-SUFFIX,vevo.com,Proxies",
	"DOMAIN-SUFFIX,viu.tv,Proxies",
	"DOMAIN-SUFFIX,vk.com,Proxies",
	"DOMAIN-SUFFIX,vsco.co,Proxies",
	"DOMAIN-SUFFIX,xfinity.com,Proxies",
	"DOMAIN-SUFFIX,zattoo.com,Proxies",
	"DOMAIN,abc.com,Proxies",
	"DOMAIN,abc.go.com,Proxies",
	"DOMAIN,abc.net.au,Proxies",
	"DOMAIN,wego.here.com,Proxies",
	// "USER-AGENT,Roam*,Proxies",

	// (The Most Popular Sites)
	// > Apple
	// > Apple URL Shortener
	"DOMAIN-SUFFIX,appsto.re,Proxies",
	// > TestFlight
	"DOMAIN,beta.itunes.apple.com,Proxies",
	// > iBooks Store download
	"DOMAIN,books.itunes.apple.com,Proxies",
	// > iTunes Store Moveis Trailers
	"DOMAIN,hls.itunes.apple.com,Proxies",
	// App Store Preview
	"DOMAIN,itunes.apple.com,Proxies",
	// > Spotlight
	"DOMAIN,api-glb-sea.smoot.apple.com,Proxies",
	// > Dictionary
	"DOMAIN,lookup-api.apple.com,Proxies",
	"PROCESS-NAME,LookupViewService,Proxies",
	// > Google
	"DOMAIN-SUFFIX,abc.xyz,Proxies",
	"DOMAIN-SUFFIX,android.com,Proxies",
	"DOMAIN-SUFFIX,androidify.com,Proxies",
	"DOMAIN-SUFFIX,dialogflow.com,Proxies",
	"DOMAIN-SUFFIX,autodraw.com,Proxies",
	"DOMAIN-SUFFIX,capitalg.com,Proxies",
	"DOMAIN-SUFFIX,certificate-transparency.org,Proxies",
	"DOMAIN-SUFFIX,chrome.com,Proxies",
	"DOMAIN-SUFFIX,chromeexperiments.com,Proxies",
	"DOMAIN-SUFFIX,chromestatus.com,Proxies",
	"DOMAIN-SUFFIX,chromium.org,Proxies",
	"DOMAIN-SUFFIX,creativelab5.com,Proxies",
	"DOMAIN-SUFFIX,debug.com,Proxies",
	"DOMAIN-SUFFIX,deepmind.com,Proxies",
	"DOMAIN-SUFFIX,firebaseio.com,Proxies",
	"DOMAIN-SUFFIX,getmdl.io,Proxies",
	"DOMAIN-SUFFIX,ggpht.com,Proxies",
	"DOMAIN-SUFFIX,gmail.com,Proxies",
	"DOMAIN-SUFFIX,gmodules.com,Proxies",
	"DOMAIN-SUFFIX,godoc.org,Proxies",
	"DOMAIN-SUFFIX,golang.org,Proxies",
	"DOMAIN-SUFFIX,gstatic.com,Proxies",
	"DOMAIN-SUFFIX,gv.com,Proxies",
	"DOMAIN-SUFFIX,gwtproject.org,Proxies",
	"DOMAIN-SUFFIX,itasoftware.com,Proxies",
	"DOMAIN-SUFFIX,madewithcode.com,Proxies",
	"DOMAIN-SUFFIX,material.io,Proxies",
	"DOMAIN-SUFFIX,polymer-project.org,Proxies",
	"DOMAIN-SUFFIX,admin.recaptcha.net,Proxies",
	"DOMAIN-SUFFIX,recaptcha.net,Proxies",
	"DOMAIN-SUFFIX,shattered.io,Proxies",
	"DOMAIN-SUFFIX,synergyse.com,Proxies",
	"DOMAIN-SUFFIX,tensorflow.org,Proxies",
	"DOMAIN-SUFFIX,tiltbrush.com,Proxies",
	"DOMAIN-SUFFIX,waveprotocol.org,Proxies",
	"DOMAIN-SUFFIX,waymo.com,Proxies",
	"DOMAIN-SUFFIX,webmproject.org,Proxies",
	"DOMAIN-SUFFIX,webrtc.org,Proxies",
	"DOMAIN-SUFFIX,whatbrowser.org,Proxies",
	"DOMAIN-SUFFIX,widevine.com,Proxies",
	"DOMAIN-SUFFIX,x.company,Proxies",
	"DOMAIN-SUFFIX,youtu.be,Proxies",
	"DOMAIN-SUFFIX,yt.be,Proxies",
	"DOMAIN-SUFFIX,ytimg.com,Proxies",
	// > Steam
	"DOMAIN-SUFFIX,steampowered.com,Direct",
	// > Other
	"DOMAIN-SUFFIX,0rz.tw,Proxies",
	"DOMAIN-SUFFIX,4bluestones.biz,Proxies",
	"DOMAIN-SUFFIX,9bis.net,Proxies",
	"DOMAIN-SUFFIX,allconnected.co,Proxies",
	"DOMAIN-SUFFIX,amazonaws.com,Proxies",
	"DOMAIN-SUFFIX,aol.com,Proxies",
	"DOMAIN-SUFFIX,bcc.com.tw,Proxies",
	"DOMAIN-SUFFIX,bit.ly,Proxies",
	"DOMAIN-SUFFIX,bitshare.com,Proxies",
	"DOMAIN-SUFFIX,blog.jp,Proxies",
	"DOMAIN-SUFFIX,blogimg.jp,Proxies",
	"DOMAIN-SUFFIX,blogtd.org,Proxies",
	"DOMAIN-SUFFIX,broadcast.co.nz,Proxies",
	"DOMAIN-SUFFIX,camfrog.com,Proxies",
	"DOMAIN-SUFFIX,cfos.de,Proxies",
	"DOMAIN-SUFFIX,citypopulation.de,Proxies",
	"DOMAIN-SUFFIX,cloudfront.net,Proxies",
	"DOMAIN-SUFFIX,ctitv.com.tw,Proxies",
	"DOMAIN-SUFFIX,cuhk.edu.hk,Proxies",
	"DOMAIN-SUFFIX,cusu.hk,Proxies",
	"DOMAIN-SUFFIX,discuss.com.hk,Proxies",
	"DOMAIN-SUFFIX,dropboxapi.com,Proxies",
	"DOMAIN-SUFFIX,edditstatic.com,Proxies",
	"DOMAIN-SUFFIX,flickriver.com,Proxies",
	"DOMAIN-SUFFIX,focustaiwan.tw,Proxies",
	"DOMAIN-SUFFIX,free.fr,Proxies",
	"DOMAIN-SUFFIX,ftchinese.com,Proxies",
	"DOMAIN-SUFFIX,gigacircle.com,Proxies",
	"DOMAIN-SUFFIX,gov,Proxies",
	"DOMAIN-SUFFIX,hk-pub.com,Proxies",
	"DOMAIN-SUFFIX,hosting.co.uk,Proxies",
	"DOMAIN-SUFFIX,hwcdn.net,Proxies",
	"DOMAIN-SUFFIX,jtvnw.net,Proxies",
	"DOMAIN-SUFFIX,linksalpha.com,Proxies",
	"DOMAIN-SUFFIX,manyvids.com,Proxies",
	"DOMAIN-SUFFIX,myactimes.com,Proxies",
	"DOMAIN-SUFFIX,newsblur.com,Proxies",
	"DOMAIN-SUFFIX,now.im,Proxies",
	"DOMAIN-SUFFIX,redditlist.com,Proxies",
	"DOMAIN-SUFFIX,signal.org,Proxies",
	"DOMAIN-SUFFIX,sparknotes.com,Proxies",
	"DOMAIN-SUFFIX,streetvoice.com,Proxies",
	"DOMAIN-SUFFIX,ttvnw.net,Proxies",
	"DOMAIN-SUFFIX,tv.com,Proxies",
	"DOMAIN-SUFFIX,twitchcdn.net,Proxies",
	"DOMAIN-SUFFIX,typepad.com,Proxies",
	"DOMAIN-SUFFIX,udnbkk.com,Proxies",
	"DOMAIN-SUFFIX,whispersystems.org,Proxies",
	"DOMAIN-SUFFIX,wikia.com,Proxies",
	"DOMAIN-SUFFIX,wn.com,Proxies",
	"DOMAIN-SUFFIX,wolframalpha.com,Proxies",
	"DOMAIN-SUFFIX,x-art.com,Proxies",
	"DOMAIN-SUFFIX,yimg.com,Proxies",

	// Local Area Network
	"DOMAIN-KEYWORD,announce,DIRECT",
	"DOMAIN-KEYWORD,torrent,DIRECT",
	"DOMAIN-KEYWORD,tracker,DIRECT",
	"DOMAIN-SUFFIX,smtp,DIRECT",
	"DOMAIN-SUFFIX,local,DIRECT",
	"IP-CIDR,192.168.0.0/16,DIRECT",
	"IP-CIDR,10.0.0.0/8,DIRECT",
	"IP-CIDR,172.16.0.0/12,DIRECT",
	"IP-CIDR,127.0.0.0/8,DIRECT",
	"IP-CIDR,100.64.0.0/10,DIRECT",

	// > IQIYI
	"IP-CIDR,101.227.0.0/16,Direct",
	"IP-CIDR,101.224.0.0/13,Direct",
	"IP-CIDR,119.176.0.0/12,Direct",

	// > Youku
	"IP-CIDR,106.11.0.0/16,Direct",

	// > Telegram
	"IP-CIDR,67.198.55.0/24,Proxies",
	"IP-CIDR,91.108.4.0/22,Proxies",
	"IP-CIDR,91.108.8.0/22,Proxies",
	"IP-CIDR,91.108.12.0/22,Proxies",
	"IP-CIDR,91.108.16.0/22,Proxies",
	"IP-CIDR,91.108.56.0/22,Proxies",
	"IP-CIDR,109.239.140.0/24,Proxies",
	"IP-CIDR,149.154.160.0/20,Proxies",
	"IP-CIDR,205.172.60.0/22,Proxies",

	// (Extra IP-CIRD)
	// > Google
	"IP-CIDR,35.190.247.0/24,Proxies",
	"IP-CIDR,64.233.160.0/19,Proxies",
	"IP-CIDR,66.102.0.0/20,Proxies",
	"IP-CIDR,66.249.80.0/20,Proxies",
	"IP-CIDR,72.14.192.0/18,Proxies",
	"IP-CIDR,74.125.0.0/16,Proxies",
	"IP-CIDR,108.177.8.0/21,Proxies",
	"IP-CIDR,172.217.0.0/16,Proxies",
	"IP-CIDR,173.194.0.0/16,Proxies",
	"IP-CIDR,209.85.128.0/17,Proxies",
	"IP-CIDR,216.58.192.0/19,Proxies",
	"IP-CIDR,216.239.32.0/19,Proxies",
	// > Facebook
	"IP-CIDR,31.13.24.0/21,Proxies",
	"IP-CIDR,31.13.64.0/18,Proxies",
	"IP-CIDR,45.64.40.0/22,Proxies",
	"IP-CIDR,66.220.144.0/20,Proxies",
	"IP-CIDR,69.63.176.0/20,Proxies",
	"IP-CIDR,69.171.224.0/19,Proxies",
	"IP-CIDR,74.119.76.0/22,Proxies",
	"IP-CIDR,103.4.96.0/22,Proxies",
	"IP-CIDR,129.134.0.0/17,Proxies",
	"IP-CIDR,157.240.0.0/17,Proxies",
	"IP-CIDR,173.252.64.0/19,Proxies",
	"IP-CIDR,173.252.96.0/19,Proxies",
	"IP-CIDR,179.60.192.0/22,Proxies",
	"IP-CIDR,185.60.216.0/22,Proxies",
	"IP-CIDR,204.15.20.0/22,Proxies",
	// > Twitter
	"IP-CIDR,69.195.160.0/19,Proxies",
	"IP-CIDR,104.244.42.0/21,Proxies",
	"IP-CIDR,192.133.76.0/22,Proxies",
	"IP-CIDR,199.16.156.0/22,Proxies",
	"IP-CIDR,199.59.148.0/22,Proxies",
	"IP-CIDR,199.96.56.0/21,Proxies",
	"IP-CIDR,202.160.128.0/22,Proxies",
	"IP-CIDR,209.237.192.0/19,Proxies",

	// GeoIP China
	"GEOIP,CN,Direct",
}
