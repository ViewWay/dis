package xdlms

import (
	"encoding/binary"
	"fmt"
)

// InitiateRequest represents an InitiateRequest APDU
//
//	InitiateRequest ::= SEQUENCE {
//	    dedicated-key: OCTET STRING OPTIONAL
//	    response-allowed: BOOLEAN DEFAULT TRUE
//	    proposed-quality-of-service: IMPLICIT Integer8 OPTIONAL
//	    proposed-dlms-version-number: Integer8  # Always 6?
//	    proposed-conformance: Conformance
//	    client-max-receive-pdu-size: Unsigned16
//	}
const InitiateRequestTag = 0x01

type InitiateRequest struct {
	*BaseXDlmsApdu
	ProposedConformance       *Conformance
	ProposedQualityOfService  *int
	ClientMaxReceivePDUSize   uint16
	ProposedDlmsVersionNumber uint8
	ResponseAllowed           bool
	DedicatedKey              []byte
}

// NewInitiateRequest creates a new InitiateRequest
func NewInitiateRequest(
	proposedConformance *Conformance,
	clientMaxReceivePDUSize uint16,
	proposedDlmsVersionNumber uint8,
	responseAllowed bool,
	dedicatedKey []byte,
	proposedQualityOfService *int,
) *InitiateRequest {
	return &InitiateRequest{
		BaseXDlmsApdu: &BaseXDlmsApdu{
			Tag: InitiateRequestTag,
		},
		ProposedConformance:       proposedConformance,
		ClientMaxReceivePDUSize:   clientMaxReceivePDUSize,
		ProposedDlmsVersionNumber: proposedDlmsVersionNumber,
		ResponseAllowed:           responseAllowed,
		DedicatedKey:              dedicatedKey,
		ProposedQualityOfService:  proposedQualityOfService,
	}
}

// FromBytes creates InitiateRequest from bytes
func (i *InitiateRequest) FromBytes(data []byte) (*InitiateRequest, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("insufficient data for InitiateRequest")
	}

	apduTag := data[0]
	if apduTag != InitiateRequestTag {
		return nil, fmt.Errorf("data is not an InitiateRequest APDU, got apdu tag %d", apduTag)
	}

	data = data[1:]

	// Parse dedicated key (optional)
	var dedicatedKey []byte
	if len(data) > 0 && data[0] == 0x01 {
		data = data[1:]
		if len(data) == 0 {
			return nil, fmt.Errorf("insufficient data for dedicated key length")
		}
		keyLength := int(data[0])
		data = data[1:]
		if len(data) < keyLength {
			return nil, fmt.Errorf("insufficient data for dedicated key")
		}
		dedicatedKey = make([]byte, keyLength)
		copy(dedicatedKey, data[:keyLength])
		data = data[keyLength:]
	} else if len(data) > 0 && data[0] == 0x00 {
		data = data[1:]
	}

	// Parse response_allowed (default True, encoded as 0x00 for default)
	responseAllowed := true
	if len(data) > 0 && data[0] == 0x01 {
		data = data[1:]
		if len(data) > 0 {
			responseAllowed = data[0] != 0
			data = data[1:]
		}
	} else if len(data) > 0 && data[0] == 0x00 {
		data = data[1:]
	}

	// Parse proposed_quality_of_service (optional)
	var qualityOfService *int
	if len(data) > 0 && data[0] == 0x01 {
		data = data[1:]
		if len(data) > 0 {
			qos := int(int8(data[0]))
			qualityOfService = &qos
			data = data[1:]
		}
	} else if len(data) > 0 && data[0] == 0x00 {
		data = data[1:]
	}

	// Parse proposed_dlms_version_number
	if len(data) == 0 {
		return nil, fmt.Errorf("insufficient data for dlms version number")
	}
	dlmsVersion := data[0]
	data = data[1:]

	// Parse conformance (BER encoded)
	if len(data) < 3 {
		return nil, fmt.Errorf("insufficient data for conformance tag")
	}
	conformanceTag := data[:2]
	if string(conformanceTag) != "\x5f\x1f" {
		return nil, fmt.Errorf("didn't receive conformance tag correctly, got %v", conformanceTag)
	}
	conformanceLength := data[2]
	if len(data) < int(conformanceLength)+3 {
		return nil, fmt.Errorf("insufficient data for conformance")
	}
	// conformanceLength byte is at data[2], actual conformance data starts at data[3]
	conformanceData := data[3 : 3+conformanceLength]
	conformance, err := (&Conformance{}).FromBytes(conformanceData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse conformance: %w", err)
	}
	data = data[3+conformanceLength:]

	// Parse client_max_receive_pdu_size
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for max PDU size")
	}
	maxPDUSize := binary.BigEndian.Uint16(data[:2])

	return NewInitiateRequest(
		conformance,
		maxPDUSize,
		dlmsVersion,
		responseAllowed,
		dedicatedKey,
		qualityOfService,
	), nil
}

// ToBytes converts InitiateRequest to bytes
func (i *InitiateRequest) ToBytes() ([]byte, error) {
	result := []byte{InitiateRequestTag}

	// Dedicated key (optional)
	if len(i.DedicatedKey) > 0 {
		result = append(result, 0x01)
		result = append(result, byte(len(i.DedicatedKey)))
		result = append(result, i.DedicatedKey...)
	} else {
		result = append(result, 0x00)
	}

	// Response allowed (default True, encoded as 0x00)
	result = append(result, 0x00)

	// Proposed quality of service (optional, skip for now)

	// Proposed DLMS version number
	result = append(result, i.ProposedDlmsVersionNumber)

	// Conformance (BER encoded)
	result = append(result, 0x5f, 0x1f, 0x04)
	conformanceBytes := i.ProposedConformance.ToBytes()
	result = append(result, conformanceBytes...)

	// Client max receive PDU size
	pduSizeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(pduSizeBytes, i.ClientMaxReceivePDUSize)
	result = append(result, pduSizeBytes...)

	return result, nil
}

// GlobalCipherInitiateRequest represents a Global Cipher Initiate Request
const GlobalCipherInitiateRequestTag = 33

type GlobalCipherInitiateRequest struct {
	*BaseXDlmsApdu
	SecurityControl   interface{} // SecurityControlField - will be implemented when security module is ready
	InvocationCounter uint32
	CipheredText      []byte
}

// NewGlobalCipherInitiateRequest creates a new GlobalCipherInitiateRequest
func NewGlobalCipherInitiateRequest(
	securityControl interface{},
	invocationCounter uint32,
	cipheredText []byte,
) *GlobalCipherInitiateRequest {
	return &GlobalCipherInitiateRequest{
		BaseXDlmsApdu: &BaseXDlmsApdu{
			Tag: GlobalCipherInitiateRequestTag,
		},
		SecurityControl:   securityControl,
		InvocationCounter: invocationCounter,
		CipheredText:      cipheredText,
	}
}

// FromBytes creates GlobalCipherInitiateRequest from bytes
func (g *GlobalCipherInitiateRequest) FromBytes(data []byte) (*GlobalCipherInitiateRequest, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("insufficient data for GlobalCipherInitiateRequest")
	}

	tag := data[0]
	if tag != GlobalCipherInitiateRequestTag {
		return nil, fmt.Errorf("tag is not correct. Should be %d but got %d", GlobalCipherInitiateRequestTag, tag)
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

	return NewGlobalCipherInitiateRequest(securityControl, invocationCounter, cipheredText), nil
}

// ToBytes converts GlobalCipherInitiateRequest to bytes
func (g *GlobalCipherInitiateRequest) ToBytes() ([]byte, error) {
	result := []byte{GlobalCipherInitiateRequestTag}

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
