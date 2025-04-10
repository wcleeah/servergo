package route

import (
	"context"
	"errors"
	"fmt"
	"io"
	"testing"
)

func test(rs string, p *ResWriteParam) error {
	r, w := io.Pipe()
	defer r.Close()
	defer w.Close()
	res := NewRes("HTTP", "1.1", w)
	go func() {
		defer w.Close()
		res.Write(context.Background(), p)
	}()
	// if the result are correct, the size of the correct string must be the same as the result
	bs := make([]byte, len(p.Body))
	_, err := r.Read(bs)
	if err != nil {
		return err
	}
	str := string(bs)
	for i, r := range str {
		if string(r) != string(rs[i]) {
			return errors.New(fmt.Sprintf("Mismatch charater: expected %s, got %s", string(rs[i]), string(r)))
		}
	}
	return nil
}

func TestWriteResHappy_Header(t *testing.T) {
	bodyStr := "Hello World"
	body := []byte(bodyStr)
	statusCode := "200"
	statusStr := codeMsgMap[statusCode]
    resultStr := fmt.Sprintf("HTTP/1.1 %s %s\r\ntest: haha\r\nContent-Length: %d\r\n\r\n%s", statusCode, statusStr, len(body), bodyStr)
	rwp := ResWriteParam{
		StatusCode: statusCode,
		Body:       body,
		Ahs:        map[string]string{
            "test": "haha",
        },
	}
	err := test(resultStr, &rwp)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWriteResHappy_Body(t *testing.T) {
	bodyStr := "Hello World"
	body := []byte(bodyStr)
	statusCode := "200"
	statusStr := codeMsgMap[statusCode]
	resultStr := fmt.Sprintf("HTTP/1.1 %s %s\r\nContent-Length: %d\r\n\r\n%s", statusCode, statusStr, len(body), bodyStr)
	rwp := ResWriteParam{
		StatusCode: statusCode,
		Body:       body,
		Ahs:        map[string]string{},
	}
	err := test(resultStr, &rwp)
	if err != nil {
		t.Fatal(err)
	}
}
