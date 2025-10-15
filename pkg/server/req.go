package server 

import (
	"bufio"
	"context"
	"io"
	"strconv"
	"strings"
	"sync"

	"lwc.com/servergo/internal/http"
)

type Req struct {
	Method          string
	Url             string
	Protocol        string
	ProtocolVersion string
	ahs             map[string]string
	ctx             context.Context
	mu              sync.Mutex
	body            *http.Body
}

func newReq(ctx context.Context, method, url, protocol, protocolVersion string, ahs map[string]string, connReader *bufio.Reader) *Req {
	req := &Req{
		Method:          method,
		Url:             url,
		Protocol:        protocol,
		ProtocolVersion: protocolVersion,
		ahs:             ahs,
		ctx:             ctx,
	}

	// if content length is not provided, treat it as no body and pass zero-value
	clStr, ok := req.GetHeader("Content-Length")
	var clInt int
	if ok {
		clInt, _ = strconv.Atoi(clStr)
	}

	req.body = http.NewBody(connReader, clInt)
	return req
}

// we will assume the body length must be the same with content Content-Length header
func (r *Req) Body() io.Reader {
	return r.body
}

func (r *Req) GetHeader(key string) (string, bool) {
	h, ok := r.ahs[strings.ToLower(key)]
	return h, ok
}

func (r *Req) Ctx() context.Context {
	if r.ctx == nil {
		r.ctx = context.Background()
	}
	return r.ctx
}

func (r *Req) IsBodyRead() bool {
	return r.body.IsBodyRead()
}
