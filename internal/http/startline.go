package http

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"lwc.com/servergo/internal/logger"
)

var supportedMethods = []string{
	"GET",
	"POST",
	"PUT",
	"DELETE",
	"OPTIONS",
}

var supportedProtocolVersions = []string{
	"1.0",
	"1.1",
}

const (
	SUPPORTED_PROTOCOL       = "HTTP"
	DEFAULT_PROTOCOL_VERSION = "1.1"
)

var UnsupportedMethod = errors.New("Unsupported Method")
var UnsupportedProtocol = errors.New("Unsupported Protocol")
var UnsupportedProtocolVersion = errors.New("Unsupported Protocol Version")

func ReadStartLine(ctx context.Context, slb []byte) (string, string, string, string, error) {
	l := logger.Get(ctx)
	sls := trimCRLF(string(slb))
	l.Info(fmt.Sprintf("Start Line: %s", sls))
	// go internal uses Cut, because of standardization?
	// also, not like field lines, start lines only allow one space in between
	slsSplitted := strings.Split(sls, " ")

	if len(slsSplitted) != 3 {
		return "", "", "", "", fmt.Errorf("Start Line: malformed structure, there are %d arguments", len(slsSplitted))
	}
	method := slsSplitted[0]
	url := slsSplitted[1]
	httpVersion := slsSplitted[2]

	if !slices.Contains(supportedMethods, method) {
		return "", "", "", "", UnsupportedMethod
	}

	if !strings.HasPrefix(url, "/") {
		return "", "", "", "", fmt.Errorf("Start Line: malformed url -> %s", url)
	}

	hvSplitted := strings.Split(httpVersion, "/")
	if len(hvSplitted) != 2 {
		return "", "", "", "", fmt.Errorf("Start Line: malformed protocol and version -> %s", httpVersion)
	}

	protocol := hvSplitted[0]
	version := hvSplitted[1]
	if protocol != SUPPORTED_PROTOCOL {
		return "", "", "", "", UnsupportedProtocol
	}

	if !slices.Contains(supportedProtocolVersions, version) {
		return "", "", "", "", UnsupportedProtocolVersion
	}

	return method, url, version, protocol, nil
}
