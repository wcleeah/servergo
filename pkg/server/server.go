package server

import (
	"context"

	"github.com/google/uuid"
	"lwc.com/servergo/internal/common"
	"lwc.com/servergo/internal/logger"
)

type TCPServer struct {
	listener common.Listener
}

func NewHTTPServer(listener common.Listener) TCPServer {
	return TCPServer{
		listener: listener,
	}
}

func (s *TCPServer) Start(ctx context.Context) error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		ctx := context.WithValue(ctx, logger.TRACE_ID_KEY, uuid.NewString())
		handler := newConnHandler(ctx)

		go handler.Handle(conn)
	}
}
