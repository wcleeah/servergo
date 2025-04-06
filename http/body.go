package http

import (
	"errors"
	"unicode/utf8"
)

type Body struct {
    Raw []byte
    Text string
}

func readBody(bs []byte) (*Body, error) {
    valid := utf8.Valid(bs)
    if !valid {
        return nil, errors.New("Invalid byte slice")
    }
    if len(bs) == 0 {
        return nil, errors.New("Empty byte slice")
    }

    return &Body{
        Raw: bs,
        Text: string(bs),
    }, nil
}
