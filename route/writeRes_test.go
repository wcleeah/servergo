package route

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"testing"
)

func TestWriteResHappy(t *testing.T) {
	r, w := io.Pipe()
	defer r.Close()
	defer w.Close()

	res := NewRes("HTTP", "1.1", w)

	bodyStr := "Hello World"
	body := []byte(bodyStr)
	l := len(body)
	statusCode := "200"
	statusStr := codeMsgMap[statusCode]
	resultStr := fmt.Sprintf("HTTP/1.1 %s %s\r\nContent-Length: %d\r\n\r\n%s", statusCode, statusStr, len(body), bodyStr)
	log.Println(resultStr)
	rwp := ResWriteParam{
		StatusCode: statusCode,
		Body:       body,
		Ahs: map[string]string{
			"Content-Length": strconv.Itoa(l),
		},
	}

	c := make(chan string)
	go func() {
		bs := make([]byte, 0)
		_, err := r.Read(bs)
		if err != nil {
			c <- err.Error()
			return
		}
		str := string(bs)
		if str != resultStr {
			log.Println("expected")
			log.Println(resultStr)
			log.Println("actual")
			log.Println(str)
			c <- "result invalid"
			return
		}
		c <- ""
	}()
    go func() {
		res.Write(context.Background(), &rwp)
        r.Close()
	}()
	log.Println("wait")
	err := <-c
	if err != "" {
		t.Fatalf(err)
	}
	close(c)
}
