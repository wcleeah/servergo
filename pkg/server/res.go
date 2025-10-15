package server 

import (
	"bytes"
	"context"
	"io"
	"strconv"

	"lwc.com/servergo/internal/common"
)

type Res struct {
	Protocol        string
	ProtocolVersion string
	// using io.WriteCloser provides a few benefit
	// 1. testing will be easier because we don't need to create a fake connection
	// 2. req probably should not be reading from the connection (?)
	conn      io.WriteCloser
	w         *bytes.Buffer
	ctx       context.Context
	keepAlive bool
}

type ResWriteParam struct {
	StatusCode string
	Body       []byte
	Ahs        map[string]string
}

var (
	emptySpace = []byte(" ")
	colon      = []byte(":")
	slash      = []byte("/")
)

func newRes(ctx context.Context, protocol, protocolVersion string, keepAlive bool, conn io.WriteCloser) *Res {
	var w bytes.Buffer
	return &Res{
		Protocol:        protocol,
		ProtocolVersion: protocolVersion,
		conn:            conn,
		w:               &w,
		ctx:             ctx,
		keepAlive:       keepAlive,
	}
}

func (r *Res) Write(param *ResWriteParam) {
	r.writeStartLine(param)
	r.writeHeader(param)
	r.w.Write(param.Body)

	r.conn.Write(r.w.Bytes())
}

func (r *Res) writeStartLine(param *ResWriteParam) {
	r.w.WriteString(r.Protocol)
	r.w.Write(slash)
	r.w.WriteString(r.ProtocolVersion)
	r.w.Write(emptySpace)
	r.w.WriteString(param.StatusCode)
	r.w.Write(emptySpace)
	r.w.WriteString(codeMsgMap[param.StatusCode])
	r.w.Write(common.CRLF_BYTES)
}

func (r *Res) writeHeader(param *ResWriteParam) {
	ahs := param.Ahs
	for k, v := range ahs {
		r.w.WriteString(k)
		r.w.Write(colon)
		r.w.Write(emptySpace)
		r.w.WriteString(v)
		r.w.Write(common.CRLF_BYTES)
	}
	if _, ok := ahs["Content-Type"]; !ok {
		r.w.WriteString("Content-Type: text/plain")
		r.w.Write(common.CRLF_BYTES)
	}

    if r.keepAlive {
        r.w.WriteString("Connection: keep-alive")
        r.w.Write(common.CRLF_BYTES)
    } else {
        r.w.WriteString("Connection: close")
        r.w.Write(common.CRLF_BYTES)
	}

	contentLength := len(param.Body)
	if contentLength > 0 {
		r.w.WriteString("Content-Length")
		r.w.Write(colon)
		r.w.Write(emptySpace)
		r.w.WriteString(strconv.Itoa(contentLength))
		r.w.Write(common.CRLF_BYTES)
	}
	r.w.Write(common.CRLF_BYTES)
}
