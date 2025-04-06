package http

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net"
	"strconv"

	"lwc.com/servergo/logger"
)

func HandleConn(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	l := logger.Get(ctx)
	bufIoReader := bufio.NewReader(conn)

	sls, err := bufIoReader.ReadString('\n')
	if err != nil {
		l.Error("Startline read failed")
		return
	}

	startLine, err := readStartLine(ctx, sls)
	if err != nil {
		l.Error("Start Line read error", "err", err.Error())
		return
	}

	l.Info("Start Line", "sl", startLine)
    ahs := make(map[string]string, 0)

	for {
		headerLine, err := bufIoReader.ReadString('\n')
		if err != nil {
			l.Error("headerLine read failed")
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

    cl, ok := ahs["Content-Length"]
    if !ok {
        l.Info("Request has no body")
        return
    }

    clInt, err := strconv.Atoi(cl)
    if err != nil {
        l.Error("Content length to integer has error", "err", err.Error())
        return
    }

    bodySlice := make([]byte, clInt)
    _, err = io.ReadFull(bufIoReader, bodySlice)
    if err != nil {
        l.Error("Body read faild", "err", err.Error())
        return
    }

    bodyStr, err := readBody(bodySlice)
    if err != nil {
        l.Error("Body read error", "err", err.Error())
        return
    }
    l.Info("Body", "body", bodyStr)
}
