package http

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"lwc.com/servergo/internal/logger"
)

var (
	headerEnds = errors.New("Headers Ended")
	headerNoColon = errors.New("Header Line: there must at least be a colon in between the field key and the field value")
	headerWhiteSpaceBeforeColon = errors.New("Header Line: invalid structure, whitespace before colon is not allowed") 
)

func readHeader(ctx context.Context, hb []byte) (string, string, error) {
	l := logger.Get(ctx)

	hl := string(hb)
	l.Info(fmt.Sprintf("Header Line: %s", hl))
	hlTrim := strings.Trim(trimCRLF(hl), " ")
	if hlTrim == "" {
		return "", "", headerEnds
	}

	rawKey, rawValue, ok  := strings.Cut(hlTrim, ":")

	if !ok {
		return "", "", headerNoColon
	}

	key := strings.ToLower(rawKey)
	if strings.HasSuffix(key, " ") {
		return "", "", headerWhiteSpaceBeforeColon
	}

	value := strings.Trim(rawValue, " ")
	if key == "" || value == "" {
		return "", "", fmt.Errorf("Header Line: invalid structure, empty key: %s, or value: %s", rawKey, rawValue)
	}

	return key, value, nil
}
