package http

import (
	"context"
	"errors"
	"fmt"

	"lwc.com/servergo/logger"
)

var routes = map[string]func(string) []byte{}

func route(ctx context.Context, sl *StartLine, ahs map[string]string, body *Body) {
	l := logger.Get(ctx)
    rf, ok := routes[sl.MPlusUrl]
}

func AddRoute(key string, f func(string) []byte) error {
	_, ok := routes[key]
	if ok {
		return errors.New(fmt.Sprintf("Route duplicated for key: %s", key))
	}
	routes[key] = f
    // add options version as well

	return nil
}
