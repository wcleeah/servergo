package http

import "errors"

var headerEnds = errors.New("Headers Ended")
var unsupportedMethod = errors.New("Unsupported Method")
var unsupportedProtocol = errors.New("Unsupported Protocol")
var unsupportedProtocolVersion = errors.New("Unsupported Protocol Version")
