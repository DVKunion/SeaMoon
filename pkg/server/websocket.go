package server

import (
	"github.com/gorilla/websocket"
	"net"
)

type WebsocketServer struct {
	net.Conn
	wConn       *websocket.Conn
	messageType int
}

func NewWebsocketServer(wConn *websocket.Conn) net.Conn {
	return &WebsocketServer{
		wConn:       wConn,
		messageType: websocket.BinaryMessage,
	}
}

func (ws *WebsocketServer) RemoteAddr() net.Addr {
	return ws.wConn.RemoteAddr()
}

func (ws *WebsocketServer) Close() error {
	return ws.wConn.Close()
}

func (ws *WebsocketServer) Write(b []byte) (n int, err error) {
	err = ws.wConn.WriteMessage(ws.messageType, b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (ws *WebsocketServer) Read(b []byte) (n int, err error) {
	mt, message, err := ws.wConn.ReadMessage()
	if err != nil {
		return 0, err
	}
	ws.messageType = mt
	copy(b, message)
	return len(message), nil
}
