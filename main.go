package main

import (
	"context"
	"io"
	"net"

	"lwc.com/servergo/internal/logger"
	"lwc.com/servergo/internal/route"
	"lwc.com/servergo/internal/server"
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
				Body:       []byte("nobody nobody but you"),
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

	server := server.NewTCPServer(listener)

	server.Start(ctx)
}
