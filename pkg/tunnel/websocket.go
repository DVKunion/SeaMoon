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
