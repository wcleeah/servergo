package http

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type VerySlowConn struct {
	Limit        int
	SendBytes    []byte
	ReceiveBytes []byte
	cur          int
	Closed       bool
	Recursive    bool
}

// assume target will always be larger then limit
func (vsc VerySlowConn) Read(target []byte) (int, error) {
	if vsc.Closed {
		return 0, io.EOF
	}
	start := vsc.cur
	end := vsc.cur + vsc.Limit
	if end > len(vsc.SendBytes) {
		rem := len(vsc.SendBytes) - end
		end -= rem
		send := vsc.SendBytes[start:end]
		start = 0
		end = rem
		send = append(send, vsc.SendBytes[start:end]...)

		vsc.setCur(end)
		copy(target, send)
		return vsc.Limit, nil
	}
	send := vsc.SendBytes[vsc.cur : vsc.cur+vsc.Limit]
	vsc.setCur(vsc.cur + vsc.Limit)

	copy(target, send)
	return vsc.Limit, nil
}

func (vsc VerySlowConn) Write(p []byte) (int, error) {
	if vsc.Closed {
		return 0, io.EOF
	}
	vsc.setReceiveBytes(append(vsc.ReceiveBytes, p...))

	return len(p), nil
}

func (vsc VerySlowConn) Close() error {
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
	vsc.ReceiveBytes = bytes
}

// in all the following case, i will probably not check the content,
// those are individual reading function test's job
// will focus more on if the req is read completely
// will i get a response anyways
// and is the connection handled properly
func TestSmallPacketsHandling(t *testing.T) {
	requestStr := "GET /bro http/1.1\r\nNoway: hahahaha\r\nConnection: Close\r\n\r\n\r\n"
	conn := VerySlowConn{
		Limit:     8,
		SendBytes: []byte(requestStr),
		Closed:    false,
		Recursive: true,
	}

	handler := NewConnHandler()
	handler.Handle(context.Background(), conn)

	assert.Equal(t, stateHeaderDone, handler.state)
	assert.Equal(t, "GET", handler.req.Method)
	assert.Equal(t, "/bro", handler.req.Url)
	assert.Equal(t, "http", handler.req.Protocol)
	assert.Equal(t, "1.1", handler.req.ProtocolVersion)
	assert.Equal(t, "hahahaha", handler.req.GetHeader("Noway"))
	assert.Equal(t, "Close", handler.req.GetHeader("Connection"))
	assert.Equal(t, true, handler.req.IsBodyRead())
}
func TestBadRequest(t *testing.T) {}
func TestKeepAlive(t *testing.T)  {}
func TestClose(t *testing.T)      {}
func TestEOF(t *testing.T)        {}
func TestCleanUpBody(t *testing.T)        {}
