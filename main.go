package main

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/google/uuid"
	"lwc.com/servergo/internal/http"
	"lwc.com/servergo/internal/logger"
	"lwc.com/servergo/internal/route"
)

func main() {
	logger.Setup()
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		fmt.Printf("net.Listen failed, shutting down...")
		panic(err)
	}
	defer listener.Close()
	timeoutCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	l := logger.Get(timeoutCtx)
	l.Info("TCP listening on 3000")

    route.AddRoute("GET /health", func(req *route.Req, res *route.Res) {
        res.Write(&route.ResWriteParam{
            StatusCode: "200",
            Ahs: map[string]string{
                "Custom-Header": "hello",
            },
            Body: []byte("okokokokokokokok"),
        })
    })

    route.AddRoute("POST /user", func(req *route.Req, res *route.Res) {

		l := logger.Get(req.Ctx())
		body, err := io.ReadAll(req.Body())
		if err != nil {
			res.Write(&route.ResWriteParam{
				StatusCode: "400",
				Body: []byte("nobody nobody but you"),
			})
			return
		}

		l.Info("Body", "Body", string(body))
		res.Write(&route.ResWriteParam{
			StatusCode: "200",
			Ahs: map[string]string{
				"Custom-Header": "yoooooooooooo",
			},
			Body: []byte("useruseruseruser"),
		})
    })

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error when accepting connection: %s", err.Error())
		}
		ctx := context.WithValue(timeoutCtx, logger.TRACE_ID_KEY, uuid.NewString())
		handler := http.NewConnHandler(ctx)

		go handler.Handle(conn)
	}
}
