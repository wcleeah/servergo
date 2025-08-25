package http

import (
	"bufio"
	"context"
	"errors"
	"io"
	"strings"

	"lwc.com/servergo/logger"
	"lwc.com/servergo/route"
)

type State string

const (
	bytesLimit int = 3 
)

type ConnHandler struct {
	conn      io.ReadWriteCloser
	keepAlive bool
	req       *route.Req
	res       *route.Res
	bufReader *bufio.Reader
}

func NewConnHandler() *ConnHandler {
	return &ConnHandler{
		keepAlive: true,
	}
}

// TODO: timeout case
func (ch *ConnHandler) Handle(ctx context.Context, conn io.ReadWriteCloser) {
	l := logger.Get(ctx)
	l.Info("New Connection")
	ch.conn = conn
	ch.bufReader = bufio.NewReader(ch.conn)
	defer ch.conn.Close()

	var err error
	for ch.keepAlive {
		ch.req = nil
		ch.res = route.NewRes(ctx, SUPPORTED_PROTOCOL, SUPPORTED_PROTOCOL_VERSION[0], false, ch.conn)
		err = ch.listen(ctx)
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

	if errors.Is(err, unsupportedMethod) {
		l.Error("unsupportedMethod")
		ch.res.Write(&route.ResWriteParam{
			StatusCode: "405",
		})
		return
	}

	if errors.Is(err, unsupportedProtocolVersion) {
		l.Error("unsupportedProtocolVersion")
		ch.res.Write(&route.ResWriteParam{
			StatusCode: "505",
		})
		return
	}

	if err != nil {
		l.Error("other errors", "error", err.Error())
		ch.res.Write(&route.ResWriteParam{
			StatusCode: "400",
			Body:       []byte(err.Error()),
		})
		return
	}

	l.Info("Connection Closed")
}

func (ch *ConnHandler) listen(ctx context.Context) error {

	startLineBytes := make([]byte, 0)
	var startLine *StartLine

	headerBytes := make([]byte, 0)
	ahs := make(map[string]string, 0)


	// BUGS: (fixed)
	// - if allBytes is big enough, the state machine will be blocked by the second read as all bytes all being read in the first call already
	// - header reading might acceidentally read some part of body
	// Solution:
	// - add a new byte array specifically for storing read result, the array is not fixed length, so can check the length
	// - ⭐ or, use bufio.ReadLine, which is what golang stdlib does, this solve both problems
	//  - bufio.Reader can specify buf size, so it will read til the line
	//  - ReadLine already handles buf size issue
	for {
		// the go internal package use read line also but it ignores isPrefix
		// it assumes either the returned line has the whole startline, or it is too big
		// i am not implementing size limit here, adding a inf for loop to handle isPrefix = true case
		// (i probably should)

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

		startLine, err = readStartLine(ctx, startLineBytes)
		if err != nil {
			return err
		}
		break

	}

	for {
		line, isPrefix, err := ch.bufReader.ReadLine()

		if err != nil {
			return err
		}

		headerBytes = append(headerBytes, line...)
		if isPrefix {
			continue
		}

		key, value, err := readHeader(ctx, headerBytes)
		if err != nil {
			if errors.Is(err, headerEnds) {
				break
			}
			return err
		}

		if v, ok := ahs[key]; ok {
			ahs[key] = v + ", " + value
		} else {
			ahs[key] = value
		}


		headerBytes = make([]byte, 0)
	}

	if v, ok := ahs["Connection"]; ok {
		// in http/1.1, default behaviour is to keep the connection alive
		// only close it when the header specifies close
		// it does not state how to handle if multiple values are sent with the Connection header
		// i will just take precendence on the Close value here
		ch.keepAlive = strings.Contains(v, "Close")
	}

	ch.req = route.NewReq(ctx, startLine.Method, startLine.Url, startLine.Protocol, startLine.ProtocolVersion, ahs, ch.conn, ch.bufReader)
	ch.res = route.NewRes(ctx, startLine.Protocol, startLine.ProtocolVersion, ch.keepAlive, ch.conn)
	route.Route(ch.req, ch.res)

	// since the next request will be using the same connection
	// new bytes from the new request cannot be read unless the prev body bytes is read
	// offtopic: https://github.com/golang/go/issues/60240
	// Open source maintainence is really hard
	if !ch.req.IsBodyRead() && ch.keepAlive {
		// if the request is processed, but reading body throws error
		// will just let the error floats and closing the connection
		// not the server responsibility to handle timing difference between body reading and handler process imo
		return ch.req.CleanUpBodyBytes()
	}

	return nil
}
