package xdlms

import (
	"encoding/binary"
	"fmt"
)

// InitiateResponse represents an InitiateResponse APDU
// InitiateResponse ::= SEQUENCE {
//     negotiated-quality-of-service [0] IMPLICIT Integer8 OPTIONAL,
//     negotiated-dlms-version-number  Unsigned8,
//     negotiated-conformance  Conformance, -- Shall be encoded in BER
//     server-max-receive-pdu-size   Unsigned16,
//     vaa-name  ObjectName
// }
// When using LN referencing the value if vaa-name is always 0x0007
const InitiateResponseTag = 0x08

type InitiateResponse struct {
	*BaseXDlmsApdu
	NegotiatedConformance       *Conformance
	ServerMaxReceivePDUSize      uint16
	NegotiatedDlmsVersionNumber  uint8
	NegotiatedQualityOfService   uint8
}

// NewInitiateResponse creates a new InitiateResponse
func NewInitiateResponse(
	negotiatedConformance *Conformance,
	serverMaxReceivePDUSize uint16,
	negotiatedDlmsVersionNumber uint8,
	negotiatedQualityOfService uint8,
) *InitiateResponse {
	return &InitiateResponse{
		BaseXDlmsApdu: &BaseXDlmsApdu{
			Tag: InitiateResponseTag,
		},
		NegotiatedConformance:      negotiatedConformance,
		ServerMaxReceivePDUSize:   serverMaxReceivePDUSize,
		NegotiatedDlmsVersionNumber: negotiatedDlmsVersionNumber,
		NegotiatedQualityOfService: negotiatedQualityOfService,
	}
}

// FromBytes creates InitiateResponse from bytes
func (i *InitiateResponse) FromBytes(data []byte) (*InitiateResponse, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for InitiateResponse")
	}
	
	// Check vaa-name at the end
	if len(data) < 2 || data[len(data)-2] != 0x00 || data[len(data)-1] != 0x07 {
		return nil, fmt.Errorf("vaa-name in InitiateResponse is not 0x0007")
	}
	
	data = data[:len(data)-2]
	
	tag := data[0]
	if tag != InitiateResponseTag {
		return nil, fmt.Errorf("data is not an InitiateResponse APDU, got apdu tag %d", tag)
	}
	
	data = data[1:]
	
	// Parse negotiated_quality_of_service (optional)
	qualityOfService := uint8(0)
	if len(data) > 0 && data[0] == 0x01 {
		data = data[1:]
		if len(data) > 0 {
			qualityOfService = data[0]
			data = data[1:]
		}
	} else if len(data) > 0 && data[0] == 0x00 {
		data = data[1:]
	}
	
	// Parse negotiated_dlms_version_number
	if len(data) == 0 {
		return nil, fmt.Errorf("insufficient data for dlms version number")
	}
	dlmsVersion := data[0]
	data = data[1:]
	
	// Parse conformance (BER encoded)
	if len(data) < 3 {
		return nil, fmt.Errorf("insufficient data for conformance tag")
	}
	conformanceTagAndLength := data[:3]
	if string(conformanceTagAndLength) != "\x5f\x1f\x04" {
		return nil, fmt.Errorf("not correct conformance tag and length: %v", conformanceTagAndLength)
	}
	
	conformanceData := data[3:7] // 4 bytes: unused bits + 3 bytes for conformance
	conformance, err := (&Conformance{}).FromBytes(conformanceData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse conformance: %w", err)
	}
	data = data[7:]
	
	// Parse server_max_receive_pdu_size
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for max PDU size")
	}
	maxPDUSize := binary.BigEndian.Uint16(data[:2])
	
	return NewInitiateResponse(
		conformance,
		maxPDUSize,
		dlmsVersion,
		qualityOfService,
	), nil
}

// ToBytes converts InitiateResponse to bytes
func (i *InitiateResponse) ToBytes() ([]byte, error) {
	result := []byte{InitiateResponseTag}
	
	// Negotiated quality of service (optional, skip for now - use 0x00)
	result = append(result, 0x00)
	
	// Negotiated DLMS version number
	result = append(result, i.NegotiatedDlmsVersionNumber)
	
	// Conformance (BER encoded)
	result = append(result, 0x5f, 0x1f, 0x04)
	conformanceBytes := i.NegotiatedConformance.ToBytes()
	result = append(result, conformanceBytes...)
	
	// Server max receive PDU size
	pduSizeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(pduSizeBytes, i.ServerMaxReceivePDUSize)
	result = append(result, pduSizeBytes...)
	
	// VAA name (always 0x0007 for LN referencing)
	result = append(result, 0x00, 0x07)
	
	return result, nil
}

// GlobalCipherInitiateResponse represents a Global Cipher Initiate Response
const GlobalCipherInitiateResponseTag = 40

type GlobalCipherInitiateResponse struct {
	*BaseXDlmsApdu
	SecurityControl   interface{} // SecurityControlField - will be implemented when security module is ready
	InvocationCounter uint32
	CipheredText      []byte
}

// NewGlobalCipherInitiateResponse creates a new GlobalCipherInitiateResponse
func NewGlobalCipherInitiateResponse(
	securityControl interface{},
	invocationCounter uint32,
	cipheredText []byte,
) *GlobalCipherInitiateResponse {
	return &GlobalCipherInitiateResponse{
		BaseXDlmsApdu: &BaseXDlmsApdu{
			Tag: GlobalCipherInitiateResponseTag,
		},
		SecurityControl:   securityControl,
		InvocationCounter: invocationCounter,
		CipheredText:      cipheredText,
	}
}

// FromBytes creates GlobalCipherInitiateResponse from bytes
func (g *GlobalCipherInitiateResponse) FromBytes(data []byte) (*GlobalCipherInitiateResponse, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("insufficient data for GlobalCipherInitiateResponse")
	}
	
	tag := data[0]
	if tag != GlobalCipherInitiateResponseTag {
		return nil, fmt.Errorf("tag is not correct. Should be %d but got %d", GlobalCipherInitiateResponseTag, tag)
	}
	
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for length")
	}
	length := data[1]
	
	if len(data) < int(length)+2 {
		return nil, fmt.Errorf("insufficient data: need %d bytes, got %d", length+2, len(data))
	}
	
	octetStringData := data[2 : 2+length]
	if len(octetStringData) < 5 {
		return nil, fmt.Errorf("insufficient data in octet string")
	}
	
	// Security control (1 byte) - TODO: parse when security module is ready
	securityControl := octetStringData[0]
	
	// Invocation counter (4 bytes)
	invocationCounter := binary.BigEndian.Uint32(octetStringData[1:5])
	
	// Ciphered text (remaining bytes)
	cipheredText := make([]byte, len(octetStringData)-5)
	copy(cipheredText, octetStringData[5:])
	
	return NewGlobalCipherInitiateResponse(securityControl, invocationCounter, cipheredText), nil
}

// ToBytes converts GlobalCipherInitiateResponse to bytes
func (g *GlobalCipherInitiateResponse) ToBytes() ([]byte, error) {
	result := []byte{GlobalCipherInitiateResponseTag}
	
	octetStringData := make([]byte, 0)
	
	// Security control (1 byte) - TODO: convert when security module is ready
	if sc, ok := g.SecurityControl.(byte); ok {
		octetStringData = append(octetStringData, sc)
	} else {
		return nil, fmt.Errorf("security control must be byte for now")
	}
	
	// Invocation counter (4 bytes)
	icBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(icBytes, g.InvocationCounter)
	octetStringData = append(octetStringData, icBytes...)
	
	// Ciphered text
	octetStringData = append(octetStringData, g.CipheredText...)
	
	result = append(result, byte(len(octetStringData)))
	result = append(result, octetStringData...)
	
	return result, nil
}

