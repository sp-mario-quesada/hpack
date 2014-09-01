package hpack

import (
	"fmt"
	"net/http"
	"strings"
)

// A header set is a potentially ordered group of header fields that are encoded jointly.
// A complete set of key-value pairs contained in a HTTP request or response is a header set.
type HeaderList []HeaderField

func NewHeaderList() *HeaderList {
	return new(HeaderList)
}

func ToHeaderList(header http.Header) HeaderList {
	hl := *new(HeaderList)
	for key, values := range header {
		key := strings.ToLower(key)
		for _, value := range values {
			hl = append(hl, *NewHeaderField(key, value))
		}
	}
	return hl
}

func (hl *HeaderList) Emit(hf *HeaderField) {
	*hl = append(*hl, *hf)
}

func (hl *HeaderList) Len() int {
	return len(*hl)
}

// Sort Interface
func (hl *HeaderList) Swap(i, j int) {
	h := *hl
	h[i], h[j] = h[j], h[i]
}

func (hl *HeaderList) Less(i, j int) bool {
	h := *hl
	if h[i].Name == h[j].Name {
		return h[i].Value < h[j].Value
	}
	return h[i].Name < h[j].Name
}

// convert to http.Header
func (hl HeaderList) ToHeader() http.Header {
	header := make(http.Header)
	for _, hf := range hl {
		header.Add(hf.Name, hf.Value)
	}
	return header
}

func (hl HeaderList) String() (str string) {
	str += fmt.Sprintf("\n--------- HL ---------\n")
	for i, v := range hl {
		str += fmt.Sprintln(i, v)
	}
	str += "--------------------------------\n"
	return str
}
