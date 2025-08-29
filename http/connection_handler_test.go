package http

import (
	"context"
	"io"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"lwc.com/servergo/logger"
	"lwc.com/servergo/route"
)

type VerySlowConn struct {
	Limit         int
	SentBytes     []byte
	ReceivedBytes []byte
	ReceivedCount int
	cur           int
	Closed        bool
	Recursive     bool
}

// assume target will always be larger then limit
func (vsc *VerySlowConn) Read(target []byte) (int, error) {
	if vsc.Closed || vsc.cur >= len(vsc.SentBytes) {
		return 0, io.EOF
	}
	start := vsc.cur
	end := min(vsc.cur + vsc.Limit, len(vsc.SentBytes))
	send := vsc.SentBytes[start:end]
	vsc.setCur(end)

	copy(target, send)
	return vsc.Limit, nil
}

func (vsc *VerySlowConn) Write(p []byte) (int, error) {
	if vsc.Closed {
		return 0, io.EOF
	}
	vsc.ReceivedCount++
	vsc.setReceiveBytes(append(vsc.ReceivedBytes, p...))

	return len(p), nil
}

func (vsc *VerySlowConn) Close() error {
	vsc.setClosed(true)
	return nil
}

func (vsc *VerySlowConn) setCur(newCur int) {
	vsc.cur = newCur
}

func (vsc *VerySlowConn) setClosed(closed bool) {
	vsc.Closed = closed
}

func (vsc *VerySlowConn) setReceiveBytes(bytes []byte) {
	vsc.ReceivedBytes = bytes
}

type ReadWriteCloser struct {
	Conn *VerySlowConn
}

func (rwc ReadWriteCloser) Read(bs []byte) (int, error) {
	return rwc.Conn.Read(bs)
}

func (rwc ReadWriteCloser) Write(bs []byte) (int, error) {
	return rwc.Conn.Write(bs)
}

func (rwc ReadWriteCloser) Close() error {
	return rwc.Conn.Close()
}

// in all the following case, i will probably not check the content,
// those are individual reading function test's job
// will focus more on if the req is read completely
// will i get a response anyways
// and is the connection handled properly

// continuous small packet
func TestSmallPacketsHandling(t *testing.T) {
	requestStr := "GET /bro HTTP/1.1\r\nNoway: hahahaha\r\nConnection: Close\r\n\r\n"
	conn := &VerySlowConn{
		Limit:     3,
		SentBytes: []byte(requestStr),
		Closed:    false,
		Recursive: true,
	}

	rwc := ReadWriteCloser{
		Conn: conn,
	}

	ctx := context.WithValue(context.Background(), logger.TRACE_ID_KEY, uuid.NewString())
	handler := NewConnHandler(ctx)

	// Connection: Close should terminate the inf loop
	handler.Handle(rwc)

	assert.NotNil(t, handler.req)
	assert.NotNil(t, handler.res)
	assert.Equal(t, "GET", handler.req.Method)
	assert.Equal(t, "/bro", handler.req.Url)
	assert.Equal(t, "HTTP", handler.req.Protocol)
	assert.Equal(t, "1.1", handler.req.ProtocolVersion)
	assert.Equal(t, "hahahaha", handler.req.GetHeader("Noway"))
	assert.Equal(t, "Close", handler.req.GetHeader("Connection"))
	assert.True(t, conn.Closed)
}

// connection reuse and connection close
func TestKeepAlive(t *testing.T) {
	requestStr := "GET /bro HTTP/1.1\r\nNoway: hahahaha\r\n\r\nGET /closer HTTP/1.1\r\nConnection: Close\r\nNoway: hahahaha\r\n\r\n"
	conn := &VerySlowConn{
		Limit:     8,
		SentBytes: []byte(requestStr),
		Closed:    false,
		Recursive: true,
	}

	rwc := ReadWriteCloser{
		Conn: conn,
	}

	traceId := uuid.NewString()
	ctx := context.WithValue(context.Background(), logger.TRACE_ID_KEY, traceId)
	handler := NewConnHandler(ctx)

	route.AddRoute("GET /bro", func(req *route.Req, res *route.Res) {
		assert.True(t, handler.keepAlive)
		assert.Equal(t, traceId, req.Ctx().Value(logger.TRACE_ID_KEY))
		res.Write(&route.ResWriteParam{
			StatusCode: "200",
		})
	})

	route.AddRoute("GET /closer", func(req *route.Req, res *route.Res) {
		assert.False(t, handler.keepAlive)
		assert.Equal(t, traceId, req.Ctx().Value(logger.TRACE_ID_KEY))
		res.Write(&route.ResWriteParam{
			StatusCode: "200",
		})
	})

	handler.Handle(rwc)

	assert.Equal(t, 2, conn.ReceivedCount)
	assert.False(t, handler.keepAlive)
}

// if the body is not read by handler, it should clean it up
// three cases, all of them should be successfully proceeded:
// - Have Body, handler did read
// - Have Body, handler did not read
// - No Body, handler did not read
func TestBodyRead(t *testing.T) {
	requestStr := "POST /readBody HTTP/1.1\r\nNoway: hahahaha\r\nContent-Length: 9\r\n\r\n123123123POST /noReadBody HTTP/1.1\r\nNoway: hahahaha\r\nContent-Length: 9\r\n\r\n123123123POST /noReadBody HTTP/1.1\r\nConnection: Close\r\nNoway: hahahaha\r\n\r\n"
	conn := &VerySlowConn{
		Limit:     8,
		SentBytes: []byte(requestStr),
		Closed:    false,
		Recursive: true,
	}

	rwc := ReadWriteCloser{
		Conn: conn,
	}

	traceId := uuid.NewString()
	ctx := context.WithValue(context.Background(), logger.TRACE_ID_KEY, traceId)
	handler := NewConnHandler(ctx)

	route.AddRoute("POST /readBody", func(req *route.Req, res *route.Res) {
		_, err := io.ReadAll(req.Body())
		assert.NoError(t, err)

		res.Write(&route.ResWriteParam{
			StatusCode: "200",
		})
	})

	route.AddRoute("POST /noReadBody", func(req *route.Req, res *route.Res) {
		res.Write(&route.ResWriteParam{
			StatusCode: "200",
		})
	})

	handler.Handle(rwc)

	assert.Equal(t, 3, conn.ReceivedCount)
}

// EOF
func TestEOF(t *testing.T) {
	// trigger an EOF immediately, readLine will read til the EOF, and not return error(?)
	requestStr := ""
	conn := &VerySlowConn{
		Limit:     8,
		SentBytes: []byte(requestStr),
		Closed:    false,
		Recursive: true,
	}

	rwc := ReadWriteCloser{
		Conn: conn,
	}

	traceId := uuid.NewString()
	ctx := context.WithValue(context.Background(), logger.TRACE_ID_KEY, traceId)
	handler := NewConnHandler(ctx)

	handler.Handle(rwc)

	assert.Equal(t, 0, conn.ReceivedCount)
	assert.True(t, conn.Closed)
}
