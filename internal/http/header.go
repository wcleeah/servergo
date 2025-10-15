package http

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"lwc.com/servergo/internal/logger"
)

var (
	HeaderNoColon = errors.New("Header Line: there must at least be a colon in between the field key and the field value")
	HeaderWhiteSpaceBeforeColon = errors.New("Header Line: invalid structure, whitespace before colon is not allowed") 
)

func ReadHeader(ctx context.Context, hb []byte) (string, string, bool, error) {
	l := logger.Get(ctx)

	hl := string(hb)
	l.Info(fmt.Sprintf("Header Line: %s", hl))
	hlTrim := strings.Trim(trimCRLF(hl), " ")
	if hlTrim == "" {
		return "", "", true, nil 
	}

	rawKey, rawValue, ok  := strings.Cut(hlTrim, ":")

	if !ok {
		return "", "", true, HeaderNoColon
	}

	key := strings.ToLower(rawKey)
	if strings.HasSuffix(key, " ") {
		return "", "", true, HeaderWhiteSpaceBeforeColon
	}

	value := strings.Trim(rawValue, " ")
	if key == "" || value == "" {
		return "", "", true, fmt.Errorf("Header Line: invalid structure, empty key: %s, or value: %s", rawKey, rawValue)
	}

	return key, value, false, nil
}
