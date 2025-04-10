package route

import (
	"bytes"
	"context"
	"net"
	"strconv"
)

type Res struct {
	Protocol        string
	ProtocolVersion string
	Conn            net.Conn
	w               *bytes.Buffer
}

type ResWriteParam struct {
	Protocol        string
	ProtocolVersion string
	StatusCode      string
	Body            []byte
	Ahs             map[string]string
}

var (
	crlf       = []byte("\r\n")
	emptySpace = []byte(" ")
	colon      = []byte(":")
	slash      = []byte("/")
)

func NewRes(protocol, protocolVersion string, conn net.Conn) *Res {
	var w bytes.Buffer
	return &Res{
		Protocol:        protocol,
		ProtocolVersion: protocolVersion,
		Conn:            conn,
		w:               &w,
	}
}

func (r *Res) Write(ctx context.Context, param *ResWriteParam) {
	r.writeStartLine(param)
	r.writeHeader(param)
	r.w.Write(param.Body)

	r.Conn.Write(r.w.Bytes())
}

func (r *Res) writeStartLine(param *ResWriteParam) {
	r.w.WriteString(param.Protocol)
	r.w.Write(slash)
	r.w.WriteString(param.ProtocolVersion)
	r.w.Write(emptySpace)
	r.w.WriteString(param.StatusCode)
	r.w.Write(emptySpace)
	r.w.WriteString(codeMsgMap[param.StatusCode])
	r.w.Write(crlf)
}

func (r *Res) writeHeader(param *ResWriteParam) {
	ahs := param.Ahs
	for k, v := range ahs {
		r.w.WriteString(k)
		r.w.Write(colon)
		r.w.Write(emptySpace)
		r.w.WriteString(v)
		r.w.Write(crlf)
	}
	contentLength := len(param.Body)
	if contentLength > 0 {
		r.w.WriteString("Content-Length")
		r.w.Write(colon)
		r.w.Write(emptySpace)
		r.w.WriteString(strconv.Itoa(contentLength))
		r.w.Write(crlf)
	}
	r.w.Write(crlf)
}
