package xlog

// Info Log

const (
	CONTROLLER_START string = "start control api service at"
)

// Error Log
const (
	PROXY_ADDR_ERROR            string = "ProxyAddr Empty, Please Check Your Config File"
	CA_ERROR                    string = "CA Init Error"
	DEFAULT_HTTP                string = "SeaMoon Listening For The HTTP/S Request"
	DEFAULT_SOCKS               string = "SeaMoon Listening For The SOCKS5 Request"
	LISTEN_ERROR                string = "Client Start Listen Error"
	LISTEN_STOP_ERROR           string = "Client Stop Listen Error"
	HTTP_ACCEPT_ERROR           string = "[Http] Server Accept Error"
	SOCKS5_LISTEN_ERROR         string = "[Socks5] Client Start Listen Error"
	ACCEPT_ERROR                string = "Client Accept Error"
	SOCKS5_READ_METHOD_ERROR    string = "[Socks5] Read Methods Failed"
	SOCKS5_WRITE_METHOD_ERROR   string = "[Socks5] Write Method Failed"
	SOCKS5_METHOD_ERROR         string = "[Socks5] Methods Is Not Acceptable"
	SOCKS5_READ_COMMAND_ERROR   string = "[Socks5] Read Command Failed"
	SOCKS5_WRITE_COMMAND_ERROR  string = "[Socks5] Write Command Replay Failed"
	SOCKS5_CONNECT_DIAL_ERROR   string = "[Socks5] Connect Dial Remote Failed"
	SOCKS5_CONNECT_WRITE_ERROR  string = "[Socks5] Connect Write Reply Failed"
	CONNECT_RMOET_ERROR         string = "Connect Remote Failed"
	CONNECT_TRANS_ERROR         string = "Connect Transport Failed"
	SOCKS5_BIND_LISTEN_ERROR    string = "[Socks5] Bind Failed On Listen"
	SOCKS5_BIND_WRITE_ERROR     string = "[Socks5] Bind Write Reply Failed"
	SOCKS5_BIND_ACCEPT_ERROR    string = "[Socks5] Bind Failed On Accept"
	SOCKS5_BIND_TRANS_ERROR     string = "[Socks5] Bind Transport Failed"
	SOCKS5_UDP2TCP_LISTEN_ERROR string = "[Socks5] Udp-Over-Tcp UDP Associate Failed On Listen"
	SOCKS5_UDP2TCP_WRITE_ERROR  string = "[Socks5] Udp-Over-Tcp Write Reply Failed "
	SOCKS_UDP2TCP_UDP_ERROR     string = "[Socks5] Udp-Over-Tcp Tunnel UDP Failed"
	SOCKS_UPGRADE_ERROR         string = "[Socks5] WebSocket Upgrade error:"

	FORWARD_ACTION_EMPTY string = "[Forward] Empty ACTION"

	CLIENT_PROTOCOL_UNSUPPORT_ERROR string = "Protocol Not Support"
)

// Info Log
const (
	CA_NOT_EXIST           string = "Ca Not Exists, Run Auto Generate"
	CA_LOAD_SUCCESS        string = "Ca Loaded Success"
	PROXY_ADDR             string = "Proxy Addr"
	LISTEN_START           string = "Client Start Listen At"
	LISTEN_STOP            string = "Client Stop Listen"
	HTTP_ACCEPT_START      string = "[Http] Server Accept Conn From"
	HTTP_CONNECT_STOP_WAIT string = "[Http] Server Stopping Conn, Please Wait..."
	HTTP_BODY_DIS          string = "[Http] Server Conn Disconnected"
	SOCKS5_LISTEN_START    string = "[Socks5] Client Start Listen At"
	SOCKS5_LISTEN_STOP     string = "[Socks5] Client Stop Listen"
	SOCKS5_ACCEPT_START    string = "[Socks5] Client Accept Conn From"
	SOCKS5_CONNECT_SERVER  string = "[Socks5] Server Connect"
	SOCKS5_CONNECT_ESTAB   string = "[Socks5] Connect Tunnel Established"
	SOCKS5_CONNECT_DIS     string = "[Socks5] Connect Tunnel Disconnected"
	SOCKS5_BIND_SERVER     string = "[Socks5] Bind For"
	SOCKS5_BIND_ESTAB      string = "[Socks5] Bind Tunnel Established"
	SOCKS5_BIND_DIS        string = "[Socks5] Bind Tunnel Disconnected"
	SOCKS_UPD2TCP_SERVER   string = "[Socks5] Udp-Over-Tcp Associate UDP For"
	SOCKS5_UPD2TCP_ESTAB   string = "[Socks5] Udp-Over-Tcp Tunnel Established"
	SOCKS5_UPD2TCP_DIS     string = "[Socks5] Udp-Over-Tcp Tunnel Disconnected"
)
