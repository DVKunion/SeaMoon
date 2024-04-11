package network

import (
	"io"
	"net"
	"sync"

	"github.com/gorilla/websocket"
)

const bufferSize = 64 * 1024

// Transport rw1 and rw2
func Transport(src, dest net.Conn) (int64, int64, error) {
	// 这是第一版本的老代码，是可以满足需求但是很丑陋
	errIn := make(chan error, 1)
	errOut := make(chan error, 1)
	done := make(chan error, 1<<8)

	var wg sync.WaitGroup
	var mu sync.Mutex // 用于保护下面的共享变量

	var inbound int64 = 0
	var outbound int64 = 0

	wg.Add(2)
	go func() {
		defer wg.Done()
		w, err := CopyBuffer(src, dest, bufferSize)
		mu.Lock()
		inbound += w
		mu.Unlock()
		errIn <- err
	}()

	go func() {
		defer wg.Done()
		w, err := CopyBuffer(dest, src, bufferSize)
		mu.Lock()
		outbound += w
		mu.Unlock()
		errOut <- err
	}()

	for {
		select {
		case e := <-errIn:
			src.Close()
			dest.Close()
			done <- e
		case e := <-errOut:
			src.Close()
			dest.Close()
			done <- e
		case e := <-done:
			wg.Wait()
			// 忽略 websocket 正常断开
			if opErr, ok := e.(*net.OpError); ok {
				if closeErr, ok := opErr.Err.(*websocket.CloseError); ok {
					if closeErr.Code == websocket.CloseNormalClosure || closeErr.Code == websocket.CloseAbnormalClosure {
						e = nil
					}
				}
			}
			return inbound, outbound, e
		}
	}
}

func CopyBuffer(dst io.Writer, src io.Reader, bufSize int) (int64, error) {
	buf := GetBuffer(bufSize)
	defer PutBuffer(buf)

	return io.CopyBuffer(dst, src, buf)
}
