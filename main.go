package main

import (
	"context"
	"io"
	"net"

	"lwc.com/servergo/internal/logger"
	"lwc.com/servergo/pkg/server"
)

func main() {
	ctx := context.Background()

	logger.Setup()
	l := logger.Get(ctx)

	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		l.Error("net.Listen failed, shutting down...")
		panic(err)
	}
	defer listener.Close()
	l.Info("TCP listening on 3000")

	server.AddRoute("GET /health", func(req *server.Req, res *server.Res) {
		res.Write(&server.ResWriteParam{
			StatusCode: "200",
			Ahs: map[string]string{
				"Custom-Header": "hello",
			},
			Body: []byte("okokokokokokokok"),
		})
	})

	server.AddRoute("POST /user", func(req *server.Req, res *server.Res) {

		l := logger.Get(req.Ctx())
		body, err := io.ReadAll(req.Body())
		if err != nil {
			res.Write(&server.ResWriteParam{
				StatusCode: "400",
				Body:       []byte("nobody nobody but you"),
			})
			return
		}

		l.Info("Body", "Body", string(body))
		res.Write(&server.ResWriteParam{
			StatusCode: "200",
			Ahs: map[string]string{
				"Custom-Header": "yoooooooooooo",
			},
			Body: []byte("useruseruseruser"),
		})
	})

	server := server.NewHTTPServer(listener)

	server.Start(ctx)
}
