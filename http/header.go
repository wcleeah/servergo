package http

import (
	"context"
	"errors"
	"fmt"
	"strings"
)


func readHeader(ctx context.Context, hl string) (string, string, error) {
	hlTrim := trimCRLF(hl)
	if hlTrim == "" {
		return "", "", headerEnds
	}

	hSplitted := strings.Split(hlTrim, ": ")

	if len(hSplitted) != 2 {
		return "", "", errors.New(fmt.Sprintf("Header Line: invalid structure, ther are %d arguments", len(hSplitted)))
	}
	key := strings.Trim(hSplitted[0], " ")
	value := strings.Trim(hSplitted[1], " ")

	return key, value, nil
}
