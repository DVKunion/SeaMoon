package network

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"time"
)

// DoHttp the basic of http request ability
func DoHttp(req *http.Request) (*http.Response, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetUrl 从 req 类中获取完整 URL信息
func GetUrl(req *http.Request) string {
	// TODO check more
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

// UnGzip do UnGzip
func UnGzip(body []byte) []byte {
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

// HealthResponse wrapper a response with plain message
func HealthResponse(message string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(message))
	if err != nil {
		// TODO log
		return
	}
}

func ErrorResponse(message string, statusCode int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(message))
	if err != nil {
		// TODO log
		return
	}
}
