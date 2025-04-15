# Http server from scratch

## HTTP vs TCP
- HTTP operates on Layer 7
- TCP operates on Layer 4
- TCP controls to accept and send packet
- HTTP defines whats inside a packet

## HTTP Packet
- there are three parts
  - the start line (method, url, status etc)
  - the headers
  - the body
- none of the above is fixed size, parser uses \r\n to determine which part is which
- there is usually a content length header so that the parser knows the body size

### Headers
- CRLF is not allowed (99.9% of use case) in header value
- Content-Length must be in ASCII format, it cannot be binary represented

## Some reflections
- It is quite hard to design a strategy to handle the request and response
  - at last i just follow the go standard library style
- The design comes very natuarally, it kinda looks like the go std lib http mux now

## HTTP 1.1: main difference
- host header is now required
- persistent connection is now supported
  - by using Connection header and default timeout
- chunked transfer encoding is now supported
