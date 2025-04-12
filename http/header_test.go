package http

import (
	"context"
	"errors"
	"testing"
)

// correct value
func TestHappyCV(t *testing.T) {
	hl := "Content-Length: 12345\r\n"

	ctx := context.Background()
	key, value, err := readHeader(ctx, hl)
	if err != nil {
		t.Fatalf("unexpected error %s", err.Error())
	}

	if key != "Content-Length" {
		t.Fatalf("key invalid, expected %s, got %s", "Content-Length", key)
	}

	if value != "12345" {
		t.Fatalf("value invalid, expected %s, got %s", "12345", value)
	}
}

// EOF
func TestHappyEOF(t *testing.T) {
	hl := "\r\n"

	ctx := context.Background()
	key, value, err := readHeader(ctx, hl)
	if err != nil {
		if !errors.Is(headerEnds, err) {
			t.Fatalf("unexpected error %s", err.Error())
		}
	}

	if key != "" {
		t.Fatalf("key invalid, expected %s, got %s", "Content-Length", key)
	}

	if value != "" {
		t.Fatalf("value invalid, expected %s, got %s", "12345", value)
	}
}

// Host header with port
func TestHappyHost(t *testing.T) {
	hl := "Host: localhost:3000"

	ctx := context.Background()
	key, value, err := readHeader(ctx, hl)
	if err != nil {
		t.Fatalf("unexpected error %s", err.Error())
	}

	if key != "Host" {
		t.Fatalf("key invalid, expected %s, got %s", "Content-Length", key)
	}

	if value != "localhost:3000" {
		t.Fatalf("value invalid, expected %s, got %s", "12345", value)
	}
}

// incorrect structure, no colon
func TestSadNoColon(t *testing.T) {
	hl := "Content-Length 12345"

	ctx := context.Background()
	_, _, err := readHeader(ctx, hl)
	if err == nil {
		t.Fatalf("it should throw error")
	}
}

// incorrect structure, too many colon
func TestSadTooManyColon(t *testing.T) {
	hl := "Content-Length: 12345: "

	ctx := context.Background()
	_, _, err := readHeader(ctx, hl)
	if err == nil {
		t.Fatalf("it should throw error")
	}
}
