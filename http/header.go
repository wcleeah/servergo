package http

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"lwc.com/servergo/logger"
)

var headerEOF = errors.New("Headers ended")

func readHeader(ctx context.Context, hl string) (string, string, error) {
	l := logger.Get(ctx)
	hlTrim := trimCRLF(hl)
	if hlTrim == "" {
		return "", "", headerEOF
	}

	l.Info("Header Line", "hl", hlTrim)

	hSplitted := strings.Split(hlTrim, ": ")

	if len(hSplitted) != 2 {
		return "", "", errors.New(fmt.Sprintf("Header Line: invalid structure, ther are %d arguments", len(hSplitted)))
	}
	key := strings.Trim(hSplitted[0], " ")
	value := strings.Trim(hSplitted[1], " ")

	return key, value, nil
}
