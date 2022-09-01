package server

import (
	"github.com/DVKunion/SeaMoon/pkg/consts"
	"github.com/DVKunion/SeaMoon/pkg/utils"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var HttpForwardHeader = "SM-HOST"

type HttpServer struct {
	CloudServer
	remoteUrl  *url.URL
	remoteHost string
}

func (s *Server) HttpServerTransfer(r *http.Request) CloudServer {
	return newHttpServer(r)
}

func (s *HttpServer) Verification(w http.ResponseWriter) bool {
	// forward host 不为空
	// 环路判断问题: 因为代理发出的请求是不带 SM-HOST 标识的，所以这个return会直接阻止了环路。
	// 也就是说，最多只会重复请求两次，即代理请求一次自己后，返回了此处的异常。
	if s.remoteHost == "" {
		utils.HealthResponse(consts.DEFAULT_HTTP, w)
		return false
	}
	var err error
	s.remoteUrl, err = url.Parse(s.remoteHost)
	if err != nil {
		utils.HealthResponse(consts.DEFAULT_HTTP, w)
		return false
	}
	return true
}

func (s *HttpServer) Serve(w http.ResponseWriter, req *http.Request) {
	req.Header.Del(HttpForwardHeader)
	req.Header.Set("Host", s.remoteHost)
	req.URL = s.remoteUrl
	req.Host = s.remoteUrl.Host
	//s.req.RequestURI = host
	err := s.httpHandler(w, req)
	if err != nil {
		log.Error("TODO")
	}
}

func newHttpServer(r *http.Request) *HttpServer {
	host := r.Header.Get(HttpForwardHeader)
	return &HttpServer{
		remoteHost: host,
	}
}

func (s *HttpServer) httpHandler(w http.ResponseWriter, req *http.Request) error {
	proxyReq, err := http.NewRequest(req.Method, s.remoteHost, req.Body)
	if err != nil {
		return err
	}
	proxyReq.Header = req.Header
	proxyReq.Proto = req.Proto
	response, err := utils.DoHttp(proxyReq)
	if err != nil {
		utils.ErrorResponse("SeaMoon HTTP Error: "+err.Error(), http.StatusBadRequest, w)
		return err
	} else {
		body, errRead := ioutil.ReadAll(response.Body)
		if errRead != nil {
			return errRead
		}
		defer response.Body.Close()
		w.WriteHeader(response.StatusCode)

		if strings.Contains(response.Header.Get("Content-Encoding"), "gzip") {
			unzip := utils.UnGzip(body)
			w.Write(unzip)
		} else {
			w.Write(body)
		}

		for key, value := range response.Header {
			for i := 0; i < len(value); i++ {
				w.Header().Add(key, value[i])
			}
		}
	}
	return nil
}
