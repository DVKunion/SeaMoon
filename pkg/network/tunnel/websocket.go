package tunnel

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type websocketConn struct {
	*websocket.Conn
	rb  []byte
	mux sync.Mutex
}

func (c *websocketConn) Delay() int64 {

	if c.Conn == nil {
		return 0
	}

	// 设置一个超时时间
	const timeout = 10 * time.Second

	// 创建一个通道用于接收pong消息的回复时间
	pongReceived := make(chan time.Time, 1) // 缓冲为1，避免潜在的阻塞

	// 设置Pong处理程序
	c.Conn.SetPongHandler(func(appData string) error {
		pongReceived <- time.Now()
		return nil
	})

	pingSentTime := time.Now()
	if err := c.WriteMessage(websocket.PingMessage, []byte("")); err != nil {
		return 0
	}

	select {
	case pongTime := <-pongReceived:
		res := pongTime.Sub(pingSentTime).Milliseconds()
		return res
	case <-time.After(timeout):
		return 99999
	}
}

func WsWrapConn(conn *websocket.Conn) Tunnel {
	return &websocketConn{
		Conn: conn,
	}
}

func (c *websocketConn) Read(b []byte) (n int, err error) {
	if len(c.rb) == 0 {
		_, c.rb, err = c.Conn.ReadMessage()
	}
	n = copy(b, c.rb)
	c.rb = c.rb[n:]
	return
}

func (c *websocketConn) Write(b []byte) (n int, err error) {
	err = c.WriteMessage(websocket.BinaryMessage, b)
	n = len(b)
	return
}

func (c *websocketConn) WriteMessage(messageType int, data []byte) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.Conn.WriteMessage(messageType, data)
}

func (c *websocketConn) SetDeadline(t time.Time) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	if err := c.SetReadDeadline(t); err != nil {
		return err
	}
	return c.SetWriteDeadline(t)
}

func (c *websocketConn) SetWriteDeadline(t time.Time) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.Conn.SetWriteDeadline(t)
}
