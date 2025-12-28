package xdlms

// XDlmsApdu is the interface for all xDLMS APDUs
type XDlmsApdu interface {
	FromBytes(data []byte) (XDlmsApdu, error)
	ToBytes() ([]byte, error)
	GetTag() uint8
}

// BaseXDlmsApdu is the base struct for xDLMS APDUs
type BaseXDlmsApdu struct {
	Tag uint8
}

// GetTag returns the tag
func (b *BaseXDlmsApdu) GetTag() uint8 {
	return b.Tag
}

