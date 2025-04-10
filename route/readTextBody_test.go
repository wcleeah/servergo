package route

import (
	"io"
	"strconv"
	"strings"
	"testing"
)

func (r *Req) TestHappy(t *testing.T) {
    str := "Hello World"
    body := []byte(str)
    l := len(body)

    rc := io.NopCloser(strings.NewReader(str))

    req := NewReq("GET", "/", "HTTP/1.1", "1.1", map[string]string{"Content-Length": strconv.Itoa(l)}, rc)
    s, err := req.ReadTextBody()
    if err != nil {
        t.Fatalf("Error reading body: %v", err.Error())
    }

    if s != str {
        t.Fatalf("Body not equal")
    }
}
