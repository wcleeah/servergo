package route

import (
	"errors"
	"fmt"
)

// Add mutex...?
var routes = map[string]func(req *Req, res *Res){}

func Route(req *Req, res *Res) {
	// assert method and url
	mPlusUrl := req.Method + " " + req.Url
	rf, ok := routes[mPlusUrl]
	if !ok {
		res.Write(&ResWriteParam{
			StatusCode: "404",
			Body:       []byte("Wrong url / method?"),
		})
	}
	rf(req, res)
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
