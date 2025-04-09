package route

import (
	"net"

)

type Req struct {
	Method          string
	Url             string
	Protocol        string
	ProtocolVersion string
	Ahs             map[string]string
	Conn            net.Conn
    Body            []byte
}

func NewReq(method, url, protocol, protocolVersion string, ahs map[string]string, conn net.Conn) *Req {
    return &Req{
		Method:          method,
		Url:             url,
		Protocol:        protocol,
		ProtocolVersion: protocolVersion,
		Ahs:             ahs,
		Conn:            conn,
        Body:            nil,
	}
}

// read as text

// read json

// read byte array

