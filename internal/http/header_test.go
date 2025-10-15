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
	key, value, noMoreHeader, err := ReadHeader(ctx, []byte(hl))

	assert.Equal(t, "content-length", key)
	assert.Equal(t, "12345", value)
	assert.False(t, noMoreHeader)
	assert.Nil(t, err)

	hl = "custoM: abcde\r\n"

	key, value, noMoreHeader, err = ReadHeader(ctx, []byte(hl))

	assert.Equal(t, "custom", key)
	assert.Equal(t, "abcde", value)
	assert.False(t, noMoreHeader)
	assert.Nil(t, err)
}

// OWC
func TestOWC(t *testing.T) {
	// many OWC
	hl := "Content-Length:                             12345                              \r\n"

	ctx := context.Background()
	key, value, noMoreHeader, err := ReadHeader(ctx, []byte(hl))

	assert.Equal(t, "content-length", key)
	assert.Equal(t, "12345", value)
	assert.False(t, noMoreHeader)
	assert.Nil(t, err)

	// no OWC
	hl = "custom:abcde\r\n"

	key, value, noMoreHeader, err = ReadHeader(ctx, []byte(hl))

	assert.Equal(t, "custom", key)
	assert.Equal(t, "abcde", value)
	assert.False(t, noMoreHeader)
	assert.Nil(t, err)
}

// EOF
func TestHeaderEOF(t *testing.T) {
	hl := "\r\n"

	ctx := context.Background()
	key, value, noMoreHeader, err := ReadHeader(ctx, []byte(hl))

	assert.Empty(t, key)
	assert.Empty(t, value)
	assert.True(t, noMoreHeader)
	assert.Nil(t, err)
}

// Header value contains colon
func TestColonValue(t *testing.T) {
	hl := "Host: localhost:3000"

	ctx := context.Background()
	key, value, noMoreHeader, err := ReadHeader(ctx, []byte(hl))

	assert.Equal(t, "host", key)
	assert.Equal(t, "localhost:3000", value)
	assert.False(t, noMoreHeader)
	assert.Nil(t, err)
}

// Space before colon
func TestSpaceBeforeColon(t *testing.T) {
	hl := "Content-Length : 123\r\n"

	ctx := context.Background()
	key, value, noMoreHeader, err := ReadHeader(ctx, []byte(hl))

	assert.Empty(t, key)
	assert.Empty(t, value)
	assert.True(t, noMoreHeader)
	assert.Error(t, err)
	assert.Equal(t, HeaderWhiteSpaceBeforeColon, err)
}

// incorrect structure, no colon
func TestNoColon(t *testing.T) {
	hl := "Content-Length 12345"

	ctx := context.Background()
	key, value, noMoreHeader, err := ReadHeader(ctx, []byte(hl))

	assert.Empty(t, key)
	assert.Empty(t, value)
	assert.True(t, noMoreHeader)
	assert.Error(t, err)
	assert.Equal(t, HeaderNoColon, err)
}

// incorrect structure, only colon
func TestOnlyColon(t *testing.T) {
	hl := ":"

	ctx := context.Background()
	key, value, noMoreHeader, err := ReadHeader(ctx, []byte(hl))

	assert.Empty(t, key)
	assert.Empty(t, value)
	assert.True(t, noMoreHeader)
	assert.Error(t, err)
	hl = "abc:"

	key, value, noMoreHeader, err = ReadHeader(ctx, []byte(hl))

	assert.Empty(t, key)
	assert.Empty(t, value)
	assert.True(t, noMoreHeader)
	assert.Error(t, err)

	hl = ":abc"

	key, value, noMoreHeader, err = ReadHeader(ctx, []byte(hl))

	assert.Empty(t, key)
	assert.Empty(t, value)
	assert.True(t, noMoreHeader)
	assert.Error(t, err)
}
