package hpac

import (
	"bytes"
	"encoding/binary"
	"log"
)

func DecodeHeader(buf *bytes.Buffer) Frame {
	log.SetFlags(log.Lshortfile)
	var types uint8
	if err := binary.Read(buf, binary.BigEndian, &types); err != nil {
		log.Println("binary.Read failed:", err)
	}
	if types > 0x80 {

		frame := &IndexedHeader{}
		frame.Flag1 = 1
		frame.Index = types & 0x7F

		log.Println("Indexed Header Representation")
		return frame

	} else if types == 0 {

		frame := &NewNameWithSubstitutionIndexing{}
		frame.Flag1 = 0
		frame.Flag2 = 0
		frame.Index = 0
		frame.NameLength = DecodePrefixedInteger(buf, 8)
		frame.NameString = DecodeString(buf, frame.NameLength)
		frame.SubstitutedIndex = DecodePrefixedInteger(buf, 8)
		frame.ValueLength = DecodePrefixedInteger(buf, 8)
		frame.ValueString = DecodeString(buf, frame.ValueLength)

		log.Println("Literal Header with Substitution Indexing - New Name")
		return frame

	} else if types == 0x40 {

		frame := &NewNameWithIncrementalIndexing{}
		frame.Flag1 = 0
		frame.Flag2 = 1
		frame.Flag3 = 0
		frame.Index = 0
		frame.NameLength = DecodePrefixedInteger(buf, 8)
		frame.NameString = DecodeString(buf, frame.NameLength)
		frame.ValueLength = DecodePrefixedInteger(buf, 8)
		frame.ValueString = DecodeString(buf, frame.ValueLength)

		log.Println("Literal Header with Incremental Indexing - New Name")
		return frame

	} else if types == 0x60 {

		var frame = &NewNameWithoutIndexing{}
		frame.Flag1 = 0
		frame.Flag2 = 1
		frame.Flag3 = 1
		frame.Index = 0
		frame.NameLength = DecodePrefixedInteger(buf, 8)
		frame.NameString = DecodeString(buf, frame.NameLength)
		frame.ValueLength = DecodePrefixedInteger(buf, 8)
		frame.ValueString = DecodeString(buf, frame.ValueLength)

		log.Println("Literal Header without Indexing - New Name")
		return frame

	} else if types>>5 == 0x2 {

		// unread first byte for parse frame
		buf.UnreadByte()

		var frame = &IndexedNameWithIncrementalIndexing{}
		frame.Flag1 = 0
		frame.Flag2 = 1
		frame.Flag3 = 0
		// 0 describes "not in the header table", but index of Header Table start with 0
		// so Index is represented as +1 integer
		frame.Index = DecodePrefixedInteger(buf, 5) - 1
		frame.ValueLength = DecodePrefixedInteger(buf, 8)
		frame.ValueString = DecodeString(buf, frame.ValueLength)

		log.Println("Literal Header with Incremental Indexing - Indexed Name")
		return frame

	} else if types&0x60 == 0x60 {

		buf.UnreadByte()

		var frame = &IndexedNameWithoutIndexing{}
		frame.Flag1 = 0
		frame.Flag2 = 1
		frame.Flag3 = 1
		frame.Index = DecodePrefixedInteger(buf, 5) - 1
		frame.ValueLength = DecodePrefixedInteger(buf, 8)
		frame.ValueString = DecodeString(buf, frame.ValueLength)

		log.Println("Literal Header without Indexing - Indexed Name")
		return frame

	} else {

		// unread first byte for parse frame
		buf.UnreadByte()

		var frame = &IndexedNameWithSubstitutionIndexing{}
		frame.Flag1 = 0
		frame.Flag2 = 0
		frame.Index = DecodePrefixedInteger(buf, 6) - 1
		frame.SubstitutedIndex = DecodePrefixedInteger(buf, 8)
		frame.ValueLength = DecodePrefixedInteger(buf, 8)
		frame.ValueString = DecodeString(buf, frame.ValueLength)

		log.Println("Literal Header with Substitution Indexing - Indexed Name")
		return frame

	}
	return nil
}

func DecodeString(buf *bytes.Buffer, n uint32) string {
	valueBytes := make([]byte, n)
	binary.Read(buf, binary.BigEndian, &valueBytes) // err
	return string(valueBytes)
}
