package http

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"lwc.com/servergo/logger"
)

type StartLine struct {
	Method      string
	Url         string
    MPlusUrl     string
	HttpVersion string
}

func readStartLine(ctx context.Context, sls string) (*StartLine, error) {
	l := logger.Get(ctx)
    slsTrim := trimCRLF(sls)
	l.Info(fmt.Sprintf("Start Line: %s", slsTrim))
	slsSplitted := strings.Split(slsTrim, " ")

	if len(slsSplitted) != 3 {
		return nil, errors.New(fmt.Sprintf("Start Line: invalid structure, there are %d arguments", len(slsSplitted)))
	}
	method := slsSplitted[0]
	url := slsSplitted[1]
	httpVersion := slsSplitted[2]

	if !strings.HasPrefix(url, "/") {
		return nil, errors.New(fmt.Sprintf("Start Line: invalid url -> %s", url))
	}

	hvSplitted := strings.Split(httpVersion, "/")
	if len(hvSplitted) != 2 {
		return nil, errors.New(fmt.Sprintf("Start Line: invalid protocol and version -> %s", httpVersion))
	}

	protocol := hvSplitted[0]
	version := hvSplitted[1]
	if protocol != SUPPORTED_PROTOCOL {
		return nil, errors.New(fmt.Sprintf("Start Line: invalid protocol -> %s", protocol))
	}

    if !slices.Contains(SUPPORTED_PROTOCOL_VERSION, version) {
		return nil, errors.New(fmt.Sprintf("Start Line: invalid protocol version -> %s", version))
    }

	return &StartLine{
        Method: method,
        Url: url,
        MPlusUrl: method + url,
        HttpVersion: version,
    }, nil
}
