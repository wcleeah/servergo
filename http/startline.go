package http

import (
	"context"
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

func readStartLine(ctx context.Context, slb []byte) (*StartLine, error) {
	l := logger.Get(ctx)
	sls := trimCRLF(string(slb))
	l.Info(fmt.Sprintf("Start Line: %s", sls))
	// go internal uses Cut, because of standardization?
	// also, not like field lines, start lines only allow one space in between
	slsSplitted := strings.Split(sls, " ")

	if len(slsSplitted) != 3 {
		return nil, fmt.Errorf("Start Line: malformed structure, there are %d arguments", len(slsSplitted))
	}
	method := slsSplitted[0]
	url := slsSplitted[1]
	httpVersion := slsSplitted[2]

    if !slices.Contains(SUPPORTED_METHOD, method) {
        return nil, unsupportedMethod
    }

	if !strings.HasPrefix(url, "/") {
		return nil, fmt.Errorf("Start Line: malformed url -> %s", url)
	}

	hvSplitted := strings.Split(httpVersion, "/")
	if len(hvSplitted) != 2 {
		return nil, fmt.Errorf("Start Line: malformed protocol and version -> %s", httpVersion)
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
