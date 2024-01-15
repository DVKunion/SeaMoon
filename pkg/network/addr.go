package network

import (
	"encoding/binary"
	"net"
	"strconv"
)

/*
Addr has following struct

	+------+----------+----------+
	| Type |   ADDR   |   PORT   |
	+------+----------+----------+
	|  1   | Variable |    2     |
	+------+----------+----------+
*/
type Addr struct {
	Type uint8
	Host string
	Port uint16
}

func IsIPv4(address string) bool {
	return address != "" && address[0] != ':' && address[0] != '['
}

// NewAddr creates an address object
func NewAddr(sa string) (addr *Addr, err error) {
	host, sport, err := net.SplitHostPort(sa)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(sport)
	if err != nil {
		return nil, err
	}

	addr = NewAddrFromPair(host, port)
	return
}

// NewAddrFromPair creates an address object from host and port pair
func NewAddrFromPair(host string, port int) (addr *Addr) {
	addr = &Addr{
		Type: AddrDomain,
		Host: host,
		Port: uint16(port),
	}

	if ip := net.ParseIP(host); ip != nil {
		if ip.To4() != nil {
			addr.Type = AddrIPv4
		} else {
			addr.Type = AddrIPv6
		}
	}

	return
}

// NewAddrFromAddr creates an address object
func NewAddrFromAddr(ln, conn net.Addr) (addr *Addr, err error) {
	_, sport, err := net.SplitHostPort(ln.String())
	if err != nil {
		return nil, err
	}
	host, _, err := net.SplitHostPort(conn.String())
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(sport)
	if err != nil {
		return nil, err
	}

	addr = NewAddrFromPair(host, port)
	return
}

// Decode an address from the stream
func (addr *Addr) Decode(b []byte) error {
	addr.Type = b[0]
	pos := 1
	switch addr.Type {
	case AddrIPv4:
		addr.Host = net.IP(b[pos : pos+net.IPv4len]).String()
		pos += net.IPv4len
	case AddrIPv6:
		addr.Host = net.IP(b[pos : pos+net.IPv6len]).String()
		pos += net.IPv6len
	case AddrDomain:
		addrlen := int(b[pos])
		pos++
		addr.Host = string(b[pos : pos+addrlen])
		pos += addrlen
	default:
		return ErrBadAddrType
	}

	addr.Port = binary.BigEndian.Uint16(b[pos:])

	return nil
}

// Encode an address to the stream
func (addr *Addr) Encode(b []byte) (int, error) {
	b[0] = addr.Type
	pos := 1
	switch addr.Type {
	case AddrIPv4:
		ip4 := net.ParseIP(addr.Host).To4()
		if ip4 == nil {
			ip4 = net.IPv4zero.To4()
		}
		pos += copy(b[pos:], ip4)
	case AddrDomain:
		b[pos] = byte(len(addr.Host))
		pos++
		pos += copy(b[pos:], []byte(addr.Host))
	case AddrIPv6:
		ip16 := net.ParseIP(addr.Host).To16()
		if ip16 == nil {
			ip16 = net.IPv6zero.To16()
		}
		pos += copy(b[pos:], ip16)
	default:
		b[0] = AddrIPv4
		copy(b[pos:pos+4], net.IPv4zero.To4())
		pos += 4
	}
	binary.BigEndian.PutUint16(b[pos:], addr.Port)
	pos += 2

	return pos, nil
}

// Length of the address
func (addr *Addr) Length() (n int) {
	switch addr.Type {
	case AddrIPv4:
		n = 10
	case AddrIPv6:
		n = 22
	case AddrDomain:
		n = 7 + len(addr.Host)
	default:
		n = 10
	}
	return
}

func (addr *Addr) String() string {
	return net.JoinHostPort(addr.Host, strconv.Itoa(int(addr.Port)))
}
