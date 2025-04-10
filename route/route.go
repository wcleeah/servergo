package route

import (
	"context"
	"errors"
	"fmt"

	"lwc.com/servergo/logger"
)

// Add mutex...?
var routes = map[string]func(req *Req, res *Res){}

func Route(ctx context.Context, req *Req, res *Res) error {
	l := logger.Get(ctx)
	// assert method and url
	mPlusUrl := req.Method + " " + req.Url
	rf, ok := routes[mPlusUrl]
	if !ok {
		l.Info("No route found for method and url", "method", req.Method, "url", req.Url)
		return nil
	}
	rf(req, res)
	return nil
}

func AddRoute(key string, f func(req *Req, res *Res)) error {
	// assert key format
	_, ok := routes[key]
	if ok {
		return errors.New(fmt.Sprintf("Route duplicated for key: %s", key))
	}
	routes[key] = f
	// add options version as well

	return nil
}
