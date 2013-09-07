package hpac

import (
	"net/http"
	"strings"
)

type HeaderSet map[string]string

var MustHeader = map[string]string{
	"Scheme": ":scheme",
	"Method": ":method",
	"Path":   ":path",
	"Host":   ":host",
	"Status": ":status",
}

func NewHeaderSet(header http.Header) HeaderSet {
	// method, scheme, host, path, status
	headerSet := make(HeaderSet, len(header))
	for name, value := range header {
		mustname, ok := MustHeader[name]
		if ok {
			name = mustname
		}
		headerSet[name] = strings.Join(value, ",")
	}
	return headerSet
}
