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
	Method          string
	Url             string
	ProtocolVersion string
	Protocol        string
}

func readStartLine(ctx context.Context, sls string) (*StartLine, error) {
	l := logger.Get(ctx)
	slsTrim := trimCRLF(sls)
	l.Info(fmt.Sprintf("Start Line: %s", slsTrim))
	slsSplitted := strings.Split(slsTrim, " ")

	if len(slsSplitted) != 3 {
		return nil, errors.New(fmt.Sprintf("Start Line: malformed structure, there are %d arguments", len(slsSplitted)))
	}
	method := slsSplitted[0]
	url := slsSplitted[1]
	httpVersion := slsSplitted[2]

    if !slices.Contains(SUPPORTED_METHOD, method) {
        return nil, unsupportedMethod
    }

	if !strings.HasPrefix(url, "/") {
		return nil, errors.New(fmt.Sprintf("Start Line: malformed url -> %s", url))
	}

	hvSplitted := strings.Split(httpVersion, "/")
	if len(hvSplitted) != 2 {
		return nil, errors.New(fmt.Sprintf("Start Line: malformed protocol and version -> %s", httpVersion))
	}

	protocol := hvSplitted[0]
	version := hvSplitted[1]
	if protocol != SUPPORTED_PROTOCOL {
		return nil, unsupportedProtocol
	}

	if !slices.Contains(SUPPORTED_PROTOCOL_VERSION, version) {
		return nil, unsupportedProtocolVersion
	}

	return &StartLine{
		Method:          method,
		Url:             url,
		ProtocolVersion: version,
		Protocol:        protocol,
	}, nil
}
