package server

import (
	"fmt"
	"net"
)

type (
	network string
	port    string
)

type NetWork struct {
	network
	port
}

func New(n network, p port) NetWork {
	return NetWork{
		n, p,
	}
}

func (n NetWork) Run() error {
	listner, listnerErr := net.Listen(string(n.network), string(n.port))
	if listnerErr != nil {
		return fmt.Errorf("error during %v server init %v", n.network, listnerErr)
	}
	defer listner.Close()

	for {
		con, conErr := listner.Accept()
		if conErr != nil {
			// do nothing for now
		}
		go handle(con)
	}
}
