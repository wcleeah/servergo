package server 

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHeaderIgnoreCasing(t *testing.T) {
	headers := map[string]string{
		"content-length": "123",
	}

	req := &Req{
		ahs: headers,
	}

	hv, ok := req.GetHeader("CONTENT-length")
	assert.Equal(t, hv, "123")
	assert.True(t, ok)

	hv, ok = req.GetHeader("CONTENT-LENGTH")
	assert.Equal(t, hv, "123")
	assert.True(t, ok)

	hv, ok = req.GetHeader("content-length")
	assert.Equal(t, hv, "123")
	assert.True(t, ok)
}

func TestGetHeaderNoHeader(t *testing.T) {
	headers := map[string]string{
		"content-length": "123",
	}

	req := &Req{
		ahs: headers,
	}

	hv, ok := req.GetHeader("haha")
	assert.Empty(t, hv)
	assert.False(t, ok)
}
