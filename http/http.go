package http

import (
	"bufio"
	"context"
	"errors"
	"net"

	"lwc.com/servergo/logger"
	"lwc.com/servergo/route"
)

func HandleConn(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	l := logger.Get(ctx)
	bufIoReader := bufio.NewReader(conn)

	l.Info("Reading start line")
	sls, err := bufIoReader.ReadString('\n')
	if err != nil {
        // write bad request
		l.Error("Startline read failed")
		return
	}

	startLine, err := readStartLine(ctx, sls)
	if err != nil {
		l.Error("Start Line read error", "err", err.Error())
        // write bad request
		return
	}

	l.Info("Start Line", "sl", startLine)

	l.Info("Reading headers")
	ahs := make(map[string]string, 0)

	for {
		headerLine, err := bufIoReader.ReadString('\n')
		if err != nil {
			l.Error("headerLine read failed")
            // write bad request
			return
		}
		key, value, err := readHeader(ctx, headerLine)
		if err != nil {
			if errors.Is(err, headerEOF) {
				l.Info("Header section ended")
				break
			}
			l.Error("Header Line read error", "err", err.Error())
			return
		}
		ahs[key] = value
	}

	l.Info("All Header", "header", ahs)
    req := route.NewReq(startLine.Method, startLine.Url, startLine.Protocol, startLine.ProtocolVersion, ahs, conn)
    res := route.NewRes(startLine.Protocol, startLine.ProtocolVersion, conn)
    route.Route(ctx, req, res)
}
