package http

import "strings"

func trimCRLF(s string) string {
    return strings.Trim(s, "\r\n")
}
