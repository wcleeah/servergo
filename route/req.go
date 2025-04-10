package route

import (
	"bufio"
	"errors"
	"io"
	"net"
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
	conn            net.Conn
	body            []byte
	textBody        string
}

func NewReq(method, url, protocol, protocolVersion string, ahs map[string]string, conn net.Conn) *Req {
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
			return errors.New("Content length not specified in request")
		}

		clInt, err := strconv.Atoi(cl)
		if err != nil {
			return errors.New("Content length malformed")
		}

		bs := make([]byte, clInt)
		bufIoReader := bufio.NewReader(r.conn)
		_, err = io.ReadFull(bufIoReader, bs)
		if err != nil {
			return errors.New("Body malformed")
		}

		valid := utf8.Valid(bs)
		if !valid {
			return errors.New("Body malformed")
		}
		if len(bs) == 0 {
			return errors.New("Body malformed")
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
