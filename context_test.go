package hpack

import (
	"net/http"
	"reflect"
	"testing"
)

func TestRequestWithoutHuffman(t *testing.T) {

	client := NewContext()

	/**
	 * First Request
	 */
	buf := []byte{
		0x82, 0x87, 0x86, 0x04,
		0x0f, 0x77, 0x77, 0x77,
		0x2e, 0x65, 0x78, 0x61,
		0x6d, 0x70, 0x6c, 0x65,
		0x2e, 0x63, 0x6f, 0x6d,
	}

	expectedHeader := http.Header{
		"Method":    []string{"GET"},
		"Scheme":    []string{"http"},
		"Path":      []string{"/"},
		"Authority": []string{"www.example.com"},
	}

	expectedHeaderFields := []HeaderField{
		HeaderField{":authority", "www.example.com"},
		HeaderField{":path", "/"},
		HeaderField{":scheme", "http"},
		HeaderField{":method", "GET"},
	}

	client.Decode(buf)

	// test Header Table
	if client.HT.Size() != 180 {
		t.Errorf("got %v\nwant %v", client.HT.Size(), 180)
	}

	// test Header Table
	for i, hf := range expectedHeaderFields {
		if !(*client.HT.HeaderFields[i] == hf) {
			t.Errorf("got %v\nwant %v", *client.HT.HeaderFields[i], hf)
		}
	}

	// test Emitted Set
	if !reflect.DeepEqual(client.ES.Header, expectedHeader) {
		t.Errorf("got %v\nwant %v", client.ES.Header, expectedHeader)
	}

	// TOOD: test Reference Set
	for i, hf := range *client.RS {
		t.Log(i, hf)
	}

	/**
	 * Second Request
	 */
	buf = []byte{
		0x1b, 0x08, 0x6e, 0x6f,
		0x2d, 0x63, 0x61, 0x63,
		0x68, 0x65,
	}

	client.Decode(buf)

	expectedHeader = http.Header{
		"Method":        []string{"GET"},
		"Scheme":        []string{"http"},
		"Path":          []string{"/"},
		"Authority":     []string{"www.example.com"},
		"Cache-Control": []string{"no-cache"},
	}

	expectedHeaderFields = []HeaderField{
		HeaderField{"cache-control", "no-cache"},
		HeaderField{":authority", "www.example.com"},
		HeaderField{":path", "/"},
		HeaderField{":scheme", "http"},
		HeaderField{":method", "GET"},
	}

	// test Header Table
	if client.HT.Size() != 233 {
		t.Errorf("got %v\nwant %v", client.HT.Size(), 233)
	}

	// test Header Table
	for i, hf := range expectedHeaderFields {
		if !(*client.HT.HeaderFields[i] == hf) {
			t.Errorf("got %v\nwant %v", *client.HT.HeaderFields[i], hf)
		}
	}

	// TOOD: test Reference Set
	for i, hf := range *client.RS {
		t.Log(i, hf)
	}



}
