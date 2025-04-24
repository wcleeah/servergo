package http

import (
	"context"
	"errors"
	"testing"
)

func TestHappy(t *testing.T) {
	startLine := "POST /123/123/123 HTTP/1.0"

	ctx := context.Background()
	sl, err := readStartLine(ctx, startLine)
	if err != nil {
		t.Fatalf("unexpected error %s", err.Error())
	}

	if sl.ProtocolVersion != "1.0" {
		t.Fatalf("version invalid, expected %s, got %s", "1.0", sl.ProtocolVersion)
	}

	if sl.Method != "POST" {
		t.Fatalf("method invalid, expected %s, got %s", "GET", sl.Method)
	}

	if sl.Url != "/123/123/123" {
		t.Fatalf("url invalid, expected %s, got %s", "/123/123/123", sl.Url)
	}
}

// Invalid structure 1
func TestInSt1(t *testing.T) {
	startLine := "GET /123/123/123 http/1.0 123123123123"

	ctx := context.Background()
	_, err := readStartLine(ctx, startLine)
	if err == nil {
		t.Fatalf("expected error %s", "invalid structure")
	}
}

// Invalid structure 2
func TestInSt2(t *testing.T) {
	startLine := "GET /123/123/123"

	ctx := context.Background()
	_, err := readStartLine(ctx, startLine)
	if err == nil {
		t.Fatalf("expected error %s", "invalid structure")
	}
}

func TestUnsupportedMethod(t *testing.T) {
	startLine := "ABC /123/123/123 HTTP/1.0"

	ctx := context.Background()
	_, err := readStartLine(ctx, startLine)
	if err == nil {
		t.Fatalf("expected error %s", "invalid structure")
	}
    if !errors.Is(err, unsupportedMethod) {
        t.Fatalf("wrong error, expected: %s, got: %s", unsupportedMethod.Error(), err)
    }
}

func TestUnsupportedProtocol(t *testing.T) {
	startLine := "GET /123/123/123 HAHAHAHAHAH/1.0"

	ctx := context.Background()
	_, err := readStartLine(ctx, startLine)
	if err == nil {
		t.Fatalf("expected error %s", "invalid structure")
	}
    if !errors.Is(err, unsupportedProtocol) {
        t.Fatalf("wrong error, expected: %s, got: %s", unsupportedProtocol.Error(), err)
    }
}

func TestUnsupportedProtocolVersion(t *testing.T) {
	startLine := "GET /123/123/123 HTTP/3.0"

	ctx := context.Background()
	_, err := readStartLine(ctx, startLine)
	if err == nil {
		t.Fatalf("expected error %s", "invalid structure")
	}
    if !errors.Is(err, unsupportedProtocolVersion) {
        t.Fatalf("wrong error, expected: %s, got: %s", unsupportedProtocolVersion.Error(), err)
    }
}
