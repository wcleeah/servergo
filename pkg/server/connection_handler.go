package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"lwc.com/servergo/internal/http"
	"lwc.com/servergo/internal/logger"
)

type connHandler struct {
	conn      io.ReadWriteCloser
	keepAlive bool
	req       *Req
	res       *Res
	bufReader *bufio.Reader
	ctx       context.Context
}

func newConnHandler(ctx context.Context) *connHandler {
	return &connHandler{
		keepAlive: true,
		ctx:       ctx,
	}
}

// TODO: timeout case
func (ch *connHandler) Handle(conn io.ReadWriteCloser) {
	l := logger.Get(ch.ctx)
	l.Info("New Connection")
	ch.conn = conn
	ch.bufReader = bufio.NewReader(ch.conn)
	defer ch.conn.Close()

	var err error
	for ch.keepAlive {
		ch.req = nil
		err = ch.handle()
		if err != nil {
			l.Error("Error, breaking")
			break
		}
	}

	l.Info("Closing connection... handling error if any")

	if errors.Is(err, io.EOF) {
		l.Info("io.EOF, nothing we can do for an EOF since the connection is already closed by client side")
		return
	}

	if ch.res == nil {
		ch.res = newRes(ch.ctx, http.SUPPORTED_PROTOCOL, http.DEFAULT_PROTOCOL_VERSION, false, ch.conn)
	}

	if errors.Is(err, http.UnsupportedMethod) {
		l.Error("unsupportedMethod")
		ch.res.Write(&ResWriteParam{
			StatusCode: "405",
		})
		return
	}

	if errors.Is(err, http.UnsupportedProtocolVersion) {
		l.Error("unsupportedProtocolVersion")
		ch.res.Write(&ResWriteParam{
			StatusCode: "505",
		})
		return
	}

	if err != nil {

		l.Error("other errors", "error", err.Error())
		ch.res.Write(&ResWriteParam{
			StatusCode: "400",
			Body:       []byte(err.Error()),
		})
		return
	}

	l.Info("Connection Closed")
}

func (ch *connHandler) handle() error {
	l := logger.Get(ch.ctx)

	startLineBytes := make([]byte, 0)

	headerBytes := make([]byte, 0)
	ahs := make(map[string]string, 0)

	// Notes:
	// - originally, i was using prime's http course method, with a state machine
	// - but if allBytes is big enough, the state machine will be blocked by the second read as all bytes all being read in the first call already
	// - buio.Reader header reading might acceidentally read some part of body
	// - i want a solution that can handle reading any number of bytes a time
	// Solution:
	// - use bufio.ReadLine, which is what golang stdlib does, this solve both problems
	//  - bufio.Reader can specify buf size, so it will read til the line
	//  - ReadLine already handles buf size issue
	// - the go internal package use textProto.Reader to do read line
	// - it handle the isPrefix logic inside textProto.Reader, here i just lay it out instead of using an extra struct
	for {

		// ReadLine will read until:
		// - it got the terminator "\n"
		// - or the read bytes size exceeds the size of the buffer, in this case it is bytesLimit
		line, isPrefix, err := ch.bufReader.ReadLine()
		if err != nil {
			return err
		}

		startLineBytes = append(startLineBytes, line...)
		if isPrefix {
			continue
		}
		break
	}

	method, url, protocolVersion, protocol, err := http.ReadStartLine(ch.ctx, startLineBytes)
	if err != nil {
		return err
	}

	for {
		line, isPrefix, err := ch.bufReader.ReadLine()

		if err != nil {
			return err
		}

		headerBytes = append(headerBytes, line...)

		l.Info("bytes", "headerBytes", string(headerBytes), "line", string(line))
		if isPrefix {
			continue
		}

		key, value, noMoreHeader, err := http.ReadHeader(ch.ctx, headerBytes)
		if err != nil {
			return err
		}
		if noMoreHeader {
			break
		}

		if v, ok := ahs[key]; ok {
			ahs[key] = v + ", " + value
		} else {
			ahs[key] = value
		}

		headerBytes = make([]byte, 0)
	}

	if v, ok := ahs["connection"]; ok {
		// in http/1.1, default behaviour is to keep the connection alive
		// only close it when the header specifies close
		// it does not state how to handle if multiple values are sent with the Connection header
		// i will just take precendence on the Close value here
		ch.keepAlive = !strings.Contains(v, "Close")
	}

	if v, ok := ahs["content-length"]; ok {
		_, err := strconv.Atoi(v)
		if err != nil {
			return errors.New("Content-Length is malformed")
		}
	}

	l.Info("before routing info", "startLine", fmt.Sprintf("%s %s, %s/%s", method, url, protocol, protocolVersion), "headers", fmt.Sprintf("%+v", ahs), "keep alive", ch.keepAlive)

	ch.req = newReq(ch.ctx, method, url, protocol, protocolVersion, ahs, ch.bufReader)
	ch.res = newRes(ch.ctx, protocol, protocolVersion, ch.keepAlive, ch.conn)

	route(ch.req, ch.res)

	// since the next request will be using the same connection
	// new bytes from the new request cannot be read unless the prev body bytes is read
	// offtopic: https://github.com/golang/go/issues/60240
	// Open source maintainence is really hard
	if !ch.req.IsBodyRead() && ch.keepAlive {
		// if the request is processed, but reading body throws error
		// will just let the error floats and closing the connection
		// not the server responsibility to handle timing difference between body reading and handler process imo
		_, err := io.ReadAll(ch.req.Body())
		if err != nil {
			return err
		}
	}

	return nil
}
