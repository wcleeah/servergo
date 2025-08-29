package route

import (
	"bufio"
	"context"
	"io"
	"strconv"
	"strings"
	"sync"
)

type Req struct {
	Method          string
	Url             string
	Protocol        string
	ProtocolVersion string
	ahs             map[string]string
	ctx             context.Context
	mu              sync.Mutex
	bodyReader      *bodyReader
}

func NewReq(ctx context.Context, method, url, protocol, protocolVersion string, ahs map[string]string, connReader *bufio.Reader) *Req {
	req := &Req{
		Method:          method,
		Url:             url,
		Protocol:        protocol,
		ProtocolVersion: protocolVersion,
		ahs:             ahs,
		ctx:             ctx,
		bodyReader: &bodyReader{
			bufioReader: connReader,
		},
	}

	// here we assumed the content length header has been checked in connectionHandler
	clStr := req.GetHeader("Content-Length")
	var clInt int
	if clStr != "" {
		clInt, _ = strconv.Atoi(clStr)
	}

	req.bodyReader.contentLength = clInt
	return req
}

// we will assume the body length must be the same with content Content-Length header
func (r *Req) Body() io.Reader {
	return r.bodyReader
}

func (r *Req) GetHeader(key string) string {
	h, ok := r.ahs[strings.ToLower(key)]
	if ok {
		return h
	}
	return ""
}

func (r *Req) Ctx() context.Context {
	if r.ctx == nil {
		return context.Background()
	}
	return r.ctx
}

func (r *Req) IsBodyRead() bool {
	return r.bodyReader.IsBodyRead()
}

func (r *Req) CleanUpBodyBytes() error {
	if r.bodyReader.IsBodyRead() {
		return nil
	}

	_, err := io.ReadAll(r.Body())
	return err
}
