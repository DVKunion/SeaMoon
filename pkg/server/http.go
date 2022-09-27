package server

import (
	"errors"
	"fmt"
	"github.com/DVKunion/SeaMoon/pkg/consts"
	"github.com/DVKunion/SeaMoon/pkg/utils"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var HttpForwardHeader = map[string]string{"SM-HOST": "", "Proxy-Authenticate": "", "Proxy-Authorization": "",
	"Connection": "", "Keep-Alive": "",
	"Proxy-Connection": "",
	"Te":               "", "Trailer": "", "Transfer-Encoding": "", "Content-Encoding": ""}

type HttpServer struct {
	CloudServer
	remoteUrl  *url.URL
	remoteHost string
}

func (s *Server) HttpServerTransfer(r *http.Request) CloudServer {
	return newHttpServer(r)
}

func (s *HttpServer) Verification(w http.ResponseWriter) (bool, error) {
	// forward host 不为空
	// 环路判断问题: 因为代理发出的请求是不带 SM-HOST 标识的，所以这个return会直接阻止了环路。
	// 也就是说，最多只会重复请求两次，即代理请求一次自己后，返回了此处的异常。
	if s.remoteHost == "" {
		var errMsg = "no sm-ost in request"
		utils.HealthResponse(fmt.Sprintf(consts.DEFAULT_HTTP, errMsg), w)
		return false, errors.New(errMsg)
	}
	var err error
	s.remoteUrl, err = url.Parse(s.remoteHost)
	if err != nil {
		utils.HealthResponse(consts.DEFAULT_HTTP, w)
		return false, err
	}
	return true, nil
}

func (s *HttpServer) Serve(w http.ResponseWriter, req *http.Request) {
	conn := req.Header.Get("Proxy-Connection")
	for key, _ := range HttpForwardHeader {
		req.Header.Del(key)
	}
	req.Header.Set("Host", s.remoteHost)
	if conn != "" {
		req.Header.Set("Connection", conn)
	}
	req.URL = s.remoteUrl
	req.Host = s.remoteUrl.Host
	err := s.httpHandler(w, req)
	if err != nil {
		log.Error(err)
	}
}

func newHttpServer(r *http.Request) *HttpServer {
	host := r.Header.Get("SM-HOST")
	return &HttpServer{
		remoteHost: host,
	}
}

func (s *HttpServer) httpHandler(w http.ResponseWriter, req *http.Request) error {
	log.Infof(consts.HTTP_ACCEPT_START, req.RemoteAddr)
	proxyReq, err := http.NewRequest(req.Method, s.remoteHost, req.Body)
	if err != nil {
		return err
	}
	proxyReq.Header = req.Header
	proxyReq.Proto = req.Proto
	response, err := utils.DoHttp(proxyReq)
	if err != nil {
		utils.ErrorResponse("SeaMoon HTTP Error: "+err.Error(), http.StatusBadRequest, w)
		log.Errorf(consts.HTTP_ACCEPT_ERROR, err.Error())
		return err
	} else {
		// head
		for key, value := range response.Header {
			for i := 0; i < len(value); i++ {
				if _, ok := HttpForwardHeader[key]; ok {
					continue
				}
				if key == "Content-Length" {
					continue
				}
				w.Header().Add(key, value[i])
			}
		}
		// proto
		w.WriteHeader(response.StatusCode)

		// body
		body, errRead := ioutil.ReadAll(response.Body)
		if errRead != nil {
			return errRead
		}
		defer response.Body.Close()

		if strings.Contains(response.Header.Get("Content-Encoding"), "gzip") {
			body = utils.UnGzip(body)
		}
		// write
		if _, err := w.Write(body); err != nil {
			return err
		}
	}
	log.Infof(consts.HTTP_BODY_DIS, response.Request.Host)
	return nil
}
