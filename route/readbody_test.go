package route

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"
	"testing"
)

func TestReadBodyHappy(t *testing.T) {
	str := "Hello World"
	body := []byte(str)
	l := len(body)

	rc := io.NopCloser(strings.NewReader(str))

	req := NewReq("GET", "/", "HTTP/1.1", "1.1", map[string]string{"Content-Length": strconv.Itoa(l)}, rc)
	b, err := req.ReadBody()
	if err != nil {
		t.Fatalf("Error reading body: %v", err.Error())
	}
	if l != len(b) {
		t.Fatalf("Body length not equal")
	}
	if string(b) != str {
		t.Fatalf("Body not equal")
	}
}

func TestReadBodySad_NoContentLength(t *testing.T) {
	str := "Hello World"
	rc := io.NopCloser(strings.NewReader(str))

	req := NewReq("GET", "/", "HTTP/1.1", "1.1", map[string]string{}, rc)
	_, err := req.ReadBody()
	if err == nil {
		t.Fatal("ReadBody should have failed")
	}

	if !errors.Is(err, ContentLengthNotSpecified) {
		t.Fatalf("Error should be %s, but it is %s now", ContentLengthNotSpecified.Error(), err.Error())
	}
}

func TestReadBodySad_ContentLengthMalformed(t *testing.T) {
	str := "Hello World"
	rc := io.NopCloser(strings.NewReader(str))

	req := NewReq("GET", "/", "HTTP/1.1", "1.1", map[string]string{"Content-Length": "abc"}, rc)
	_, err := req.ReadBody()
	if err == nil {
		t.Fatal("ReadBody should have failed")
	}

	if !errors.Is(err, ContentLengthMalformed) {
		t.Fatalf("Error should be %s, but it is %s now", ContentLengthMalformed.Error(), err.Error())
	}
}

func TestReadBodySad_BodyInvalidUTF8(t *testing.T) {
	invalidBytes := []byte{0xC0, 0x80}
	rc := io.NopCloser(bytes.NewReader(invalidBytes))
	l := len(invalidBytes)

	req := NewReq("GET", "/", "HTTP/1.1", "1.1", map[string]string{"Content-Length": strconv.Itoa(l)}, rc)
	_, err := req.ReadBody()
	if err == nil {
		t.Fatal("ReadBody should have failed")
	}

	if !errors.Is(err, BodyMalformed) {
		t.Fatalf("Error should be %s, but it is %s now", BodyMalformed.Error(), err.Error())
	}
}

func TestReadBodySad_BodyEmpty(t *testing.T) {
	invalidBytes := []byte{}
	rc := io.NopCloser(bytes.NewReader(invalidBytes))
	l := len(invalidBytes)

	req := NewReq("GET", "/", "HTTP/1.1", "1.1", map[string]string{"Content-Length": strconv.Itoa(l)}, rc)
	_, err := req.ReadBody()
	if err == nil {
		t.Fatal("ReadBody should have failed")
	}

	if !errors.Is(err, BodyMalformed) {
		t.Fatalf("Error should be %s, but it is %s now", BodyMalformed.Error(), err.Error())
	}
}
