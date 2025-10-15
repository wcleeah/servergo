package http

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// correct value
func TestHeaderCorrectValue(t *testing.T) {
	hl := "Content-Length: 12345\r\n"

	ctx := context.Background()
	key, value, err := readHeader(ctx, []byte(hl))

	assert.Nil(t, err)
	assert.Equal(t, "content-length", key)
	assert.Equal(t, "12345", value)

	hl = "custoM: abcde\r\n"

	key, value, err = readHeader(ctx, []byte(hl))

	assert.Nil(t, err)
	assert.Equal(t, "custom", key)
	assert.Equal(t, "abcde", value)
}

// OWC
func TestOWC(t *testing.T) {
	// many OWC
	hl := "Content-Length:                             12345                              \r\n"

	ctx := context.Background()
	key, value, err := readHeader(ctx, []byte(hl))

	assert.Nil(t, err)
	assert.Equal(t, "content-length", key)
	assert.Equal(t, "12345", value)

	// no OWC
	hl = "custom:abcde\r\n"

	key, value, err = readHeader(ctx, []byte(hl))

	assert.Nil(t, err)
	assert.Equal(t, "custom", key)
	assert.Equal(t, "abcde", value)
}

// EOF
func TestHeaderEOF(t *testing.T) {
	hl := "\r\n"

	ctx := context.Background()
	key, value, err := readHeader(ctx, []byte(hl))

	assert.Empty(t, key)
	assert.Empty(t, value)
	assert.Error(t, err)
	assert.Equal(t, headerEnds, err)
}

// Space before colon
func TestSpaceBeforeColon(t *testing.T) {
	hl := "Content-Length : 123\r\n"

	ctx := context.Background()
	key, value, err := readHeader(ctx, []byte(hl))

	assert.Empty(t, key)
	assert.Empty(t, value)
	assert.Error(t, err)
	assert.Equal(t, headerWhiteSpaceBeforeColon, err)
}

// Header value contains colon
func TestColonValue(t *testing.T) {
	hl := "Host: localhost:3000"

	ctx := context.Background()
	key, value, err := readHeader(ctx, []byte(hl))

	assert.Nil(t, err)
	assert.Equal(t, "host", key)
	assert.Equal(t, "localhost:3000", value)
}

// incorrect structure, no colon
func TestNoColon(t *testing.T) {
	hl := "Content-Length 12345"

	ctx := context.Background()
	key, value, err := readHeader(ctx, []byte(hl))

	assert.Empty(t, key)
	assert.Empty(t, value)
	assert.Error(t, err)
	assert.Equal(t, headerNoColon, err)
}

// incorrect structure, only colon
func TestOnlyColon(t *testing.T) {
	hl := ":"

	ctx := context.Background()
	key, value, err := readHeader(ctx, []byte(hl))

	assert.Empty(t, key)
	assert.Empty(t, value)
	assert.Error(t, err)
	hl = "abc:"

	key, value, err = readHeader(ctx, []byte(hl))

	assert.Empty(t, key)
	assert.Empty(t, value)
	assert.Error(t, err)

	hl = ":abc"

	key, value, err = readHeader(ctx, []byte(hl))

	assert.Empty(t, key)
	assert.Empty(t, value)
	assert.Error(t, err)
}
