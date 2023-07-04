package service

type Transport string

const (
	Tcp  Transport = "tcp"
	Unix           = "unix"
)

func (t Transport) String() string {
	return (string)(t)
}

type TcpTransport struct{}

func (p *TcpTransport) ServiceTransport() Transport {
	return Tcp
}

type UnixTransport struct{}

func (u *UnixTransport) ServiceTransport() Transport {
	return Unix
}
