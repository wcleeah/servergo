package http

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net"
	"os"
	"time"

	"lwc.com/servergo/logger"
	"lwc.com/servergo/route"
)

func handleConn(ctx context.Context, conn net.Conn, bufIoReader *bufio.Reader) (bool, error) {
    l := logger.Get(ctx)

	sls, err := bufIoReader.ReadString('\n')
	if err != nil {
        l.Info("hahahahahahahahahahahahahahah")
		return false, err
	}

	startLine, err := readStartLine(ctx, sls)
	if err != nil {
		return false, err
	}

	ahs := make(map[string]string, 0)
	keepAlive := false

	for {
		headerLine, err := bufIoReader.ReadString('\n')
		if err != nil {
			return false, err
		}
		key, value, err := readHeader(ctx, headerLine)
		if err != nil {
			if errors.Is(err, headerEnds) {
				break
			}
			return false, err
		}
		ahs[key] = value
		if key == "Connection" {
			keepAlive = value != "Close"
		}
	}

	req := route.NewReq(ctx, startLine.Method, startLine.Url, startLine.Protocol, startLine.ProtocolVersion, ahs, conn)
	res := route.NewRes(ctx, startLine.Protocol, startLine.ProtocolVersion, keepAlive, conn)
	err = route.Route(req, res)
	return keepAlive, err
}

func HandleConn(ctx context.Context, conn net.Conn) {
	l := logger.Get(ctx)
	l.Info("New Connection")
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(time.Second * 10))
	bufIoReader := bufio.NewReader(conn)
	for {
		defaultRes := route.NewRes(ctx, SUPPORTED_PROTOCOL, SUPPORTED_PROTOCOL_VERSION[0], false, conn)
		keepAlive, err := handleConn(ctx, conn, bufIoReader)

		// lets assume if the connection is expired, we dun keep it open for now
		// instead the client will reconnect
		// context: net.Conn can be reused even if the connection is expired
		if errors.Is(err, os.ErrDeadlineExceeded) {
			l.Info("Connection Expired")
			break
		}

		if errors.Is(err, io.EOF) {
            l.Info("1")
			defaultRes.Write(&route.ResWriteParam{
				StatusCode: "400",
				Body:       []byte("Request malformed, unexpected EOF"),
			})
			break
		}

		if errors.Is(err, unsupportedMethod) {
			defaultRes.Write(&route.ResWriteParam{
				StatusCode: "405",
				Body:       []byte("Request malformed, unexpected EOF"),
			})
			break
		}

		if errors.Is(err, unsupportedProtocolVersion) {
			defaultRes.Write(&route.ResWriteParam{
				StatusCode: "405",
				Body:       []byte("Request malformed, unexpected EOF"),
			})
			break
		}

		if err != nil {
			defaultRes.Write(&route.ResWriteParam{
				StatusCode: "400",
				Body:       []byte(err.Error()),
			})
			break
		}

		if !keepAlive {
			l.Info("Connection Closed By Client")
			break
		}

		conn.SetDeadline(time.Now().Add(time.Second * 10))
	}

	l.Info("Connection Closed for whatever reason")
}
