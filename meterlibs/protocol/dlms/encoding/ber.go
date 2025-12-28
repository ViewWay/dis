package encoding

import (
	"fmt"
)

// BER provides Basic Encoding Rules encoding/decoding
// BER encoding consists of a TAG ID, Length and data
type BER struct{}

// Encode encodes data using BER encoding
// tag can be either an int (single byte) or bytes
// data must be bytes or bytearray
func (b *BER) Encode(tag interface{}, data []byte) ([]byte, error) {
	var tagBytes []byte

	switch t := tag.(type) {
	case int:
		tagBytes = []byte{byte(t)}
	case uint8:
		tagBytes = []byte{byte(t)}
	case []byte:
		tagBytes = t
	default:
		return nil, fmt.Errorf("BER encoding requires int, uint8 or []byte for tag, got %T", tag)
	}

	if data == nil {
		return []byte{}, nil
	}

	if len(data) == 0 {
		return []byte{}, nil
	}

	length := byte(len(data))
	result := make([]byte, 0, len(tagBytes)+1+len(data))
	result = append(result, tagBytes...)
	result = append(result, length)
	result = append(result, data...)

	return result, nil
}

// Decode decodes BER encoded data
// Returns tag, length, and data
func (b *BER) Decode(data []byte, tagLength int) ([]byte, uint8, []byte, error) {
	if len(data) < tagLength+1 {
		return nil, 0, nil, fmt.Errorf("insufficient data for BER decoding")
	}

	input := make([]byte, len(data))
	copy(input, data)

	tag := make([]byte, tagLength)
	for i := 0; i < tagLength; i++ {
		tag[i] = input[0]
		input = input[1:]
	}

	length := input[0]
	input = input[1:]

	if len(input) != int(length) {
		return nil, 0, nil, fmt.Errorf("BER-decoding failed. Length byte %d does not correspond to length of data %d", length, len(input))
	}

	return tag, length, input, nil
}

// NewBER creates a new BER encoder/decoder
func NewBER() *BER {
	return &BER{}
}

