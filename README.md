# ServerGo
A experimental server so that i can understand how a http server works, it starts a tcp server, and some hand craft http parsing. 
It is absolutely not for production, it is only for me to try building a http server in go.

## Product
- A HTTP server from scratch (TCP server -> response)
- CORS
- based on region and mode return a restaurant name

## TODO
- [x] parse body
- [x] unit test for route.Req and route.Res 
- [x] assert / error handling
- [x] add route and see if the whole thing works 
- [x] some more unit testing
- [ ] 1.1 support
  - [ ] Connection: Keep-Alive
  - [ ] Mandatory Host header
  - [ ] Chunked Transfer Encoding
- [ ] 2 support
- [ ] CORS support
- [ ] load test
- [ ] add the route i need to use lol

## Future
- [ ] full compatible with the specification
- [ ] streaming? (but http/2 already kind of does that)
- [ ] websocket
