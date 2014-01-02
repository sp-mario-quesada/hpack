package hpack

// A name-value pair.
// Both name and value are sequences of octets.
type HeaderField struct {
	Name  string
	Value string
}

// Add prefix if name is Must Header
func NewHeaderField(name, value string) *HeaderField {
	return &HeaderField{name, value}
}

// The size of an entry is the sum of its name's length in octets
// of its value's length in octets and of 32 octets.
func (h *HeaderField) Size() int {
	return len(h.Name) + len(h.Value) + 32
}
