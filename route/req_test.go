package route

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHeader(t *testing.T) {
	headers := map[string]string{
		"content-length": "123",
	}

	req := &Req{
		ahs: headers,
	}

	hv := req.GetHeader("CONTENT-length")
	assert.Equal(t, hv, "123")

	hv = req.GetHeader("CONTENT-LENGTH")
	assert.Equal(t, hv, "123")

	hv = req.GetHeader("content-length")
	assert.Equal(t, hv, "123")
}
