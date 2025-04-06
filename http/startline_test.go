package http

import (
	"context"
	"testing"
)

func TestHappy(t *testing.T) {
    startLine := "POST /123/123/123 HTTP/1.1"

    ctx := context.Background()
    sl, err := readStartLine(ctx, startLine)
    if err != nil {
        t.Fatalf("unexpected error %s", err.Error())
    }

    if sl.HttpVersion != "1.1" {
        t.Fatalf("version invalid, expected %s, got %s", "1.1", sl.HttpVersion)
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
    startLine := "GET /123/123/123 http/1.1 123123123123"

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
