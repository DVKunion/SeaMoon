package transfer

import (
	"bufio"
	"net"
	"net/http"

	"github.com/DVKunion/SeaMoon/pkg/network"
)

func HttpTransport(conn net.Conn) error {
	// 接收客户端的连接，并从第一条消息中获取目标地址
	request, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		return err
	}

	defer request.Body.Close()

	targetAddr := request.Host

	if targetAddr == "" || err != nil {
		return err
	}

	if _, port, _ := net.SplitHostPort(targetAddr); port == "" {
		targetAddr = net.JoinHostPort(targetAddr, "80")
	}

	// 连接到目标服务器
	destConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		return err
	}

	defer destConn.Close()

	resp := &http.Response{
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	if resp.Header == nil {
		resp.Header = http.Header{}
	}

	if request.Method == http.MethodConnect {
		resp.StatusCode = http.StatusOK
		resp.Status = "200 Connection established"

		if err = resp.Write(conn); err != nil {
			return err
		}
	} else {
		if err = request.Write(destConn); err != nil {
			return err
		}
	}

	// 同时处理客户端到服务器和服务器到客户端的数据传输
	if _, _, err := network.Transport(destConn, conn); err != nil {
		return err
	}

	return nil
}
