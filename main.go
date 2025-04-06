package main

import (
	"context"
	"fmt"
	"net"

	"github.com/google/uuid"
	"lwc.com/servergo/http"
	"lwc.com/servergo/logger"
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

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Printf("Error when accepting connection: %s", err.Error())
        }
        ctx := context.WithValue(timeoutCtx, logger.TRACE_ID_KEY, uuid.NewString())

        go http.HandleConn(ctx, conn)
    }
}
