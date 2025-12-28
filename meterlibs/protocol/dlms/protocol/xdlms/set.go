package xdlms

import (
	"fmt"

	"github.com/yimiliya/idis/meterlibs/protocol/dlms/cosem"
	"github.com/yimiliya/idis/meterlibs/protocol/dlms/enumerations"
)

// SetRequestNormal represents a Set request normal
// Set-Request-Normal ::= SEQUENCE
// {
//     invoke-id-and-priority          Invoke-Id-And-Priority,
//     cosem-attribute-descriptor      Cosem-Attribute-Descriptor,
//     access-selection                Selective-Access-Descriptor OPTIONAL,
//     value                           Data
// }
const SetRequestTag = 193

type SetRequestNormal struct {
	*BaseXDlmsApdu
	CosemAttribute      *cosem.CosemAttribute
	Data                []byte
	AccessSelection     interface{} // Optional selective access
	InvokeIdAndPriority *InvokeIdAndPriority
}

// NewSetRequestNormal creates a new SetRequestNormal
func NewSetRequestNormal(
	cosemAttribute *cosem.CosemAttribute,
	data []byte,
	accessSelection interface{},
	invokeIdAndPriority *InvokeIdAndPriority,
) *SetRequestNormal {
	return &SetRequestNormal{
		BaseXDlmsApdu: &BaseXDlmsApdu{
			Tag: SetRequestTag,
		},
		CosemAttribute:      cosemAttribute,
		Data:                data,
		AccessSelection:     accessSelection,
		InvokeIdAndPriority: invokeIdAndPriority,
	}
}

// FromBytes creates SetRequestNormal from bytes
func (s *SetRequestNormal) FromBytes(data []byte) (*SetRequestNormal, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for SetRequest")
	}
	
	tag := data[0]
	if tag != SetRequestTag {
		return nil, fmt.Errorf("tag for SetRequest is not correct. Got %d, should be %d", tag, SetRequestTag)
	}
	
	typeChoice := enumerations.SetRequestType(data[1])
	if typeChoice != enumerations.SetRequestTypeNormal {
		return nil, fmt.Errorf("the type of the SetRequest is not for a SetRequestNormal")
	}
	
	data = data[2:]
	
	// Parse invoke_id_and_priority
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for invoke_id_and_priority")
	}
	invokeIdAndPriority, err := (&InvokeIdAndPriority{}).FromBytes(data[:1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse invoke_id_and_priority: %w", err)
	}
	data = data[1:]
	
	// Parse cosem_attribute (9 bytes)
	if len(data) < 9 {
		return nil, fmt.Errorf("insufficient data for cosem_attribute")
	}
	cosemAttribute, err := (&cosem.CosemAttribute{}).FromBytes(data[:9])
	if err != nil {
		return nil, fmt.Errorf("failed to parse cosem_attribute: %w", err)
	}
	data = data[9:]
	
	// Parse access_selection (optional)
	var accessSelection interface{}
	if len(data) > 0 {
		hasAccessSelection := data[0] != 0
		data = data[1:]
		if hasAccessSelection {
			// Parse access descriptor using factory
			// Note: FromBytes will validate data length internally, but we check here for clarity
			if len(data) == 0 {
				return nil, fmt.Errorf("insufficient data for access descriptor: expected access descriptor data, got empty buffer")
			}
			factory := cosem.NewAccessDescriptorFactory()
			parsedAccess, bytesConsumed, err := factory.FromBytes(data)
			if err != nil {
				return nil, fmt.Errorf("failed to parse access selection: %w", err)
			}
			accessSelection = parsedAccess
			
			// Validate that we have enough data before advancing pointer
			// This is a defensive check - FromBytes should have already validated this
			if bytesConsumed < 0 {
				return nil, fmt.Errorf("invalid bytes consumed: %d", bytesConsumed)
			}
			if len(data) < bytesConsumed {
				return nil, fmt.Errorf("insufficient data for access descriptor: need %d bytes, got %d", bytesConsumed, len(data))
			}
			// Advance data pointer by the number of bytes consumed by the access descriptor
			data = data[bytesConsumed:]
		}
	}
	
	// Remaining data is the value
	valueData := make([]byte, len(data))
	copy(valueData, data)
	
	return NewSetRequestNormal(cosemAttribute, valueData, accessSelection, invokeIdAndPriority), nil
}

// ToBytes converts SetRequestNormal to bytes
func (s *SetRequestNormal) ToBytes() ([]byte, error) {
	result := []byte{SetRequestTag}
	result = append(result, byte(enumerations.SetRequestTypeNormal))
	
	invokeBytes := s.InvokeIdAndPriority.ToBytes()
	result = append(result, invokeBytes...)
	
	cosemBytes := s.CosemAttribute.ToBytes()
	result = append(result, cosemBytes...)
	
	if s.AccessSelection != nil {
		result = append(result, 0x01)
		// Serialize access selection based on its type
		switch accessSel := s.AccessSelection.(type) {
		case *cosem.RangeDescriptor:
			result = append(result, accessSel.ToBytes()...)
		case *cosem.EntryDescriptor:
			result = append(result, accessSel.ToBytes()...)
		default:
			return nil, fmt.Errorf("unsupported access selection type: %T", s.AccessSelection)
		}
	} else {
		result = append(result, 0x00)
	}
	
	result = append(result, s.Data...)
	
	return result, nil
}

// SetResponseNormal represents a Set response normal
const SetResponseTag = 197

type SetResponseNormal struct {
	*BaseXDlmsApdu
	InvokeIdAndPriority *InvokeIdAndPriority
	Result              enumerations.DataAccessResult
}

// NewSetResponseNormal creates a new SetResponseNormal
func NewSetResponseNormal(
	invokeIdAndPriority *InvokeIdAndPriority,
	result enumerations.DataAccessResult,
) *SetResponseNormal {
	return &SetResponseNormal{
		BaseXDlmsApdu: &BaseXDlmsApdu{
			Tag: SetResponseTag,
		},
		InvokeIdAndPriority: invokeIdAndPriority,
		Result:              result,
	}
}

// FromBytes creates SetResponseNormal from bytes
func (s *SetResponseNormal) FromBytes(data []byte) (*SetResponseNormal, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for SetResponse")
	}
	
	tag := data[0]
	if tag != SetResponseTag {
		return nil, fmt.Errorf("tag for SetResponse is not correct. Got %d, should be %d", tag, SetResponseTag)
	}
	
	typeChoice := enumerations.SetResponseType(data[1])
	if typeChoice != enumerations.SetResponseTypeNormal {
		return nil, fmt.Errorf("the type of the SetResponse is not for a SetResponseNormal")
	}
	
	data = data[2:]
	
	// Parse invoke_id_and_priority
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for invoke_id_and_priority")
	}
	invokeIdAndPriority, err := (&InvokeIdAndPriority{}).FromBytes(data[:1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse invoke_id_and_priority: %w", err)
	}
	data = data[1:]
	
	// Parse result
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for result")
	}
	result := enumerations.DataAccessResult(data[0])
	
	return NewSetResponseNormal(invokeIdAndPriority, result), nil
}

// ToBytes converts SetResponseNormal to bytes
func (s *SetResponseNormal) ToBytes() ([]byte, error) {
	result := []byte{SetResponseTag}
	result = append(result, byte(enumerations.SetResponseTypeNormal))
	
	invokeBytes := s.InvokeIdAndPriority.ToBytes()
	result = append(result, invokeBytes...)
	
	result = append(result, byte(s.Result))
	
	return result, nil
}
