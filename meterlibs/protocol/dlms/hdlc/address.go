package hdlc

import (
	"fmt"
)

// AddressType represents the type of HDLC address
type AddressType string

const (
	AddressTypeClient AddressType = "client"
	AddressTypeServer AddressType = "server"
)

// HdlcAddress represents an HDLC address
// A client address shall always be expressed on one byte.
// To enable addressing more than one logical device within a single physical device
// and to support the multi-drop configuration the server address may be divided in
// two parts:
// The logical address to address a logical device (separate addressable entity
// within a physical device) makes up the upper HDLC address
// The logical address must always be present.
// The physical address is used to address a physical device (a physical device on
// a multi-drop)
// The physical address can be omitted if not used.
type HdlcAddress struct {
	LogicalAddress    int
	PhysicalAddress   *int // nil if not used
	AddressType       AddressType
	ExtendedAddressing bool
}

// NewHdlcAddress creates a new HDLC address
func NewHdlcAddress(logicalAddress int, physicalAddress *int, addressType AddressType, extendedAddressing bool) (*HdlcAddress, error) {
	if err := validateHdlcAddress(logicalAddress); err != nil {
		return nil, fmt.Errorf("invalid logical address: %w", err)
	}
	if physicalAddress != nil {
		if err := validateHdlcAddress(*physicalAddress); err != nil {
			return nil, fmt.Errorf("invalid physical address: %w", err)
		}
	}
	if err := validateHdlcAddressType(addressType); err != nil {
		return nil, fmt.Errorf("invalid address type: %w", err)
	}

	return &HdlcAddress{
		LogicalAddress:    logicalAddress,
		PhysicalAddress:   physicalAddress,
		AddressType:       addressType,
		ExtendedAddressing: extendedAddressing,
	}, nil
}

// Length returns the number of bytes the address makes up
func (a *HdlcAddress) Length() int {
	bytes := a.ToBytes()
	return len(bytes)
}

// ToBytes converts the HDLC address to bytes
func (a *HdlcAddress) ToBytes() []byte {
	var out []byte

	if a.AddressType == AddressTypeClient {
		// shift left 1 bit and set the lsb to mark end of address
		out = append(out, byte((a.LogicalAddress<<1)|0b00000001))
	} else {
		// server address type
		logicalHigher, logicalLower := a.splitAddress(a.LogicalAddress)

		if a.PhysicalAddress != nil {
			physicalHigher, physicalLower := a.splitAddress(*a.PhysicalAddress)
			// mark physical lower as end
			physicalLower = physicalLower | 0b00000001
			out = append(out, logicalHigher, logicalLower, physicalHigher, physicalLower)
		} else {
			// no physical address so mark the logical as end
			logicalLower = logicalLower | 0b00000001
			out = append(out, logicalHigher, logicalLower)
		}
	}

	var outBytes []byte
	for _, addr := range out {
		if addr == 0 {
			if a.ExtendedAddressing {
				outBytes = append(outBytes, addr)
			}
			// else skip
		} else {
			outBytes = append(outBytes, addr)
		}
	}

	return outBytes
}

// splitAddress splits an address into higher and lower parts
func (a *HdlcAddress) splitAddress(address int) (byte, byte) {
	var higher byte
	var lower byte

	if address > 0b01111111 {
		lower = byte((address & 0b0000000001111111) << 1)
		higher = byte((address & 0b0011111110000000) >> 6)
	} else {
		lower = byte(address << 1)
		higher = 0
	}

	return higher, lower
}

// DestinationFromBytes creates an HDLC address from frame bytes (destination address)
func DestinationFromBytes(frameBytes []byte, addressType AddressType) (*HdlcAddress, error) {
	destData, _, err := FindAddressInFrameBytes(frameBytes)
	if err != nil {
		return nil, err
	}

	destLogical, destPhysical, _ := destData
	var physicalAddr *int
	if destPhysical != nil {
		physicalAddr = destPhysical
	}

	return NewHdlcAddress(destLogical, physicalAddr, addressType, false)
}

// SourceFromBytes creates an HDLC address from frame bytes (source address)
func SourceFromBytes(frameBytes []byte, addressType AddressType) (*HdlcAddress, error) {
	_, sourceData, err := FindAddressInFrameBytes(frameBytes)
	if err != nil {
		return nil, err
	}

	sourceLogical, sourcePhysical, sourceLength := sourceData
	extendedAddress := sourceLength == 4

	var physicalAddr *int
	if sourcePhysical != nil {
		physicalAddr = sourcePhysical
	}

	return NewHdlcAddress(sourceLogical, physicalAddr, addressType, extendedAddress)
}

// ExtractAddressBytes extracts address bytes from input data
func ExtractAddressBytes(inData []byte) ([]byte, []byte, error) {
	var address []byte
	data := make([]byte, len(inData))
	copy(data, inData)

	foundWholeAddress := false
	for !foundWholeAddress {
		if len(data) == 0 {
			return nil, nil, fmt.Errorf("insufficient data to extract address")
		}

		byteVal := data[0]
		data = data[1:]

		address = append(address, byteVal)
		if byteVal&0b00000001 != 0 {
			foundWholeAddress = true
		}

		if len(address) > 4 {
			return nil, nil, fmt.Errorf("recovered an HDLC address of length longer than 4 bytes")
		}
	}

	return address, data, nil
}

// AddressData represents address information
type AddressData struct {
	Logical  int
	Physical *int
	Length   int
}

// FindAddressInFrameBytes finds destination and source addresses in HDLC frame bytes
// Address can be 1, 2 or 4 bytes long. The end byte is indicated by the
// last byte LSB being 1
// The first address is the destination address and the second is the source address.
func FindAddressInFrameBytes(hdlcFrameBytes []byte) (AddressData, AddressData, error) {
	if len(hdlcFrameBytes) < 4 {
		return AddressData{}, AddressData{}, fmt.Errorf("frame too short")
	}

	// Find destination address
	destinationLength := 1
	destinationPositions := [][]int{{3, 1}, {4, 2}, {6, 4}}

	for _, pos := range destinationPositions {
		posIdx, length := pos[0], pos[1]
		if posIdx >= len(hdlcFrameBytes) {
			continue
		}
		endByte := hdlcFrameBytes[posIdx]
		if endByte&0b00000001 != 0 {
			destinationLength = length
			break
		}
	}

	var destinationLogical int
	var destinationPhysical *int

	switch destinationLength {
	case 1:
		if len(hdlcFrameBytes) < 4 {
			return AddressData{}, AddressData{}, fmt.Errorf("frame too short for destination address")
		}
		destinationLogical = int(hdlcFrameBytes[3] >> 1)
		destinationPhysical = nil

	case 2:
		if len(hdlcFrameBytes) < 5 {
			return AddressData{}, AddressData{}, fmt.Errorf("frame too short for 2-byte destination address")
		}
		destinationLogical = int(hdlcFrameBytes[3] >> 1)
		physical := int(hdlcFrameBytes[4] >> 1)
		destinationPhysical = &physical

	case 4:
		if len(hdlcFrameBytes) < 7 {
			return AddressData{}, AddressData{}, fmt.Errorf("frame too short for 4-byte destination address")
		}
		destBytes := hdlcFrameBytes[3:7]
		destinationLogical = parseTwoByteAddress(destBytes[:2])
		physical := parseTwoByteAddress(destBytes[2:])
		destinationPhysical = &physical
	}

	// Find source address
	sourceLength := 1
	sourceStartPos := 3 + destinationLength
	sourcePositions := [][]int{
		{sourceStartPos, 1},
		{sourceStartPos + 1, 2},
		{sourceStartPos + 3, 4},
	}

	for _, pos := range sourcePositions {
		posIdx, length := pos[0], pos[1]
		if posIdx >= len(hdlcFrameBytes) {
			continue
		}
		endByte := hdlcFrameBytes[posIdx]
		if endByte&0b00000001 != 0 {
			sourceLength = length
			break
		}
	}

	var sourceLogical int
	var sourcePhysical *int

	switch sourceLength {
	case 1:
		if len(hdlcFrameBytes) < sourceStartPos+1 {
			return AddressData{}, AddressData{}, fmt.Errorf("frame too short for source address")
		}
		sourceLogical = int(hdlcFrameBytes[sourceStartPos] >> 1)
		sourcePhysical = nil

	case 2:
		if len(hdlcFrameBytes) < sourceStartPos+2 {
			return AddressData{}, AddressData{}, fmt.Errorf("frame too short for 2-byte source address")
		}
		sourceLogical = int(hdlcFrameBytes[sourceStartPos] >> 1)
		physical := int(hdlcFrameBytes[sourceStartPos+1] >> 1)
		sourcePhysical = &physical

	case 4:
		if len(hdlcFrameBytes) < sourceStartPos+4 {
			return AddressData{}, AddressData{}, fmt.Errorf("frame too short for 4-byte source address")
		}
		sourceBytes := hdlcFrameBytes[sourceStartPos : sourceStartPos+4]
		sourceLogical = parseTwoByteAddress(sourceBytes[:2])
		physical := parseTwoByteAddress(sourceBytes[2:])
		sourcePhysical = &physical
	}

	destData := AddressData{
		Logical:  destinationLogical,
		Physical: destinationPhysical,
		Length:   destinationLength,
	}

	sourceData := AddressData{
		Logical:  sourceLogical,
		Physical: sourcePhysical,
		Length:   sourceLength,
	}

	return destData, sourceData, nil
}

// parseTwoByteAddress parses a two-byte address
func parseTwoByteAddress(addressBytes []byte) int {
	if len(addressBytes) != 2 {
		panic("can only parse 2 bytes for address")
	}
	upper := addressBytes[0] >> 1
	lower := addressBytes[1] >> 1
	return int(lower) + (int(upper) << 7)
}

// validateHdlcAddress validates an HDLC address value
func validateHdlcAddress(value int) error {
	if value < 0 || value > 127 {
		return fmt.Errorf("HDLC address must be between 0 and 127, got %d", value)
	}
	return nil
}

// validateHdlcAddressType validates an address type
func validateHdlcAddressType(addressType AddressType) error {
	if addressType != AddressTypeClient && addressType != AddressTypeServer {
		return fmt.Errorf("invalid address type: %s", addressType)
	}
	return nil
}

