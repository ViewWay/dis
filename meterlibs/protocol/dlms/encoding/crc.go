package encoding

// CRCCCITT provides CRC CCITT - HDLC Style 16-bit
// In accordance with ANSI C12.18(2006)
// Using 0xFFFF as initial value
// Running over serial so all message bytes need to be reversed
// before calculation (because least significant bit is sent first)
// resulting crc bytes needs to be reversed to become in correct order
// The reversed crc is then XOR:ed with 0xFFFF
type CRCCCITT struct {
	crcTable        [256]uint16
	startingValue   uint16
	crcCCITTConstant uint16
	initialized     bool
}

// NewCRCCCITT creates a new CRCCCITT calculator
func NewCRCCCITT() *CRCCCITT {
	crc := &CRCCCITT{
		startingValue:    0xFFFF,
		crcCCITTConstant: 0x1021,
		initialized:      false,
	}
	crc.initCRCTable()
	return crc
}

// initCRCTable initializes the pre-calculated CRC table
func (c *CRCCCITT) initCRCTable() {
	if c.initialized {
		return
	}

	for i := 0; i < 256; i++ {
		crc := uint16(0)
		cVal := uint16(i) << 8

		for j := 0; j < 8; j++ {
			if (crc^cVal)&0x8000 != 0 {
				crc = (crc << 1) ^ c.crcCCITTConstant
			} else {
				crc = crc << 1
			}
			cVal = cVal << 1
		}

		c.crcTable[i] = crc
	}
	c.initialized = true
}

// CalculateFor calculates CRC for input data
// lsbFirst indicates if the Least significant byte should be returned first (little endian)
func (c *CRCCCITT) CalculateFor(inputData []byte, lsbFirst bool) []byte {
	// need to reverse bits in bytes
	reversedData := reverseByteMessage(inputData)

	reversedCRC := c.calculate(reversedData)
	lsbRev := reversedCRC & 0x00FF
	lsb := reverseByte(byte(lsbRev))
	lsb ^= 0xFF
	lsbByte := byte(lsb)

	msbRev := (reversedCRC & 0xFF00) >> 8
	msb := reverseByte(byte(msbRev))
	msb ^= 0xFF
	msbByte := byte(msb)

	if lsbFirst {
		return []byte{lsbByte, msbByte}
	}
	return []byte{msbByte, lsbByte}
}

// calculate performs the CRC calculation
func (c *CRCCCITT) calculate(inputData []byte) uint16 {
	crcValue := c.startingValue

	for _, char := range inputData {
		tmp := ((crcValue >> 8) & 0xFF) ^ uint16(char)
		crcShifted := (crcValue << 8) & 0xFF00
		crcValue = crcShifted ^ c.crcTable[tmp]
	}

	return crcValue
}

// reverseByte reverses the bits in a byte
func reverseByte(byteToReverse byte) byte {
	andValue := byte(1)
	reversedByte := byte(0)

	for i := 0; i < 8; i++ {
		reversedByte += ((byteToReverse & andValue) >> i) * (1 << (7 - i))
		andValue <<= 1
	}

	return reversedByte
}

// reverseByteMessage reverses bits in all bytes of a message
func reverseByteMessage(msg []byte) []byte {
	reversedMsg := make([]byte, len(msg))
	for i, char := range msg {
		reversedMsg[i] = reverseByte(char)
	}
	return reversedMsg
}

