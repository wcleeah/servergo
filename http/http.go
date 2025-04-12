package http

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net"

	"lwc.com/servergo/route"
)

func handleConn(ctx context.Context, conn net.Conn) error {
	bufIoReader := bufio.NewReader(conn)

	sls, err := bufIoReader.ReadString('\n')
	if err != nil {
		return err
	}

	startLine, err := readStartLine(ctx, sls)
	if err != nil {
		return err
	}

	ahs := make(map[string]string, 0)

	for {
		headerLine, err := bufIoReader.ReadString('\n')
		if err != nil {
			return err
		}
		key, value, err := readHeader(ctx, headerLine)
		if err != nil {
			if errors.Is(err, headerEnds) {
				break
			}
			return err
		}
		ahs[key] = value
	}

	req := route.NewReq(ctx, startLine.Method, startLine.Url, startLine.Protocol, startLine.ProtocolVersion, ahs, conn)
	res := route.NewRes(ctx, startLine.Protocol, startLine.ProtocolVersion, conn)
    err = route.Route(req, res)
    return err 
}

func HandleConn(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	defaultRes := route.NewRes(ctx, SUPPORTED_PROTOCOL, SUPPORTED_PROTOCOL_VERSION[0], conn)
	err := handleConn(ctx, conn)

	if errors.Is(err, io.EOF) {
		defaultRes.Write(&route.ResWriteParam{
			StatusCode: "400",
			Body:       []byte("Request malformed, unexpected EOF"),
		})
		return
	}

    if errors.Is(err, unsupportedMethod) {
		defaultRes.Write(&route.ResWriteParam{
			StatusCode: "405",
			Body:       []byte("Request malformed, unexpected EOF"),
		})
		return
    }

    if errors.Is(err, unsupportedProtocolVersion) {
		defaultRes.Write(&route.ResWriteParam{
			StatusCode: "405",
			Body:       []byte("Request malformed, unexpected EOF"),
		})
		return
    }

	if err != nil {
		defaultRes.Write(&route.ResWriteParam{
			StatusCode: "400",
			Body:       []byte(err.Error()),
		})
		return
	}
}
