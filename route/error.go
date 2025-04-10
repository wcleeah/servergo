package route

import "errors"

var (
    ContentLengthNotSpecified = errors.New("Content length not specified in request")
    ContentLengthMalformed = errors.New("Content length malformed")
    BodyMalformed = errors.New("Body malformed")
)
