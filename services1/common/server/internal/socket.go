package internal

type Socket string

const socketSuffix = ".socket"

func NewSocket(name string) Socket {
	return Socket(name)
}

func (s Socket) ListenAddress() string {
	return (string)(s) + socketSuffix
}
