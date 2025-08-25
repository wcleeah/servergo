package route

import (
	"bytes"
	"context"
	"io"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
)

type Req struct {
	Method          string
	Url             string
	Protocol        string
	ProtocolVersion string
	ahs             map[string]string
	ctx             context.Context
	mu              sync.Mutex
	isBodyRead      bool
	body            []byte
	connReader      io.Reader
	// using io.ReadCloser provides a few benefit
	// 1. testing will be easier because we don't need to create a fake connection
	// 2. req probably should not be writing to the connection (?)
	conn io.ReadCloser
}

func NewReq(ctx context.Context, method, url, protocol, protocolVersion string, ahs map[string]string, conn io.ReadCloser, connReader io.Reader) *Req {
	return &Req{
		Method:          method,
		Url:             url,
		Protocol:        protocol,
		ProtocolVersion: protocolVersion,
		ahs:             ahs,
		conn:            conn,
		body:            []byte{},
		ctx:             ctx,
		connReader:      connReader,
	}
}

// TODO: add mutex
// Kind of understood instead of writting this way, go package uses a Reader instead
// - We do not have to return the []byte, we can just return number
// - It suites the convention of passing in a byte array, then writing init
// I am too lazy to change this part, keeping as it is for now
// we will assume the body length must be the same with content Content-Length header
func (r *Req) Body() ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isBodyRead {
		return r.body, nil
	}

	cl := r.GetHeader("Content-Length")
	if cl == "" {
		return nil, nil
	}

	clInt, err := strconv.Atoi(cl)
	if err != nil {
		return nil, ContentLengthMalformed
	}
	if clInt == 0 {
		return r.body, nil
	}

	bs := make([]byte, clInt)

	lr := io.LimitReader(r.connReader, int64(clInt))

	n, err := lr.Read(bs)
	if err != nil {
		return nil, err
	}

	if n != clInt {
		return nil, BodyMalformed
	}

	valid := utf8.Valid(bs)
	if !valid {
		return nil, BodyMalformed
	}
	if len(bs) == 0 {
		return nil, BodyMalformed
	}

	r.body = bytes.TrimSuffix(bs, crlf)
	r.isBodyRead = true

	return r.body, nil
}

func (r *Req) GetHeader(key string) string {
	h, ok := r.ahs[strings.ToLower(key)]
	if ok {
		return h
	}
	return ""
}

func (r *Req) Ctx() context.Context {
	return r.ctx
}

func (r *Req) IsBodyRead() bool {
	return r.isBodyRead
}

func (r *Req) CleanUpBodyBytes() error {
	_, err := r.Body()
	return err
}
