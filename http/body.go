package http

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

func readBody(bs []byte) (string, error) {
    valid := utf8.Valid(bs)
    if !valid {
        return "", errors.New("Invalid byte slice")
    }
    str := string(bs)
    fmt.Printf("str: %s\n", str)
    fmt.Printf("len str: %d\n", len(str))
    if str == "" {
        return "", errors.New("Empty byte slice")
    }

    return str, nil
}
