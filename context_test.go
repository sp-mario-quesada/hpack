package hpack

import (
	assert "github.com/jxck/assertion"
	. "github.com/jxck/color"
	. "github.com/jxck/logger"
	"sort"
	"testing"
)

func TestRequestWithoutHuffman(t *testing.T) {
	// D.2.  Request Examples without Huffman
	var (
		context    *Context
		buf        []byte
		expectedES *HeaderSet
		expectedHT *HeaderTable
		expectedRS ReferenceSet
	)

	context = NewContext(DEFAULT_HEADER_TABLE_SIZE)

	/**
	 * D.2.1.  First request
	 */
	Debug(Pink("\n========== First Request ==============="))

	buf = []byte{
		0x82, 0x87,
		0x86, 0x04,
		0x0f, 0x77,
		0x77, 0x77,
		0x2e, 0x65,
		0x78, 0x61,
		0x6d, 0x70,
		0x6c, 0x65,
		0x2e, 0x63,
		0x6f, 0x6d,
	}

	expectedES = &HeaderSet{
		HeaderField{":method", "GET"},
		HeaderField{":scheme", "http"},
		HeaderField{":path", "/"},
		HeaderField{":authority", "www.example.com"},
	}

	expectedHT = NewHeaderTable(DEFAULT_HEADER_TABLE_SIZE)
	expectedHT.HeaderFields = []*HeaderField{
		NewHeaderField(":authority", "www.example.com"),
		NewHeaderField(":path", "/"),
		NewHeaderField(":scheme", "http"),
		NewHeaderField(":method", "GET"),
	}

	expectedRS = ReferenceSet{
		{NewHeaderField(":authority", "www.example.com"), false},
		{NewHeaderField(":path", "/"), false},
		{NewHeaderField(":scheme", "http"), false},
		{NewHeaderField(":method", "GET"), false},
	}

	context.Decode(buf)

	// test Header Table
	assert.Equal(t, context.HT, expectedHT)

	// test Emitted Set
	sort.Sort(context.ES)
	sort.Sort(expectedES)
	assert.Equal(t, context.ES, expectedES)

	// test Reference Set
	if expectedRS.Len() != context.RS.Len() {
		t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
	}

	for i, hf := range *context.RS {
		i := len(expectedRS) - i - 1
		if *expectedRS[i].HeaderField != *hf.HeaderField {
			t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
		}
	}

	/**
	 * D.2.2.  Second request
	 */
	Debug(Pink("\n========== Second Request ==============="))

	buf = []byte{
		0x1b, 0x08,
		0x6e, 0x6f,
		0x2d, 0x63,
		0x61, 0x63,
		0x68, 0x65,
	}

	expectedES = &HeaderSet{
		HeaderField{":method", "GET"},
		HeaderField{":scheme", "http"},
		HeaderField{":path", "/"},
		HeaderField{":authority", "www.example.com"},
		HeaderField{"cache-control", "no-cache"},
	}

	expectedHT = NewHeaderTable(DEFAULT_HEADER_TABLE_SIZE)
	expectedHT.HeaderFields = []*HeaderField{
		NewHeaderField("cache-control", "no-cache"),
		NewHeaderField(":authority", "www.example.com"),
		NewHeaderField(":path", "/"),
		NewHeaderField(":scheme", "http"),
		NewHeaderField(":method", "GET"),
	}

	expectedRS = ReferenceSet{
		{NewHeaderField("cache-control", "no-cache"), false},
		{NewHeaderField(":authority", "www.example.com"), false},
		{NewHeaderField(":path", "/"), false},
		{NewHeaderField(":scheme", "http"), false},
		{NewHeaderField(":method", "GET"), false},
	}

	context.Decode(buf)

	// test Header Table
	assert.Equal(t, context.HT, expectedHT)

	// test Emitted Set
	sort.Sort(context.ES)
	sort.Sort(expectedES)
	assert.Equal(t, context.ES, expectedES)

	// test Reference Set
	if expectedRS.Len() != context.RS.Len() {
		t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
	}

	for i, hf := range *context.RS {
		i := len(expectedRS) - i - 1
		if *expectedRS[i].HeaderField != *hf.HeaderField {
			t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
		}
	}

	/**
	 * D.2.3.  Third request
	 */
	Debug(Pink("\n========== Third Request ==============="))

	buf = []byte{
		0x80, 0x80,
		0x85, 0x8c,
		0x8b, 0x84,
		0x00, 0x0a,
		0x63, 0x75,
		0x73, 0x74,
		0x6f, 0x6d,
		0x2d, 0x6b,
		0x65, 0x79,
		0x0c, 0x63,
		0x75, 0x73,
		0x74, 0x6f,
		0x6d, 0x2d,
		0x76, 0x61,
		0x6c, 0x75,
		0x65,
	}

	expectedES = &HeaderSet{
		HeaderField{":method", "GET"},
		HeaderField{":scheme", "https"},
		HeaderField{":path", "/index.html"},
		HeaderField{":authority", "www.example.com"},
		HeaderField{"custom-key", "custom-value"},
	}

	expectedHT = NewHeaderTable(DEFAULT_HEADER_TABLE_SIZE)
	expectedHT.HeaderFields = []*HeaderField{
		NewHeaderField("custom-key", "custom-value"),
		NewHeaderField(":path", "/index.html"),
		NewHeaderField(":scheme", "https"),
		NewHeaderField("cache-control", "no-cache"),
		NewHeaderField(":authority", "www.example.com"),
		NewHeaderField(":path", "/"),
		NewHeaderField(":scheme", "http"),
		NewHeaderField(":method", "GET"),
	}

	expectedRS = ReferenceSet{
		{NewHeaderField("custom-key", "custom-value"), true},
		{NewHeaderField(":authority", "www.example.com"), true},
		{NewHeaderField(":path", "/index.html"), true},
		{NewHeaderField(":scheme", "https"), true},
		{NewHeaderField(":method", "GET"), true},
	}

	context.Decode(buf)

	// test Header Table
	assert.Equal(t, context.HT, expectedHT)

	// test Emitted Set
	sort.Sort(context.ES)
	sort.Sort(expectedES)
	assert.Equal(t, context.ES, expectedES)

	// test Reference Set
	if expectedRS.Len() != context.RS.Len() {
		t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
	}

	for i, hf := range *context.RS {
		i := len(expectedRS) - i - 1
		if *expectedRS[i].HeaderField != *hf.HeaderField {
			t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
		}
	}
}

func TestRequestWithHuffman(t *testing.T) {
	// D.3.  Request Examples with Huffman
	var (
		context    *Context
		buf        []byte
		expectedES *HeaderSet
		expectedHT *HeaderTable
		expectedRS ReferenceSet
	)

	context = NewContext(DEFAULT_HEADER_TABLE_SIZE)

	/**
	 * D.3.1.  First request
	 */
	Debug(Pink("\n========== First Request ==============="))

	buf = []byte{
		0x82, 0x87,
		0x86, 0x04,
		0x8b, 0xdb,
		0x6d, 0x88,
		0x3e, 0x68,
		0xd1, 0xcb,
		0x12, 0x25,
		0xba, 0x7f,
	}

	expectedES = &HeaderSet{
		HeaderField{":method", "GET"},
		HeaderField{":scheme", "http"},
		HeaderField{":path", "/"},
		HeaderField{":authority", "www.example.com"},
	}

	expectedHT = NewHeaderTable(DEFAULT_HEADER_TABLE_SIZE)
	expectedHT.HeaderFields = []*HeaderField{
		NewHeaderField(":authority", "www.example.com"),
		NewHeaderField(":path", "/"),
		NewHeaderField(":scheme", "http"),
		NewHeaderField(":method", "GET"),
	}

	expectedRS = ReferenceSet{
		{NewHeaderField(":authority", "www.example.com"), false},
		{NewHeaderField(":path", "/"), false},
		{NewHeaderField(":scheme", "http"), false},
		{NewHeaderField(":method", "GET"), false},
	}

	context.Decode(buf)

	// test Header Table
	assert.Equal(t, context.HT, expectedHT)

	// test Emitted Set
	sort.Sort(context.ES)
	sort.Sort(expectedES)
	assert.Equal(t, context.ES, expectedES)

	// test Reference Set
	if expectedRS.Len() != context.RS.Len() {
		t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
	}

	for i, hf := range *context.RS {
		i := len(expectedRS) - i - 1
		if *expectedRS[i].HeaderField != *hf.HeaderField {
			t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
		}
	}

	/**
	 * D.3.2.  Second request
	 */
	Debug(Pink("\n========== Second Request ==============="))

	buf = []byte{
		0x1b, 0x86,
		0x63, 0x65,
		0x4a, 0x13,
		0x98, 0xff,
	}

	expectedES = &HeaderSet{
		HeaderField{":method", "GET"},
		HeaderField{":scheme", "http"},
		HeaderField{":path", "/"},
		HeaderField{":authority", "www.example.com"},
		HeaderField{"cache-control", "no-cache"},
	}

	expectedHT = NewHeaderTable(DEFAULT_HEADER_TABLE_SIZE)
	expectedHT.HeaderFields = []*HeaderField{
		NewHeaderField("cache-control", "no-cache"),
		NewHeaderField(":authority", "www.example.com"),
		NewHeaderField(":path", "/"),
		NewHeaderField(":scheme", "http"),
		NewHeaderField(":method", "GET"),
	}

	expectedRS = ReferenceSet{
		{NewHeaderField("cache-control", "no-cache"), false},
		{NewHeaderField(":authority", "www.example.com"), false},
		{NewHeaderField(":path", "/"), false},
		{NewHeaderField(":scheme", "http"), false},
		{NewHeaderField(":method", "GET"), false},
	}

	context.Decode(buf)

	// test Header Table
	assert.Equal(t, context.HT, expectedHT)

	// test Emitted Set
	sort.Sort(context.ES)
	sort.Sort(expectedES)
	assert.Equal(t, context.ES, expectedES)

	// test Reference Set
	if expectedRS.Len() != context.RS.Len() {
		t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
	}

	for i, hf := range *context.RS {
		i := len(expectedRS) - i - 1
		if *expectedRS[i].HeaderField != *hf.HeaderField {
			t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
		}
	}

	/**
	 * D.3.3.  Third request
	 */
	Debug(Pink("\n========== Third Request ==============="))

	buf = []byte{
		0x80, 0x80,
		0x85, 0x8c,
		0x8b, 0x84,
		0x00, 0x88,
		0x4e, 0xb0,
		0x8b, 0x74,
		0x97, 0x90,
		0xfa, 0x7f,
		0x89, 0x4e,
		0xb0, 0x8b,
		0x74, 0x97,
		0x9a, 0x17,
		0xa8, 0xff,
	}

	expectedES = &HeaderSet{
		HeaderField{":method", "GET"},
		HeaderField{":scheme", "https"},
		HeaderField{":path", "/index.html"},
		HeaderField{":authority", "www.example.com"},
		HeaderField{"custom-key", "custom-value"},
	}

	expectedHT = NewHeaderTable(DEFAULT_HEADER_TABLE_SIZE)
	expectedHT.HeaderFields = []*HeaderField{
		NewHeaderField("custom-key", "custom-value"),
		NewHeaderField(":path", "/index.html"),
		NewHeaderField(":scheme", "https"),
		NewHeaderField("cache-control", "no-cache"),
		NewHeaderField(":authority", "www.example.com"),
		NewHeaderField(":path", "/"),
		NewHeaderField(":scheme", "http"),
		NewHeaderField(":method", "GET"),
	}

	expectedRS = ReferenceSet{
		{NewHeaderField("custom-key", "custom-value"), true},
		{NewHeaderField(":authority", "www.example.com"), true},
		{NewHeaderField(":path", "/index.html"), true},
		{NewHeaderField(":scheme", "https"), true},
		{NewHeaderField(":method", "GET"), true},
	}

	context.Decode(buf)

	// test Header Table
	assert.Equal(t, context.HT, expectedHT)

	// test Emitted Set
	sort.Sort(context.ES)
	sort.Sort(expectedES)
	assert.Equal(t, context.ES, expectedES)

	// test Reference Set
	if expectedRS.Len() != context.RS.Len() {
		t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
	}

	for i, hf := range *context.RS {
		i := len(expectedRS) - i - 1
		if *expectedRS[i].HeaderField != *hf.HeaderField {
			t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
		}
	}
}

func TestResponseWithoutHuffman(t *testing.T) {
	// D.4.  Response Examples without Huffman
	var (
		context    *Context
		buf        []byte
		expectedES *HeaderSet
		expectedHT *HeaderTable
		expectedRS ReferenceSet
	)

	var HeaderTableSize int = 256
	context = NewContext(HeaderTableSize)

	/**
	 * D.4.1.  First response
	 */
	Debug(Pink("\n========== First Response ==============="))

	buf = []byte{
		0x08, 0x03,
		0x33, 0x30,
		0x32, 0x18,
		0x07, 0x70,
		0x72, 0x69,
		0x76, 0x61,
		0x74, 0x65,
		0x22, 0x1d,
		0x4d, 0x6f,
		0x6e, 0x2c,
		0x20, 0x32,
		0x31, 0x20,
		0x4f, 0x63,
		0x74, 0x20,
		0x32, 0x30,
		0x31, 0x33,
		0x20, 0x32,
		0x30, 0x3a,
		0x31, 0x33,
		0x3a, 0x32,
		0x31, 0x20,
		0x47, 0x4d,
		0x54, 0x30,
		0x17, 0x68,
		0x74, 0x74,
		0x70, 0x73,
		0x3a, 0x2f,
		0x2f, 0x77,
		0x77, 0x77,
		0x2e, 0x65,
		0x78, 0x61,
		0x6d, 0x70,
		0x6c, 0x65,
		0x2e, 0x63,
		0x6f, 0x6d,
	}

	expectedES = &HeaderSet{
		HeaderField{":status", "302"},
		HeaderField{"cache-control", "private"},
		HeaderField{"date", "Mon, 21 Oct 2013 20:13:21 GMT"},
		HeaderField{"location", "https://www.example.com"},
	}

	expectedHT = NewHeaderTable(HeaderTableSize)
	expectedHT.HeaderFields = []*HeaderField{
		NewHeaderField("location", "https://www.example.com"),
		NewHeaderField("date", "Mon, 21 Oct 2013 20:13:21 GMT"),
		NewHeaderField("cache-control", "private"),
		NewHeaderField(":status", "302"),
	}

	expectedRS = ReferenceSet{
		{NewHeaderField("location", "https://www.example.com"), false},
		{NewHeaderField("date", "Mon, 21 Oct 2013 20:13:21 GMT"), false},
		{NewHeaderField("cache-control", "private"), false},
		{NewHeaderField(":status", "302"), false},
	}

	context.Decode(buf)

	// test Header Table
	assert.Equal(t, context.HT, expectedHT)

	// test Emitted Set
	sort.Sort(context.ES)
	sort.Sort(expectedES)
	assert.Equal(t, context.ES, expectedES)

	// test Reference Set
	if expectedRS.Len() != context.RS.Len() {
		t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
	}

	for i, hf := range *context.RS {
		i := len(expectedRS) - i - 1
		if *expectedRS[i].HeaderField != *hf.HeaderField {
			t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
		}
	}

	/**
	 * D.4.2.  Second response
	 */
	Debug(Pink("\n========== Second Response ==============="))

	buf = []byte{
		0x84, 0x8c,
	}

	expectedES = &HeaderSet{
		HeaderField{":status", "200"},
		HeaderField{"cache-control", "private"},
		HeaderField{"date", "Mon, 21 Oct 2013 20:13:21 GMT"},
		HeaderField{"location", "https://www.example.com"},
	}

	expectedHT = NewHeaderTable(HeaderTableSize)
	expectedHT.HeaderFields = []*HeaderField{
		NewHeaderField(":status", "200"),
		NewHeaderField("location", "https://www.example.com"),
		NewHeaderField("date", "Mon, 21 Oct 2013 20:13:21 GMT"),
		NewHeaderField("cache-control", "private"),
	}

	expectedRS = ReferenceSet{
		{NewHeaderField(":status", "200"), false},
		{NewHeaderField("location", "https://www.example.com"), false},
		{NewHeaderField("date", "Mon, 21 Oct 2013 20:13:21 GMT"), false},
		{NewHeaderField("cache-control", "private"), false},
	}

	context.Decode(buf)

	// test Header Table
	assert.Equal(t, context.HT, expectedHT)

	// test Emitted Set
	sort.Sort(context.ES)
	sort.Sort(expectedES)
	assert.Equal(t, context.ES, expectedES)

	// test Reference Set
	if expectedRS.Len() != context.RS.Len() {
		t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
	}

	for i, hf := range *context.RS {
		i := len(expectedRS) - i - 1
		if *expectedRS[i].HeaderField != *hf.HeaderField {
			t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
		}
	}

	/**
	 * D.4.3.  Third response
	 */
	Debug(Pink("\n========== Third Response ==============="))

	buf = []byte{
		0x83, 0x84,
		0x84, 0x03,
		0x1d, 0x4d,
		0x6f, 0x6e,
		0x2c, 0x20,
		0x32, 0x31,
		0x20, 0x4f,
		0x63, 0x74,
		0x20, 0x32,
		0x30, 0x31,
		0x33, 0x20,
		0x32, 0x30,
		0x3a, 0x31,
		0x33, 0x3a,
		0x32, 0x32,
		0x20, 0x47,
		0x4d, 0x54,
		0x1d, 0x04,
		0x67, 0x7a,
		0x69, 0x70,
		0x84, 0x84,
		0x83, 0x83,
		0x3a, 0x38,
		0x66, 0x6f,
		0x6f, 0x3d,
		0x41, 0x53,
		0x44, 0x4a,
		0x4b, 0x48,
		0x51, 0x4b,
		0x42, 0x5a,
		0x58, 0x4f,
		0x51, 0x57,
		0x45, 0x4f,
		0x50, 0x49,
		0x55, 0x41,
		0x58, 0x51,
		0x57, 0x45,
		0x4f, 0x49,
		0x55, 0x3b,
		0x20, 0x6d,
		0x61, 0x78,
		0x2d, 0x61,
		0x67, 0x65,
		0x3d, 0x33,
		0x36, 0x30,
		0x30, 0x3b,
		0x20, 0x76,
		0x65, 0x72,
		0x73, 0x69,
		0x6f, 0x6e,
		0x3d, 0x31,
	}

	expectedES = &HeaderSet{
		HeaderField{":status", "200"},
		HeaderField{"cache-control", "private"},
		HeaderField{"date", "Mon, 21 Oct 2013 20:13:22 GMT"},
		HeaderField{"location", "https://www.example.com"},
		HeaderField{"content-encoding", "gzip"},
		HeaderField{"set-cookie", "foo=ASDJKHQKBZXOQWEOPIUAXQWEOIU; max-age=3600; version=1"},
	}

	expectedHT = NewHeaderTable(HeaderTableSize)
	expectedHT.HeaderFields = []*HeaderField{
		NewHeaderField("set-cookie", "foo=ASDJKHQKBZXOQWEOPIUAXQWEOIU; max-age=3600; version=1"),
		NewHeaderField("content-encoding", "gzip"),
		NewHeaderField("date", "Mon, 21 Oct 2013 20:13:22 GMT"),
	}

	expectedRS = ReferenceSet{
		{NewHeaderField("set-cookie", "foo=ASDJKHQKBZXOQWEOPIUAXQWEOIU; max-age=3600; version=1"), false},
		{NewHeaderField("content-encoding", "gzip"), false},
		{NewHeaderField("date", "Mon, 21 Oct 2013 20:13:22 GMT"), false},
	}

	context.Decode(buf)

	// test Header Table
	assert.Equal(t, context.HT, expectedHT)

	// test Emitted Set
	sort.Sort(context.ES)
	sort.Sort(expectedES)
	assert.Equal(t, context.ES, expectedES)

	// test Reference Set
	if expectedRS.Len() != context.RS.Len() {
		t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
	}

	for i, hf := range *context.RS {
		i := len(expectedRS) - i - 1
		if *expectedRS[i].HeaderField != *hf.HeaderField {
			t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
		}
	}
}

func TestResponseWithHuffman(t *testing.T) {
	// D.5.  Response Examples with Huffman
	var (
		context    *Context
		buf        []byte
		expectedES *HeaderSet
		expectedHT *HeaderTable
		expectedRS ReferenceSet
	)

	var HeaderTableSize int = 256
	context = NewContext(HeaderTableSize)

	/**
	 * D.5.1.  First response
	 */
	Debug(Pink("\n========== First Response ==============="))

	buf = []byte{
		0x08, 0x03,
		0x33, 0x30,
		0x32, 0x18,
		0x07, 0x70,
		0x72, 0x69,
		0x76, 0x61,
		0x74, 0x65,
		0x22, 0x1d,
		0x4d, 0x6f,
		0x6e, 0x2c,
		0x20, 0x32,
		0x31, 0x20,
		0x4f, 0x63,
		0x74, 0x20,
		0x32, 0x30,
		0x31, 0x33,
		0x20, 0x32,
		0x30, 0x3a,
		0x31, 0x33,
		0x3a, 0x32,
		0x31, 0x20,
		0x47, 0x4d,
		0x54, 0x30,
		0x17, 0x68,
		0x74, 0x74,
		0x70, 0x73,
		0x3a, 0x2f,
		0x2f, 0x77,
		0x77, 0x77,
		0x2e, 0x65,
		0x78, 0x61,
		0x6d, 0x70,
		0x6c, 0x65,
		0x2e, 0x63,
		0x6f, 0x6d,
	}

	expectedES = &HeaderSet{
		HeaderField{":status", "302"},
		HeaderField{"cache-control", "private"},
		HeaderField{"date", "Mon, 21 Oct 2013 20:13:21 GMT"},
		HeaderField{"location", "https://www.example.com"},
	}

	expectedHT = NewHeaderTable(HeaderTableSize)
	expectedHT.HeaderFields = []*HeaderField{
		NewHeaderField("location", "https://www.example.com"),
		NewHeaderField("date", "Mon, 21 Oct 2013 20:13:21 GMT"),
		NewHeaderField("cache-control", "private"),
		NewHeaderField(":status", "302"),
	}

	expectedRS = ReferenceSet{
		{NewHeaderField("location", "https://www.example.com"), false},
		{NewHeaderField("date", "Mon, 21 Oct 2013 20:13:21 GMT"), false},
		{NewHeaderField("cache-control", "private"), false},
		{NewHeaderField(":status", "302"), false},
	}

	context.Decode(buf)

	// test Header Table
	assert.Equal(t, context.HT, expectedHT)

	// test Emitted Set
	sort.Sort(context.ES)
	sort.Sort(expectedES)
	assert.Equal(t, context.ES, expectedES)

	// test Reference Set
	if expectedRS.Len() != context.RS.Len() {
		t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
	}

	for i, hf := range *context.RS {
		i := len(expectedRS) - i - 1
		if *expectedRS[i].HeaderField != *hf.HeaderField {
			t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
		}
	}

	/**
	 * D.5.2.  Second response
	 */
	Debug(Pink("\n========== Second Response ==============="))

	buf = []byte{
		0x8c,
	}

	expectedES = &HeaderSet{
		HeaderField{":status", "200"},
		HeaderField{"cache-control", "private"},
		HeaderField{"date", "Mon, 21 Oct 2013 20:13:21 GMT"},
		HeaderField{"location", "https://www.example.com"},
	}

	expectedHT = NewHeaderTable(HeaderTableSize)
	expectedHT.HeaderFields = []*HeaderField{
		NewHeaderField(":status", "200"),
		NewHeaderField("location", "https://www.example.com"),
		NewHeaderField("date", "Mon, 21 Oct 2013 20:13:21 GMT"),
		NewHeaderField("cache-control", "private"),
	}

	expectedRS = ReferenceSet{
		{NewHeaderField(":status", "200"), false},
		{NewHeaderField("location", "https://www.example.com"), false},
		{NewHeaderField("date", "Mon, 21 Oct 2013 20:13:21 GMT"), false},
		{NewHeaderField("cache-control", "private"), false},
	}

	context.Decode(buf)

	// test Header Table
	assert.Equal(t, context.HT, expectedHT)

	// test Emitted Set
	sort.Sort(context.ES)
	sort.Sort(expectedES)
	assert.Equal(t, context.ES, expectedES)

	// test Reference Set
	if expectedRS.Len() != context.RS.Len() {
		t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
	}

	for i, hf := range *context.RS {
		i := len(expectedRS) - i - 1
		if *expectedRS[i].HeaderField != *hf.HeaderField {
			t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
		}
	}

	/**
	 * D.5.3.  Third response
	 */
	Debug(Pink("\n========== Third Response ==============="))

	buf = []byte{
		0x84, 0x84,
		0x03, 0x1d,
		0x4d, 0x6f,
		0x6e, 0x2c,
		0x20, 0x32,
		0x31, 0x20,
		0x4f, 0x63,
		0x74, 0x20,
		0x32, 0x30,
		0x31, 0x33,
		0x20, 0x32,
		0x30, 0x3a,
		0x31, 0x33,
		0x3a, 0x32,
		0x32, 0x20,
		0x47, 0x4d,
		0x54, 0x1d,
		0x04, 0x67,
		0x7a, 0x69,
		0x70, 0x84,
		0x84, 0x83,
		0x83, 0x3a,
		0x38, 0x66,
		0x6f, 0x6f,
		0x3d, 0x41,
		0x53, 0x44,
		0x4a, 0x4b,
		0x48, 0x51,
		0x4b, 0x42,
		0x5a, 0x58,
		0x4f, 0x51,
		0x57, 0x45,
		0x4f, 0x50,
		0x49, 0x55,
		0x41, 0x58,
		0x51, 0x57,
		0x45, 0x4f,
		0x49, 0x55,
		0x3b, 0x20,
		0x6d, 0x61,
		0x78, 0x2d,
		0x61, 0x67,
		0x65, 0x3d,
		0x33, 0x36,
		0x30, 0x30,
		0x3b, 0x20,
		0x76, 0x65,
		0x72, 0x73,
		0x69, 0x6f,
		0x6e, 0x3d,
		0x31,
	}

	expectedES = &HeaderSet{
		HeaderField{":status", "200"},
		HeaderField{"cache-control", "private"},
		HeaderField{"date", "Mon, 21 Oct 2013 20:13:22 GMT"},
		HeaderField{"location", "https://www.example.com"},
		HeaderField{"content-encoding", "gzip"},
		HeaderField{"set-cookie", "foo=ASDJKHQKBZXOQWEOPIUAXQWEOIU; max-age=3600; version=1"},
	}

	expectedHT = NewHeaderTable(HeaderTableSize)
	expectedHT.HeaderFields = []*HeaderField{
		NewHeaderField("set-cookie", "foo=ASDJKHQKBZXOQWEOPIUAXQWEOIU; max-age=3600; version=1"),
		NewHeaderField("content-encoding", "gzip"),
		NewHeaderField("date", "Mon, 21 Oct 2013 20:13:22 GMT"),
	}

	expectedRS = ReferenceSet{
		{NewHeaderField("set-cookie", "foo=ASDJKHQKBZXOQWEOPIUAXQWEOIU; max-age=3600; version=1"), false},
		{NewHeaderField("content-encoding", "gzip"), false},
		{NewHeaderField("date", "Mon, 21 Oct 2013 20:13:22 GMT"), false},
	}

	context.Decode(buf)

	// test Header Table
	assert.Equal(t, context.HT, expectedHT)

	// test Emitted Set
	sort.Sort(context.ES)
	sort.Sort(expectedES)
	assert.Equal(t, context.ES, expectedES)

	// test Reference Set
	if expectedRS.Len() != context.RS.Len() {
		t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
	}

	for i, hf := range *context.RS {
		i := len(expectedRS) - i - 1
		if *expectedRS[i].HeaderField != *hf.HeaderField {
			t.Errorf("\n got %v\nwant %v", context.RS.Dump(), expectedRS.Dump())
		}
	}
}

func TestEncodeDecode(t *testing.T) {
	context := NewContext(DEFAULT_HEADER_TABLE_SIZE)

	hs := HeaderSet{
		*NewHeaderField(":status", "200"),
		*NewHeaderField("cache-control", "private"),
		*NewHeaderField("date", "Mon, 21 Oct 2013 20:13:22 GMT"),
		*NewHeaderField("location", "https://www.example.com"),
		*NewHeaderField("content-encoding", "gzip"),
		*NewHeaderField("set-cookie", "foo=ASDJKHQKBZXOQWEOPIUAXQWEOIU; max-age=3600; version=1"),
	}

	encoded := context.Encode(hs)
	context.Decode(encoded)
	assert.Equal(t, hs, *context.ES)
}

func TestEncodeDecodeWithHuffman(t *testing.T) {
	context := NewContext(DEFAULT_HEADER_TABLE_SIZE)
	hs := HeaderSet{
		*NewHeaderField(":method", "GET"),
		*NewHeaderField(":scheme", "http"),
		*NewHeaderField(":authority", "g-ecx.images-amazon.com"),
		*NewHeaderField(":path", "/images/G/01/gno/beacon/BeaconSprite-US-01._V401903535_.png"),
		*NewHeaderField("user-agent", "Mozilla/5.0 NewHeaderField(Macintosh; Intel Mac OS X 10.8; rv:16.0) Gecko/20100101 Firefox/16.0"),
		*NewHeaderField("accept", "image/png,image/*;q=0.8,*/*;q=0.5"),
		*NewHeaderField("accept-language", "en-US,en;q=0.5"),
		*NewHeaderField("accept-encoding", "gzip, deflate"),
		*NewHeaderField("connection", "keep-alive"),
		*NewHeaderField("referer", "http://www.amazon.com/"),
	}

	encoded := context.Encode(hs)
	if encoded[len(encoded)-1] == 255 {
		t.Error("8bit EOS on huffman encoded result is error")
	}
	context.Decode(encoded)
	assert.Equal(t, hs, *context.ES)
}
