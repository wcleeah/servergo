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
	sl, err := readStartLine(ctx, []byte(startLine))
	
	assert.Nil(t, err)
	assert.Equal(t, "1.0", sl.ProtocolVersion)
	assert.Equal(t, "POST", sl.Method)
	assert.Equal(t, "/123/123/123", sl.Url)
}

func TestSupportedProtocolVersion(t *testing.T) {
	startLine := "GET /123/123/123 HTTP/1.0"

	ctx := context.Background()
	sl, err := readStartLine(ctx, []byte(startLine))

	assert.Nil(t, err)
	assert.Equal(t, "1.0", sl.ProtocolVersion)
	assert.Equal(t, "GET", sl.Method)
	assert.Equal(t, "/123/123/123", sl.Url)

	startLine = "GET /123/123/123 HTTP/1.1"

	sl, err = readStartLine(ctx, []byte(startLine))

	assert.Nil(t, err)
	assert.Equal(t, "1.1", sl.ProtocolVersion)
	assert.Equal(t, "GET", sl.Method)
	assert.Equal(t, "/123/123/123", sl.Url)
}

// invalid structure, too many parts 
func TestTooManyParts(t *testing.T) {
	startLine := "GET /123/123/123 http/1.0 123123123123"

	ctx := context.Background()
	sl, err := readStartLine(ctx, []byte(startLine))

	assert.Nil(t, sl)
	assert.Error(t, err)
}

// invalid structure, missing parts
func TestMissingParts(t *testing.T) {
	startLine := "GET /123/123/123"

	ctx := context.Background()
	sl, err := readStartLine(ctx, []byte(startLine))

	assert.Nil(t, sl)
	assert.Error(t, err)
}

// Path does not starts with /
func TestMissingSlashPrefixForPath(t *testing.T) {
	startLine := "GET 123/123/123 HTTP/1.1"

	ctx := context.Background()
	sl, err := readStartLine(ctx, []byte(startLine))

	assert.Nil(t, sl)
	assert.Error(t, err)
}

// unsupported method
func TestUnsupportedMethod(t *testing.T) {
	startLine := "ABC /123/123/123 HTTP/1.0"

	ctx := context.Background()
	sl, err := readStartLine(ctx, []byte(startLine))

	assert.Nil(t, sl)
	assert.Error(t, err)
	assert.Equal(t, unsupportedMethod, err)
}

// unsupported protocol
func TestUnsupportedProtocol(t *testing.T) {
	startLine := "GET /123/123/123 HAHAHAHAHAH/1.0"

	ctx := context.Background()
	sl, err := readStartLine(ctx, []byte(startLine))

	assert.Nil(t, sl)
	assert.Error(t, err)
	assert.Equal(t, unsupportedProtocol, err)
}

func TestUnsupportedProtocolVersion(t *testing.T) {
	startLine := "GET /123/123/123 HTTP/3.0"

	ctx := context.Background()
	sl, err := readStartLine(ctx, []byte(startLine))

	assert.Nil(t, sl)
	assert.Error(t, err)
	assert.Equal(t, unsupportedProtocolVersion, err)
}
