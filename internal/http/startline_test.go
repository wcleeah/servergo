package http

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// correct value
func TestStartLineCorrectValue(t *testing.T) {
	startLine := "POST /123/123/123 HTTP/1.0"

	ctx := context.Background()
	method, url, protocolVersion, protocol, err := ReadStartLine(ctx, []byte(startLine))
	
	assert.Nil(t, err)
	assert.Equal(t, "HTTP", protocol)
	assert.Equal(t, "1.0", protocolVersion)
	assert.Equal(t, "POST", method)
	assert.Equal(t, "/123/123/123", url)
}

func TestSupportedProtocolVersion(t *testing.T) {
	startLine := "GET /123/123/123 HTTP/1.0"

	ctx := context.Background()
	method, url, protocolVersion, protocol, err := ReadStartLine(ctx, []byte(startLine))

	assert.Nil(t, err)
	assert.Equal(t, "HTTP", protocol)
	assert.Equal(t, "1.0", protocolVersion)
	assert.Equal(t, "GET", method)
	assert.Equal(t, "/123/123/123", url)

	startLine = "GET /123/123/123 HTTP/1.1"

	method, url, protocolVersion, protocol, err = ReadStartLine(ctx, []byte(startLine))

	assert.Nil(t, err)
	assert.Equal(t, "HTTP", protocol)
	assert.Equal(t, "1.1", protocolVersion)
	assert.Equal(t, "GET", method)
	assert.Equal(t, "/123/123/123", url)
}

// invalid structure, too many parts 
func TestTooManyParts(t *testing.T) {
	startLine := "GET /123/123/123 http/1.0 123123123123"

	ctx := context.Background()
	method, url, protocolVersion, protocol, err := ReadStartLine(ctx, []byte(startLine))

	assert.Empty(t, method)
	assert.Empty(t, url)
	assert.Empty(t, protocolVersion)
	assert.Empty(t, protocol)
	assert.Error(t, err)
}

// invalid structure, missing parts
func TestMissingParts(t *testing.T) {
	startLine := "GET /123/123/123"

	ctx := context.Background()
	method, url, protocolVersion, protocol, err := ReadStartLine(ctx, []byte(startLine))

	assert.Empty(t, method)
	assert.Empty(t, url)
	assert.Empty(t, protocolVersion)
	assert.Empty(t, protocol)
	assert.Error(t, err)
}

// Path does not starts with /
func TestMissingSlashPrefixForPath(t *testing.T) {
	startLine := "GET 123/123/123 HTTP/1.1"

	ctx := context.Background()
	method, url, protocolVersion, protocol, err := ReadStartLine(ctx, []byte(startLine))

	assert.Empty(t, method)
	assert.Empty(t, url)
	assert.Empty(t, protocolVersion)
	assert.Empty(t, protocol)
	assert.Error(t, err)
}

// unsupported method
func TestUnsupportedMethod(t *testing.T) {
	startLine := "ABC /123/123/123 HTTP/1.0"

	ctx := context.Background()
	method, url, protocolVersion, protocol, err := ReadStartLine(ctx, []byte(startLine))

	assert.Empty(t, method)
	assert.Empty(t, url)
	assert.Empty(t, protocolVersion)
	assert.Empty(t, protocol)
	assert.Error(t, err)
	assert.Equal(t, UnsupportedMethod, err)
}

// unsupported protocol
func TestUnsupportedProtocol(t *testing.T) {
	startLine := "GET /123/123/123 HAHAHAHAHAH/1.0"

	ctx := context.Background()
	method, url, protocolVersion, protocol, err := ReadStartLine(ctx, []byte(startLine))

	assert.Empty(t, method)
	assert.Empty(t, url)
	assert.Empty(t, protocolVersion)
	assert.Empty(t, protocol)
	assert.Error(t, err)
	assert.Equal(t, UnsupportedProtocol, err)
}

func TestUnsupportedProtocolVersion(t *testing.T) {
	startLine := "GET /123/123/123 HTTP/3.0"

	ctx := context.Background()
	method, url, protocolVersion, protocol, err := ReadStartLine(ctx, []byte(startLine))

	assert.Empty(t, method)
	assert.Empty(t, url)
	assert.Empty(t, protocolVersion)
	assert.Empty(t, protocol)
	assert.Error(t, err)
	assert.Equal(t, UnsupportedProtocolVersion, err)
}
