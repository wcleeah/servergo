package http 

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type tttt struct {}

func (t *tttt) Read(p []byte) (int, error) {
	println("haha")
	return 0, nil
}

func (t tttt) Close() error {
	println("haha")
	return nil
}


func TestReadBody(t *testing.T) {
	str := "Hello World"
	strByteLen := len(str)

	bufioReader := bufio.NewReader(strings.NewReader(str))
	bodyReader := NewBody(bufioReader, strByteLen)

	body, err := io.ReadAll(bodyReader)

	assert.NoError(t, err)
	assert.Equal(t, strByteLen, len(body))
	assert.Equal(t, str, string(body))
}

func TestReadBody_BodyInvalidUTF8(t *testing.T) {
	invalidBytes := []byte{0xC0, 0x80}

	bufioReader := bufio.NewReader(bytes.NewReader(invalidBytes))
	strByteLen := len(invalidBytes)

	bodyReader := NewBody(bufioReader, strByteLen)

	_, err := io.ReadAll(bodyReader)

	assert.Error(t, err)
	assert.Equal(t, BodyMalformed, err)
}

func TestReadBody_BodyEmpty(t *testing.T) {
	emptyBytes := []byte{}
	bufioReader := bufio.NewReader(bytes.NewReader(emptyBytes))
	strByteLen := 0

	bodyReader := NewBody(bufioReader, strByteLen)

	body, err := io.ReadAll(bodyReader)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(body))
}
