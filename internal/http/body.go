package http

import (
	"bufio"
	"errors"
	"io"
	"sync"
	"unicode/utf8"
)

var BodyMalformed = errors.New("Body malformed")

type Body struct {
	mu            sync.Mutex
	bufioReader   *bufio.Reader
	contentLength int
	readN         int
}

func NewBody(bufioReader *bufio.Reader, contentLength int) *Body {
	return &Body{
		bufioReader: bufioReader,
		contentLength: contentLength,
	}
}

func (r *Body) Read(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.contentLength == 0 || r.contentLength-r.readN <= 0 {
		return 0, io.EOF
	}

	targetN := min(r.contentLength-r.readN, len(p))

	lr := io.LimitReader(r.bufioReader, int64(targetN))
	totalN := 0
	bs := make([]byte, 0)

	for {
		temp := make([]byte, r.contentLength-totalN)
		n, err := lr.Read(temp)
		totalN += n
		if err != nil {
			return totalN, err
		}

		bs = append(bs, temp...)

		if totalN == targetN {
			break
		}

		if totalN < targetN {
			continue
		}
	}

	valid := utf8.Valid(bs)
	if !valid {
		return totalN, BodyMalformed
	}
	if len(bs) == 0 {
		return totalN, BodyMalformed
	}

	r.readN += totalN
	copy(p, bs)

	return totalN, nil
}

func (r *Body) IsBodyRead() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.readN >= r.contentLength
}
