package network

type Network string

const (
	HTTP       Network = "http"
	TCP                = "tcp"
	TCP4               = "tcp4"
	TCP6               = "tcp6"
	UDP                = "udp"
	UDP4               = "udp4"
	UDP6               = "udp6"
	IP                 = "ip"
	IP4                = "ip4"
	IP6                = "ip6"
	UNIX               = "unix"
	UNIXGRAM           = "unixgram"
	UNIXPACKET         = "unixpacket"
)

func (n Network) String() string {
	return string(n)
}

func SupportedNetwork(network Network) bool {
	ns := map[Network]bool{
		HTTP:       true,
		TCP:        true,
		TCP4:       true,
		TCP6:       true,
		UDP:        true,
		UDP4:       true,
		UDP6:       true,
		IP:         true,
		IP4:        true,
		IP6:        true,
		UNIX:       true,
		UNIXGRAM:   true,
		UNIXPACKET: true}
	return ns[network]
}
