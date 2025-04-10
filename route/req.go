package route

import (
	"bufio"
	"io"
	"strconv"
	"sync"
	"unicode/utf8"
)

type Req struct {
	Method          string
	Url             string
	Protocol        string
	ProtocolVersion string
	ahs             map[string]string
    // using io.ReadCloser provides a few benefit
    // 1. testing will be easier because we don't need to create a fake connection
    // 2. req probably should not be writing to the connection (?)
	conn            io.ReadCloser
	body            []byte
	textBody        string
}

func NewReq(method, url, protocol, protocolVersion string, ahs map[string]string, conn io.ReadCloser) *Req {
	return &Req{
		Method:          method,
		Url:             url,
		Protocol:        protocol,
		ProtocolVersion: protocolVersion,
		ahs:             ahs,
		conn:            conn,
		body:            nil,
	}
}

func (r *Req) ReadBody() ([]byte, error) {
	ov := sync.OnceValue(func() error {
		cl := r.GetHeader("Content-Length")
		if cl == "" {
			return ContentLengthNotSpecified
		}

		clInt, err := strconv.Atoi(cl)
		if err != nil {
			return ContentLengthMalformed
		}

		bs := make([]byte, clInt)
		bufIoReader := bufio.NewReader(r.conn)
		_, err = io.ReadFull(bufIoReader, bs)
		if err != nil {
			return BodyMalformed
		}

		valid := utf8.Valid(bs)
		if !valid {
			return BodyMalformed
		}
		if len(bs) == 0 {
			return BodyMalformed
		}

		r.body = bs
		return nil
	})
	err := ov()
	if err != nil {
		return nil, err
	}

	return r.body, nil
}

func (r *Req) ReadTextBody() (string, error) {
	ov := sync.OnceValue(func() error {
		bs, err := r.ReadBody()
		if err != nil {
			return err
		}
		r.textBody = string(bs)
		return nil
	})
	err := ov()
	if err != nil {
		return "", err
	}
	return r.textBody, nil
}

func (r *Req) GetHeader(key string) string {
	h, ok := r.ahs[key]
	if ok {
		return h
	}
	return ""
}
