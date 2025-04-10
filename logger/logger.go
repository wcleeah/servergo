package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

const TRACE_ID_KEY = "traceId"

func Setup() {
	l := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(l)
}

func Get(ctx context.Context) *slog.Logger {
	id := ctx.Value(TRACE_ID_KEY)
	if id == nil {
		id = uuid.NewString()
	}
	return slog.With(TRACE_ID_KEY, id)
}
