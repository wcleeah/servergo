# ServerGo
A experimental server so that i can understand how a http server works, it starts a tcp server, and a hand crafted http server. 
It is absolutely not for production, it is only for me to try building a http server in go.

## Product
- A HTTP server from scratch (TCP server -> response)

## TODO
- [x] parse body
- [x] unit test for route.Req and route.Res 
- [x] assert / error handling
- [x] add route and see if the whole thing works 
- [x] some more unit testing
- [x] 1.1 support
  - [x] Connection: Keep-Alive

## Future
- [ ] 1.1 support
  - [ ] Mandatory Host header
  - [ ] Chunked Transfer Encoding
- [ ] load testing
- [ ] full compatible with the specification
- [ ] http 2
- [ ] websocket
- [ ] JSON parser
