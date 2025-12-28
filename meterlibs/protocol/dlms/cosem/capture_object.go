package cosem

// CaptureObject represents a value that is supposed to be saved in a Profile Generic.
// A data_index of 0 means the whole attribute is referenced. Otherwise it points to a
// specific element of the attribute. For example an entry in a buffer.
// DataIndex is uint16 because DLMS UnsignedLong (tag 0x12) is 2 bytes, matching Python's UnsignedLongData.
type CaptureObject struct {
	CosemAttribute *CosemAttribute
	DataIndex      uint16
}

// NewCaptureObject creates a new CaptureObject
func NewCaptureObject(cosemAttribute *CosemAttribute, dataIndex uint16) *CaptureObject {
	return &CaptureObject{
		CosemAttribute: cosemAttribute,
		DataIndex:      dataIndex,
	}
}

// ToBytes converts CaptureObject to bytes
// Returns a structure of 4 elements:
// - interface (UnsignedLong)
// - instance (OctetString - OBIS)
// - attribute (Integer)
// - data_index (UnsignedLong)
func (c *CaptureObject) ToBytes() []byte {
	// TODO: This will need to be implemented when dlmsdata package is ready
	// For now, this is a placeholder structure
	// Structure tag (0x02) + length (0x04 for 4 elements)
	result := []byte{0x02, 0x04}
	
	// Interface (UnsignedLong - 2 bytes)
	interfaceBytes := make([]byte, 2)
	interfaceBytes[0] = byte(uint16(c.CosemAttribute.Interface) >> 8)
	interfaceBytes[1] = byte(uint16(c.CosemAttribute.Interface))
	result = append(result, 0x11, 0x02) // Tag UnsignedLong + length
	result = append(result, interfaceBytes...)
	
	// Instance (OctetString - 6 bytes OBIS)
	obisBytes := c.CosemAttribute.Instance.ToBytes()
	result = append(result, 0x09, 0x06) // Tag OctetString + length
	result = append(result, obisBytes...)
	
	// Attribute (Integer - 1 byte)
	result = append(result, 0x0F, 0x01) // Tag Integer + length
	result = append(result, c.CosemAttribute.Attribute)
	
	// DataIndex (UnsignedLong - 2 bytes, tag 0x12)
	// DataIndex is uint16, so no range check needed - type system guarantees it fits in 2 bytes
	dataIndexBytes := make([]byte, 2)
	dataIndexBytes[0] = byte(c.DataIndex >> 8)
	dataIndexBytes[1] = byte(c.DataIndex & 0xFF)
	result = append(result, 0x12, 0x02) // Tag UnsignedLong (0x12) + length (0x02)
	result = append(result, dataIndexBytes...)
	
	return result
}

