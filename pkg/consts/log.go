package consts

// Error Log
const (
	PROXY_ADDR_ERROR            string = "ProxyAddr Empty, Please Check Your Config File"
	CA_ERROR                    string = "CA Init Error"
	DEFAULT_HTTP                string = "SeaMoon Listening For The HTTP/S Request, %s"
	DEFAULT_SOCKS               string = "SeaMoon Listening For The SOCKS5 Request: %s"
	HTTP_LISTEN_ERROR           string = "[Http] Client Start Listen Error: %s"
	HTTP_ACCEPT_ERROR           string = "[Http] Server Accept Error: %s"
	SOCKS5_LISTEN_ERROR         string = "[Socks5] Client Start Listen Error: %s"
	SOCKS5_ACCEPT_ERROR         string = "[Socks5] Client Accept Error: %s"
	SOCKS5_READ_METHOD_ERROR    string = "[Socks5] Read Methods Failed: %s"
	SOCKS5_WRITE_METHOD_ERROR   string = "[Socks5] Write Method Failed: %s"
	SOCKS5_METHOD_ERROR         string = "[Socks5] Methods Is Not Acceptable"
	SOCKS5_READ_COMMAND_ERROR   string = "[Socks5] Read Command Failed: %s"
	SOCKS5_WRITE_COMMAND_ERROR  string = "[Socks5] Write Command Replay Failed: %s"
	SOCKS5_CONNECT_DIAL_ERROR   string = "[Socks5] Connect Dial Remote Failed: %s"
	SOCKS5_CONNECT_WRITE_ERROR  string = "[Socks5] Connect Write Reply Failed: %s"
	SOCKS5_CONNECT_TRANS_ERROR  string = "[Socks5] Connect Transport Failed: %s"
	SOCKS5_BIND_LISTEN_ERROR    string = "[Socks5] Bind Failed On Listen: %s"
	SOCKS5_BIND_WRITE_ERROR     string = "[Socks5] Bind Write Reply Failed %s"
	SOCKS5_BIND_ACCEPT_ERROR    string = "[Socks5] Bind Failed On Accept: %s"
	SOCKS5_BIND_TRANS_ERROR     string = "[Socks5] Bind Transport Failed: %s"
	SOCKS5_UDP2TCP_LISTEN_ERROR string = "[Socks5] Udp-Over-Tcp UDP Associate Failed On Listen: %s`"
	SOCKS5_UDP2TCP_WRITE_ERROR  string = "[Socks5] Udp-Over-Tcp Write Reply Failed %s`"
	SOCKS_UDP2TCP_UDP_ERROR     string = "[Socks5] Udp-Over-Tcp Tunnel UDP Failed: %s`"
	SOCKS_UPGRADE_ERROR         string = "[Socks5] WebSocket Upgrade error: %s "

	CLIENT_PROTOCOL_UNSUPPORT_ERROR string = "Protocol Not Support: %s"
)

// Info Log
const (
	CA_NOT_EXIST           string = "Ca Not Exists, Run Auto Generate"
	CA_LOAD_SUCCESS        string = "Ca Loaded Success"
	PROXY_ADDR             string = "Proxy Addr: %s"
	HTTP_LISTEN_START      string = "[Http] Client Start Listen At: %s"
	HTTP_LISTEN_STOP       string = "[Http] Client Stop Listen"
	HTTP_ACCEPT_START      string = "[Http] Server Accept Conn From: %s"
	HTTP_CONNECT_STOP_WAIT string = "[Http] Server Stopping Conn, Please Wait..."
	HTTP_BODY_DIS          string = "[Http] Server Conn Disconnected: %s"
	SOCKS5_LISTEN_START    string = "[Socks5] Client Start Listen At: %s"
	SOCKS5_LISTEN_STOP     string = "[Socks5] Client Stop Listen"
	SOCKS5_ACCEPT_START    string = "[Socks5] Client Accept Conn From: %s"
	SOCKS5_CONNECT_SERVER  string = "[Socks5] Server Connect %s For %s"
	SOCKS5_CONNECT_ESTAB   string = "[Socks5] Connect Tunnel Established %s <-> %s"
	SOCKS5_CONNECT_DIS     string = "[Socks5] Connect Tunnel Disconnected %s >-< %s"
	SOCKS5_BIND_SERVER     string = "[Socks5] Bind For %s"
	SOCKS5_BIND_ESTAB      string = "[Socks5] Bind Tunnel Established %s <-> %s"
	SOCKS5_BIND_DIS        string = "[Socks5] Bind Tunnel Disconnected %s >-< %s"
	SOCKS_UPD2TCP_SERVER   string = "[Socks5] Udp-Over-Tcp Associate UDP For %s"
	SOCKS5_UPD2TCP_ESTAB   string = "[Socks5] Udp-Over-Tcp Tunnel Established %s <-> (UDP)%s"
	SOCKS5_UPD2TCP_DIS     string = "[Socks5] Udp-Over-Tcp Tunnel Disconnected %s  >-< (UDP)%s"

	CONTROLLER_START string = "[Control] start control service at : %s"
)

// debug log
const (
	CONFIG_NOT_FIND      string = "[Config] Load Configuration Failed, Use Default"
	SOCKS5_UNSUPPORT_UDP string = "[Socks5] Unsupported Command CmdUDP"
)
