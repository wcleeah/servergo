package common

import (
	"net"
)

type Listener interface {
	Close() error
	Accept() (net.Conn, error)
}
