package listener

import (
	"bytes"
	"context"
	"net"
	"sync"
	"time"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	db_service "github.com/DVKunion/SeaMoon/pkg/api/service"
	"github.com/DVKunion/SeaMoon/pkg/network/basic"
	"github.com/DVKunion/SeaMoon/pkg/network/tunnel/service"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

const (
	udpSessionTimeout = 60 * time.Second
	udpBufferSize     = 65535
)

type udpSession struct {
	conn       net.Conn
	lastActive time.Time
}

func UDPListen(ctx context.Context, py *models.Proxy) (net.PacketConn, error) {
	if *py.Type != enum.ProxyTypeAUTO && *py.Type != enum.ProxyTypeSOCKS5 {
		xlog.Warn("udp listener only support socks5 and auto proxy type")
		return nil, nil
	}
	server, err := net.ListenPacket("udp", py.Addr())
	if err != nil {
		return nil, err
	}

	tun, err := db_service.SVC.GetTunnelById(ctx, py.TunnelID)
	if err != nil {
		return nil, err
	}

	go listenUDP(ctx, server, py.ID, py.Type, tun)

	return server, nil
}

func listenUDP(ctx context.Context, server net.PacketConn, id uint, t *enum.ProxyType, tun *models.Tunnel) {
	var (
		sessions = make(map[string]*udpSession)
		mu       sync.Mutex
		buf      = make([]byte, udpBufferSize)
	)

	// Cleanup routine
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				mu.Lock()
				now := time.Now()
				for k, s := range sessions {
					if now.Sub(s.lastActive) > udpSessionTimeout {
						s.conn.Close()
						delete(sessions, k)
						db_service.SVC.UpdateProxyConn(ctx, id, -1)
						xlog.Debug("udp session expired", "client", k)
					}
				}
				mu.Unlock()
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			server.Close()
			return
		default:
			n, addr, err := server.ReadFrom(buf)
			if err != nil {
				// avoid log spam on close
				return
			}

			clientAddr := addr.String()
			mu.Lock()
			sess, ok := sessions[clientAddr]
			if ok {
				sess.lastActive = time.Now()
			}
			mu.Unlock()

			if !ok {
				// New Session
				srv, ok := service.Factory.Load(*tun.Type)
				if !ok {
					continue
				}

				// Dial Tunnel
				destConn, err := srv.(service.Service).Conn(ctx, *t,
					service.WithAddr(tun.GetAddr()),
					service.WithTorFlag(tun.Config.Tor),
					service.WithPath("/socks5-udp"),
				)
				if err != nil {
					xlog.Error("failed to dial tunnel udp", "err", err)
					continue
				}

				sess = &udpSession{
					conn:       destConn,
					lastActive: time.Now(),
				}

				mu.Lock()
				sessions[clientAddr] = sess
				db_service.SVC.UpdateProxyConn(ctx, id, 1)
				mu.Unlock()

				// Start Tunnel -> Client handler
				go func(c net.Conn, target net.Addr, key string) {
					defer func() {
						c.Close()
						mu.Lock()
						if _, ok := sessions[key]; ok {
							delete(sessions, key)
							db_service.SVC.UpdateProxyConn(ctx, id, -1)
						}
						mu.Unlock()
					}()

					for {
						d, err := basic.ReadUDPDatagram(c)
						if err != nil {
							return
						}

						// Reset RSV to 0 for standard SOCKS5
						d.Header.Rsv = 0

						var b bytes.Buffer
						if err := d.Write(&b); err != nil {
							continue
						}

						if _, err := server.WriteTo(b.Bytes(), target); err != nil {
							// xlog.Error("udp write to client error", "err", err)
						}
					}
				}(destConn, addr, clientAddr)
			}

			// Write to Tunnel
			if _, err := sess.conn.Write(buf[:n]); err != nil {
				xlog.Error("udp write to tunnel error", "err", err)
			}
		}
	}
}
