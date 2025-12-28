package utils

import (
	"fmt"
)

// ParseAsDlmsData parses data as DLMS data structure
// This is a placeholder - will be implemented when dlmsdata package is ready
func ParseAsDlmsData(data []byte) (interface{}, error) {
	// TODO: Implement using AXdrDecoder when available
	return nil, fmt.Errorf("ParseAsDlmsData not yet implemented")
}

// ParseDlmsObject parses DLMS object attributes that contain data structures
// The items are not self descriptive and we cannot use the normal parser
// since we don't know what items to parse as the data tag is not included.
// But it seems they are always integers so we can parse them as a list of integers.
func ParseDlmsObject(sourceBytes []byte) ([]int, error) {
	if len(sourceBytes) == 0 {
		return nil, fmt.Errorf("empty source bytes")
	}

	data := make([]byte, len(sourceBytes))
	copy(data, sourceBytes)

	tag := data[0]
	data = data[1:]

	// TagArray = 1, TagStructure = 2
	allowedTags := []byte{1, 2}

	allowed := false
	for _, allowedTag := range allowedTags {
		if tag == allowedTag {
			allowed = true
			break
		}
	}

	if !allowed {
		return nil, fmt.Errorf("cannot use dlms object parse with tag %d", tag)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("insufficient data for length")
	}

	length := int(data[0])
	data = data[1:]

	if len(data) < length {
		return nil, fmt.Errorf("insufficient data: need %d bytes, got %d", length, len(data))
	}

	values := make([]int, length)
	for i := 0; i < length; i++ {
		values[i] = int(data[i])
	}

	return values, nil
}

