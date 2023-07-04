package internal

import "net"

type Network string

const (
	TcpNetwork  Network = "tcp"
	UnixNetwork Network = "unix"
)

func (n Network) String() string {
	return (string)(n)
}

func NewTcpListener(port Port) (net.Listener, error) {
	return net.Listen(TcpNetwork.String(), port.ListenAddress())
}

func NewUnixListener(socket Socket) (net.Listener, error) {
	return net.Listen(UnixNetwork.String(), socket.ListenAddress())
}
