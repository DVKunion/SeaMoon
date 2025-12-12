package basic

// This file is modified version from https://github.com/ginuerzh/gosocks5/blob/master/socks5.go

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/DVKunion/SeaMoon/pkg/system/errors"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

// Address types
const (
	AddrIPv4   uint8 = 1
	AddrDomain       = 3
	AddrIPv6         = 4
)

// SOCKS5 Version = 5
const (
	SOCKS5Version     = 5
	SOCKS5UserPassVer = 1
)

// Methods
const (
	MethodNoAuth uint8 = iota
	MethodGSSAPI
	MethodUserPass
	MethodNoAcceptable uint8 = 0xFF
)

// SOCKS5 Commands
const (
	SOCKS5CmdConnect uint8 = iota + 1
	SOCKS5CmdBind
	SOCKS5CmdUDPOverTCP
)

// SOCKS5 Response codes
const (
	SOCKS5RespSucceeded uint8 = iota
	SOCKS5RespFailure
	SOCKS5RespAllowed
	SOCKS5RespNetUnreachable
	SOCKS5RespHostUnreachable
	SOCKS5RespConnRefused
	SOCKS5RespTTLExpired
	SOCKS5RespCmdUnsupported
	SOCKS5RespAddrUnsupported
)

const (
	smallSize = 576
	largeSize = 32 * 1024
)

/*
ReadMethods returns methods
Method selection

	+----+----------+----------+
	|VER | NMETHODS | METHODS  |
	+----+----------+----------+
	| 1  |    1     | 1 to 255 |
	+----+----------+----------+
*/
func ReadMethods(r io.Reader) ([]uint8, error) {
	b := GetBuffer(smallSize)
	defer PutBuffer(b)

	n, err := io.ReadAtLeast(r, b, 2)
	if err != nil {
		return nil, err
	}

	if b[0] != SOCKS5Version {
		return nil, errors.New(xlog.NetworkVersionError)
	}

	if b[1] == 0 {
		return nil, errors.New(xlog.NetworkMethodError)
	}

	length := 2 + int(b[1])
	if n < length {
		if _, err := io.ReadFull(r, b[n:length]); err != nil {
			return nil, err
		}
	}

	methods := make([]byte, int(b[1]))
	copy(methods, b[2:length])

	return methods, nil
}

// WriteMethod send the selected method to the client
func WriteMethod(method uint8, w io.Writer) error {
	_, err := w.Write([]byte{SOCKS5Version, method})
	return err
}

// WriteMethods send method select request to the server
func WriteMethods(methods []uint8, w io.Writer) error {
	b := make([]byte, 2+len(methods))
	b[0] = SOCKS5Version
	b[1] = uint8(len(methods))
	copy(b[2:], methods)

	_, err := w.Write(b)
	return err
}

/*
UserPassRequest Username/Password authentication request

	+----+------+----------+------+----------+
	|VER | ULEN |  UNAME   | PLEN |  PASSWD  |
	+----+------+----------+------+----------+
	| 1  |  1   | 1 to 255 |  1   | 1 to 255 |
	+----+------+----------+------+----------+
*/
type UserPassRequest struct {
	Version  byte
	Username string
	Password string
}

func NewUserPassRequest(ver byte, u, p string) *UserPassRequest {
	return &UserPassRequest{
		Version:  ver,
		Username: u,
		Password: p,
	}
}

func ReadUserPassRequest(r io.Reader) (*UserPassRequest, error) {
	b := GetBuffer(smallSize)
	defer PutBuffer(b)

	n, err := io.ReadAtLeast(r, b, 2)
	if err != nil {
		return nil, err
	}

	if b[0] != SOCKS5UserPassVer {
		return nil, errors.New(xlog.NetworkVersionError)
	}

	req := &UserPassRequest{
		Version: b[0],
	}

	ulen := int(b[1])
	length := ulen + 3

	if n < length {
		if _, err := io.ReadFull(r, b[n:length]); err != nil {
			return nil, err
		}
		n = length
	}
	req.Username = string(b[2 : 2+ulen])

	plen := int(b[length-1])
	length += plen
	if n < length {
		if _, err := io.ReadFull(r, b[n:length]); err != nil {
			return nil, err
		}
	}
	req.Password = string(b[3+ulen : length])
	return req, nil
}

func (req *UserPassRequest) Write(w io.Writer) error {
	// b := make([]byte, 513)
	b := GetBuffer(smallSize)
	defer PutBuffer(b)

	b[0] = req.Version
	ulen := len(req.Username)
	b[1] = byte(ulen)
	length := 2 + ulen
	copy(b[2:length], req.Username)

	plen := len(req.Password)
	b[length] = byte(plen)
	length++
	copy(b[length:length+plen], req.Password)
	length += plen

	_, err := w.Write(b[:length])
	return err
}

func (req *UserPassRequest) String() string {
	return fmt.Sprintf("%d %s:%s",
		req.Version, req.Username, req.Password)
}

/*
UserPassResponse Username/Password authentication response

	+----+--------+
	|VER | STATUS |
	+----+--------+
	| 1  |   1    |
	+----+--------+
*/
type UserPassResponse struct {
	Version byte
	Status  byte
}

func NewUserPassResponse(ver, status byte) *UserPassResponse {
	return &UserPassResponse{
		Version: ver,
		Status:  status,
	}
}

func ReadUserPassResponse(r io.Reader) (*UserPassResponse, error) {
	// b := make([]byte, 2)
	b := GetBuffer(smallSize)
	defer PutBuffer(b)

	if _, err := io.ReadFull(r, b[:2]); err != nil {
		return nil, err
	}

	if b[0] != SOCKS5UserPassVer {
		return nil, errors.New(xlog.NetworkVersionError)
	}

	res := &UserPassResponse{
		Version: b[0],
		Status:  b[1],
	}

	return res, nil
}

func (res *UserPassResponse) Write(w io.Writer) error {
	_, err := w.Write([]byte{res.Version, res.Status})
	return err
}

func (res *UserPassResponse) String() string {
	return fmt.Sprintf("%d %d",
		res.Version, res.Status)
}

/*
SOCKS5Request represent a socks5 request
The SOCKSv5 request

	+----+-----+-------+------+----------+----------+
	|VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	+----+-----+-------+------+----------+----------+
	| 1  |  1  | X'00' |  1   | Variable |    2     |
	+----+-----+-------+------+----------+----------+
*/
type SOCKS5Request struct {
	Cmd  uint8
	Addr *Addr
}

// NewSOCKS5Request creates an request object
func NewSOCKS5Request(cmd uint8, addr *Addr) *SOCKS5Request {
	return &SOCKS5Request{
		Cmd:  cmd,
		Addr: addr,
	}
}

// ReadSOCKS5Request reads request from the stream
func ReadSOCKS5Request(r io.Reader) (*SOCKS5Request, error) {
	// b := make([]byte, 262)
	b := GetBuffer(smallSize)
	defer PutBuffer(b)

	n, err := io.ReadAtLeast(r, b, 5)
	if err != nil {
		return nil, err
	}

	if b[0] != SOCKS5Version {
		return nil, errors.New(xlog.NetworkVersionError)
	}

	request := &SOCKS5Request{
		Cmd: b[1],
	}

	atype := b[3]
	length := 0
	switch atype {
	case AddrIPv4:
		length = 10
	case AddrIPv6:
		length = 22
	case AddrDomain:
		length = 7 + int(b[4])
	default:
		return nil, errors.New(xlog.NetworkAddrTypeError)
	}

	if n < length {
		if _, err := io.ReadFull(r, b[n:length]); err != nil {
			return nil, err
		}
	}
	addr := new(Addr)
	if err := addr.Decode(b[3:length]); err != nil {
		return nil, err
	}
	request.Addr = addr

	return request, nil
}

func (r *SOCKS5Request) Write(w io.Writer) (err error) {
	//b := make([]byte, 262)
	b := GetBuffer(smallSize)
	defer PutBuffer(b)

	b[0] = SOCKS5Version
	b[1] = r.Cmd
	b[2] = 0        //rsv
	b[3] = AddrIPv4 // default

	addr := r.Addr
	if addr == nil {
		addr = &Addr{}
	}
	n, _ := addr.Encode(b[3:])
	length := 3 + n

	_, err = w.Write(b[:length])
	return
}

func (r *SOCKS5Request) String() string {
	addr := r.Addr
	if addr == nil {
		addr = &Addr{}
	}
	return fmt.Sprintf("5 %d 0 %d %s",
		r.Cmd, addr.Type, addr.String())
}

/*
Reply is a SOCKSv5 reply

	+----+-----+-------+------+----------+----------+
	|VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	+----+-----+-------+------+----------+----------+
	| 1  |  1  | X'00' |  1   | Variable |    2     |
	+----+-----+-------+------+----------+----------+
*/
type Reply struct {
	Rep  uint8
	Addr *Addr
}

// NewReply creates a socks5 reply
func NewReply(rep uint8, addr *Addr) *Reply {
	return &Reply{
		Rep:  rep,
		Addr: addr,
	}
}

// ReadReply reads a reply from the stream
func ReadReply(r io.Reader) (*Reply, error) {
	// b := make([]byte, 262)
	b := GetBuffer(smallSize)
	defer PutBuffer(b)

	n, err := io.ReadAtLeast(r, b, 5)
	if err != nil {
		return nil, err
	}

	if b[0] != SOCKS5Version {
		return nil, errors.New(xlog.NetworkVersionError)
	}

	reply := &Reply{
		Rep: b[1],
	}

	atype := b[3]
	length := 0
	switch atype {
	case AddrIPv4:
		length = 10
	case AddrIPv6:
		length = 22
	case AddrDomain:
		length = 7 + int(b[4])
	default:
		return nil, errors.New(xlog.NetworkAddrTypeError)
	}

	if n < length {
		if _, err := io.ReadFull(r, b[n:length]); err != nil {
			return nil, err
		}
	}

	addr := new(Addr)
	if err := addr.Decode(b[3:length]); err != nil {
		return nil, err
	}
	reply.Addr = addr

	return reply, nil
}

func (r *Reply) Write(w io.Writer) (err error) {
	// b := make([]byte, 262)
	b := GetBuffer(smallSize)
	defer PutBuffer(b)

	b[0] = SOCKS5Version
	b[1] = r.Rep
	b[2] = 0        //rsv
	b[3] = AddrIPv4 // default
	length := 10
	b[4], b[5], b[6], b[7], b[8], b[9] = 0, 0, 0, 0, 0, 0 // reset address field

	if r.Addr != nil {
		n, _ := r.Addr.Encode(b[3:])
		length = 3 + n
	}
	_, err = w.Write(b[:length])

	return
}

func (r *Reply) String() string {
	addr := r.Addr
	if addr == nil {
		addr = &Addr{}
	}
	return fmt.Sprintf("5 %d 0 %d %s",
		r.Rep, addr.Type, addr.String())
}

/*
UDPHeader is the header of an UDP request

	+----+------+------+----------+----------+----------+
	|RSV | FRAG | ATYP | DST.ADDR | DST.PORT |   DATA   |
	+----+------+------+----------+----------+----------+
	| 2  |  1   |  1   | Variable |    2     | Variable |
	+----+------+------+----------+----------+----------+
*/
type UDPHeader struct {
	Rsv  uint16
	Frag uint8
	Addr *Addr
}

// NewUDPHeader creates an UDPHeader
func NewUDPHeader(rsv uint16, frag uint8, addr *Addr) *UDPHeader {
	return &UDPHeader{
		Rsv:  rsv,
		Frag: frag,
		Addr: addr,
	}
}

func (h *UDPHeader) Write(w io.Writer) error {
	b := GetBuffer(smallSize)
	defer PutBuffer(b)

	binary.BigEndian.PutUint16(b[:2], h.Rsv)
	b[2] = h.Frag

	addr := h.Addr
	if addr == nil {
		addr = &Addr{}
	}
	length, _ := addr.Encode(b[3:])

	_, err := w.Write(b[:3+length])
	return err
}

func (h *UDPHeader) String() string {
	return fmt.Sprintf("%d %d %d %s",
		h.Rsv, h.Frag, h.Addr.Type, h.Addr.String())
}

// UDPDatagram represent an UDP request
type UDPDatagram struct {
	Header *UDPHeader
	Data   []byte
}

// NewUDPDatagram creates an UDPDatagram
func NewUDPDatagram(header *UDPHeader, data []byte) *UDPDatagram {
	return &UDPDatagram{
		Header: header,
		Data:   data,
	}
}

// ReadUDPDatagram reads an UDPDatagram from the stream
func ReadUDPDatagram(r io.Reader) (*UDPDatagram, error) {
	b := GetBuffer(largeSize)
	defer PutBuffer(b)

	// when r is a streaming (such as TCP connection), we may read more than the required data,
	// but we don't know how to handle it. So we use io.ReadFull to instead of io.ReadAtLeast
	// to make sure that no redundant data will be discarded.
	n, err := io.ReadFull(r, b[:5])
	if err != nil {
		return nil, err
	}

	header := &UDPHeader{
		Rsv:  binary.BigEndian.Uint16(b[:2]),
		Frag: b[2],
	}

	atype := b[3]
	hlen := 0
	switch atype {
	case AddrIPv4:
		hlen = 10
	case AddrIPv6:
		hlen = 22
	case AddrDomain:
		hlen = 7 + int(b[4])
	default:
		return nil, errors.New(xlog.NetworkAddrTypeError)
	}

	dlen := int(header.Rsv)
	if dlen == 0 { // standard SOCKS5 UDP datagram
		// Just read to the end of buffer
		nn, err := r.Read(b[n:])
		if err != nil && err != io.EOF {
			return nil, err
		}
		n += nn
		dlen = n - hlen
	} else { // extended feature, for UDP over TCP, using reserved field as data length
		if _, err := io.ReadFull(r, b[n:hlen+dlen]); err != nil {
			return nil, err
		}
		n = hlen + dlen
	}

	header.Addr = new(Addr)
	if err := header.Addr.Decode(b[3:hlen]); err != nil {
		return nil, err
	}

	data := make([]byte, dlen)
	copy(data, b[hlen:n])

	d := &UDPDatagram{
		Header: header,
		Data:   data,
	}

	return d, nil
}

func (d *UDPDatagram) Write(w io.Writer) error {
	h := d.Header
	if h == nil {
		h = &UDPHeader{}
	}
	buf := bytes.Buffer{}
	if err := h.Write(&buf); err != nil {
		return err
	}
	if _, err := buf.Write(d.Data); err != nil {
		return err
	}

	_, err := buf.WriteTo(w)
	return err
}
