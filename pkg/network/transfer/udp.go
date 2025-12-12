package transfer

import (
	"bytes"
	"encoding/binary"
	"net"

	"github.com/DVKunion/SeaMoon/pkg/network/basic"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

// Socks5UDPTransport handles UDP forwarding over a stream connection (e.g. WebSocket)
// It listens on a local ephemeral UDP port to send/receive packets to/from the Internet.
// It reads SOCKS5 UDP Datagrams from conn, sends them to Internet.
// It reads UDP packets from Internet, wraps them in SOCKS5 UDP Datagrams, sends to conn.
func Socks5UDPTransport(conn net.Conn) error {
	// UDP socket for outgoing traffic to Targets
	udpConn, err := net.ListenPacket("udp", ":0")
	if err != nil {
		return err
	}
	defer udpConn.Close()

	return Socks5UDPProxy(conn, udpConn)
}

// Socks5UDPProxy handles the actual forwarding logic between a Stream and a PacketConn (Server Side)
func Socks5UDPProxy(conn net.Conn, udpConn net.PacketConn) error {
	errCh := make(chan error, 2)

	// 1. Client (Stream) -> Target (UDP)
	go func() {
		for {
			// Read framed UDP packet from Stream
			d, err := basic.ReadUDPDatagram(conn)
			if err != nil {
				errCh <- err
				return
			}

			xlog.Debug("udp forward", "target", d.Header.Addr.String(), "len", len(d.Data))

			targetAddr, err := net.ResolveUDPAddr("udp", d.Header.Addr.String())
			if err != nil {
				xlog.Error("resolve udp addr error", "err", err)
				continue
			}

			if _, err := udpConn.WriteTo(d.Data, targetAddr); err != nil {
				xlog.Error("udp write error", "err", err)
			}
		}
	}()

	// 2. Target (UDP) -> Client (Stream)
	go func() {
		buf := make([]byte, 65535)
		for {
			n, addr, err := udpConn.ReadFrom(buf)
			if err != nil {
				errCh <- err
				return
			}

			// Create SOCKS5 UDP Datagram
			// We need to construct the header.
			// And we should set RSV to len(data).

			udpAddr, ok := addr.(*net.UDPAddr)
			if !ok {
				continue
			}

			header := basic.NewUDPHeader(uint16(n), 0, &basic.Addr{
				Type: basic.AddrIPv4, // Default to IPv4, need to check
				Host: udpAddr.IP.String(),
				Port: uint16(udpAddr.Port),
			})
			// Fix Addr Type
			if udpAddr.IP.To4() == nil {
				header.Addr.Type = basic.AddrIPv6
			}

			d := basic.NewUDPDatagram(header, buf[:n])

			if err := d.Write(conn); err != nil {
				xlog.Error("udp write to stream error", "err", err)
				errCh <- err
				return
			}
		}
	}()

	return <-errCh
}

// ClientUDPProxy handles the forwarding for the Client side (Local Listener <-> Tunnel)
func ClientUDPProxy(udpConn net.PacketConn, tunnelConn net.Conn) error {
	var clientAddr net.Addr
	errCh := make(chan error, 2)

	// 1. UDP (Client) -> Tunnel
	go func() {
		buf := make([]byte, 65535)
		for {
			n, addr, err := udpConn.ReadFrom(buf)
			if err != nil {
				errCh <- err
				return
			}
			// Update client address
			clientAddr = addr

			if n < 4 {
				continue // Invalid
			}

			atype := buf[3]
			hlen := 0
			switch atype {
			case basic.AddrIPv4:
				hlen = 10
			case basic.AddrIPv6:
				hlen = 22
			case basic.AddrDomain:
				if n < 5 {
					continue
				}
				hlen = 7 + int(buf[4])
			default:
				continue
			}

			if n < hlen {
				continue
			}

			dataLen := n - hlen

			// Set RSV to dataLen
			binary.BigEndian.PutUint16(buf[:2], uint16(dataLen))

			// Write to Tunnel
			if _, err := tunnelConn.Write(buf[:n]); err != nil {
				errCh <- err
				return
			}
		}
	}()

	// 2. Tunnel -> UDP (Client)
	go func() {
		for {
			d, err := basic.ReadUDPDatagram(tunnelConn)
			if err != nil {
				errCh <- err
				return
			}

			d.Header.Rsv = 0

			var buf bytes.Buffer
			if err := d.Write(&buf); err != nil {
				continue
			}

			if clientAddr != nil {
				if _, err := udpConn.WriteTo(buf.Bytes(), clientAddr); err != nil {
					xlog.Error("udp write to client error", "err", err)
				}
			}
		}
	}()

	return <-errCh
}
