package route

import (
	"context"
	"io"
	"strconv"
	"strings"
	"testing"
)

func (r *Req) TestTextBodyHappy(t *testing.T) {
	str := "Hello World"
	body := []byte(str)
	l := len(body)

	rc := io.NopCloser(strings.NewReader(str))
	defer rc.Close()
    ctx := context.Background()

	req := NewReq(ctx, "GET", "/", "HTTP/1.0", "1.0", map[string]string{"Content-Length": strconv.Itoa(l)}, rc)
	s, err := req.ReadTextBody()
	if err != nil {
		t.Fatalf("Error reading body: %v", err.Error())
	}

	if s != str {
		t.Fatalf("Body not equal")
	}
}
