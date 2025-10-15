package http

import (
	"strings"

	"lwc.com/servergo/internal/common"
)

func trimCRLF(s string) string {
	return strings.Trim(s, common.CRLF)
}
