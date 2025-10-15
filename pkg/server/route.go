package server 

import (
	"fmt"
)

type RouteFunc = func(req *Req, res *Res)

// Add mutex...?
var routes = map[string]RouteFunc{}

func route(req *Req, res *Res) {
	// assert method and url
	mPlusUrl := req.Method + " " + req.Url
	rf, ok := routes[mPlusUrl]
	if !ok {
		res.Write(&ResWriteParam{
			StatusCode: "404",
			Body:       []byte("Wrong url / method?"),
		})
		return
	}
	rf(req, res)
}

func AddRoute(key string, f RouteFunc) error {
	// assert key format
	_, ok := routes[key]
	if ok {
		return fmt.Errorf("Route duplicated for key: %s", key)
	}
	routes[key] = f
	// add options version as well

	return nil
}
