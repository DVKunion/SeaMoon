package proxy

import (
	"SeaMoon/pkg/ca"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"github.com/elazarl/goproxy"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var CA = "./pkg/ca/ca.pem"
var KEY = "./pkg/ca/ca.key.pem"

// AliYunHttpHandler 阿里云HTTP代理入口
func AliYunHttpHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	host := req.Header.Get("SM-Host")
	// 环路判断问题: 因为代理发出的请求是不带 SM-Host 标识的，所以这个return会直接阻止了环路。
	// 也就是说，最多只会重复请求两次，即代理请求一次自己后，返回了此处的异常。
	if host == "" {
		return errors.New("SM-Host Not Exits")
	}
	req.Header.Del("SM-Host")
	req.Header.Set("Host", host)
	req.URL, _ = url.Parse(host)
	req.Host = req.URL.Host
	req.RequestURI = host
	return handler(w, req, false)
}

func NewHttpClient(listenAddr string, proxyAddr string, verbose bool) {
	log.Println("start http client at " + listenAddr)
	if proxyAddr == "" {
		log.Fatalf("No proxyAddr, Please confirm using -p")
	}
	server := goproxy.NewProxyHttpServer()
	_, errCa := os.Stat(CA)
	_, errKey := os.Stat(KEY)
	if os.IsNotExist(errCa) || os.IsNotExist(errKey) {
		log.Fatalf("Ca Not Exists, Run openssl-gen.sh at ./pkg/ca")
	}
	caCert, err := ioutil.ReadFile(CA)
	if err != nil {
		log.Fatal(err)
	}

	caKey, err := ioutil.ReadFile(KEY)
	if err != nil {
		log.Fatal(err)
	}

	err = ca.SetCA(caCert, caKey)
	if err != nil {
		log.Fatalf("CA Set Error")
	}
	server.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	server.Verbose = verbose
	server.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		req.Header.Set("SM-Host", getUrl(req))
		req.URL, err = url.Parse(proxyAddr)
		req.Host = req.URL.Host
		return req, nil
	})
	http.ListenAndServe(listenAddr, server)
}

func handler(w http.ResponseWriter, req *http.Request, unzip bool) error {
	resp, err := doHttp(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "application/json")
		data := map[string]string{
			"errorMessage": err.Error(),
			"errorType":    "httpError",
		}
		encoder := json.NewEncoder(w)
		err = encoder.Encode(data)
		if err != nil {
			return err
		}
		return nil
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		w.WriteHeader(resp.StatusCode)
		if unzip && resp.Header.Get("Content-Encoding") == "gzip" {
			w.Write(unGzip(body))
		} else {
			w.Write(body)
		}
		for key, value := range resp.Header {
			for i := 0; i < len(value); i++ {
				w.Header().Add(key, value[i])
			}
		}
	}
	return nil
}

func doHttp(req *http.Request) (*http.Response, error) {
	// TODO timeout setting
	client := &http.Client{Timeout: 10 * time.Second}
	url := getUrl(req)
	proxyReq, err := http.NewRequest(req.Method, url, req.Body)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("http redirect request")
	}
	proxyReq.Header = req.Header
	proxyReq.Proto = req.Proto
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(proxyReq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func getUrl(req *http.Request) string {
	url := ""
	if strings.HasPrefix(req.RequestURI, "http://") || strings.HasPrefix(req.RequestURI, "https://") {
		url = req.RequestURI
	} else {
		scheme := "http://"
		if req.TLS != nil {
			scheme = "https://"
		}
		if req.URL.Scheme != "" {
			scheme = req.URL.Scheme + "://"
		}
		url = strings.Join([]string{scheme, req.Host, req.URL.Path + "?" + req.URL.RawQuery}, "")
	}
	return url
}

func unGzip(body []byte) []byte {
	gzipReader, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		return []byte("unGzip error")
	}
	defer gzipReader.Close()
	var buf bytes.Buffer
	_, err = io.Copy(&buf, gzipReader)
	if err != nil {
		return []byte("IO Copy error")
	}
	return buf.Bytes()
}
