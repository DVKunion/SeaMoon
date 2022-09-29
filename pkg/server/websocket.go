package server

import (
	"github.com/gorilla/websocket"
	"io"
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
	// ws need send close message first to avoid err : close 1006 (abnormal closure): unexpected EOF
	// todo: panic - concurrent write to websocket connection
	err := ws.wConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "close"))
	if err != nil {
		return err
	}
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
		if wsErr, ok := err.(*websocket.CloseError); ok && wsErr.Code == websocket.CloseNormalClosure {
			return 0, io.EOF
		}
		return 0, err
	}
	ws.messageType = mt
	copy(b, message)
	return len(message), nil
}
