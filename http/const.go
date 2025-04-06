package http

var SUPPORTED_HEADERS = []string{
    "Content-Length",
}

var SUPPORTED_METHODS = []string{
    "POST",
    "GET",
    // TODO: support cors
    // "OPTIONS",
}

var SUPPORTED_PROTOCOL_VERSION = []string{
    "1.0",
    "1.1",
}


const (
    SUPPORTED_PROTOCOL = "HTTP"
)
