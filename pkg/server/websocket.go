package server

import (
	"io"
	"net"
	"sync"

	"github.com/gorilla/websocket"
)

type WebsocketServer struct {
	net.Conn
	wConn     *websocket.Conn
	writeLock *sync.Mutex

	messageType int
}

func NewWebsocketServer(wConn *websocket.Conn, lock *sync.Mutex) net.Conn {
	return &WebsocketServer{
		wConn:       wConn,
		messageType: websocket.BinaryMessage,
		writeLock:   lock,
	}
}

func (ws *WebsocketServer) RemoteAddr() net.Addr {
	return ws.wConn.RemoteAddr()
}

func (ws *WebsocketServer) Close() error {
	// ws need send close message first to avoid err : close 1006 (abnormal closure): unexpected EOF
	err := ws.write(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "close"))
	if err != nil {
		return err
	}
	return ws.wConn.Close()
}

func (ws *WebsocketServer) Write(b []byte) (n int, err error) {
	err = ws.write(ws.messageType, b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (ws *WebsocketServer) write(messageType int, data []byte) error {
	ws.writeLock.TryLock()
	defer ws.writeLock.Unlock()
	return ws.wConn.WriteMessage(messageType, data)
}

func (ws *WebsocketServer) Read(b []byte) (n int, err error) {
	mt, message, err := ws.wConn.ReadMessage()
	if err != nil {
		if wsErr, ok := err.(*websocket.CloseError); ok && wsErr.Code == websocket.CloseNormalClosure {
			return 0, io.EOF
		}
		return 0, err
	}
	ws.messageType = mt
	copy(b, message)
	return len(message), nil
}
