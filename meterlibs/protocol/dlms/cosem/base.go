package cosem

import (
	"encoding/binary"
	"fmt"

	"github.com/yimiliya/idis/meterlibs/protocol/dlms/enumerations"
)

// CosemAttribute represents a COSEM attribute descriptor
type CosemAttribute struct {
	Interface enumerations.CosemInterface
	Instance  *Obis
	Attribute uint8
}

const CosemAttributeLength = 2 + 6 + 1 // interface (2) + obis (6) + attribute (1)

// NewCosemAttribute creates a new CosemAttribute
func NewCosemAttribute(interfaceClass enumerations.CosemInterface, instance *Obis, attribute uint8) *CosemAttribute {
	return &CosemAttribute{
		Interface: interfaceClass,
		Instance:  instance,
		Attribute: attribute,
	}
}

// FromBytes creates a CosemAttribute from bytes
func (c *CosemAttribute) FromBytes(sourceBytes []byte) (*CosemAttribute, error) {
	if len(sourceBytes) != CosemAttributeLength {
		return nil, fmt.Errorf("data is not of correct length. Should be %d but is %d", CosemAttributeLength, len(sourceBytes))
	}

	interfaceClass := enumerations.CosemInterface(binary.BigEndian.Uint16(sourceBytes[:2]))
	instance, err := FromBytes(sourceBytes[2:8])
	if err != nil {
		return nil, fmt.Errorf("failed to parse OBIS: %w", err)
	}
	attribute := sourceBytes[8]

	return &CosemAttribute{
		Interface: interfaceClass,
		Instance:  instance,
		Attribute: attribute,
	}, nil
}

// ToBytes converts CosemAttribute to bytes
func (c *CosemAttribute) ToBytes() []byte {
	result := make([]byte, CosemAttributeLength)
	binary.BigEndian.PutUint16(result[0:2], uint16(c.Interface))
	copy(result[2:8], c.Instance.ToBytes())
	result[8] = c.Attribute
	return result
}

// CosemMethod represents a COSEM method descriptor
type CosemMethod struct {
	Interface enumerations.CosemInterface
	Instance  *Obis
	Method    uint8
}

const CosemMethodLength = 2 + 6 + 1 // interface (2) + obis (6) + method (1)

// NewCosemMethod creates a new CosemMethod
func NewCosemMethod(interfaceClass enumerations.CosemInterface, instance *Obis, method uint8) *CosemMethod {
	return &CosemMethod{
		Interface: interfaceClass,
		Instance:  instance,
		Method:    method,
	}
}

// FromBytes creates a CosemMethod from bytes
func (c *CosemMethod) FromBytes(sourceBytes []byte) (*CosemMethod, error) {
	if len(sourceBytes) != CosemMethodLength {
		return nil, fmt.Errorf("data is not of correct length. Should be %d but is %d", CosemMethodLength, len(sourceBytes))
	}

	interfaceClass := enumerations.CosemInterface(binary.BigEndian.Uint16(sourceBytes[:2]))
	instance, err := FromBytes(sourceBytes[2:8])
	if err != nil {
		return nil, fmt.Errorf("failed to parse OBIS: %w", err)
	}
	method := sourceBytes[8]

	return &CosemMethod{
		Interface: interfaceClass,
		Instance:  instance,
		Method:    method,
	}, nil
}

// ToBytes converts CosemMethod to bytes
func (c *CosemMethod) ToBytes() []byte {
	result := make([]byte, CosemMethodLength)
	binary.BigEndian.PutUint16(result[0:2], uint16(c.Interface))
	copy(result[2:8], c.Instance.ToBytes())
	result[8] = c.Method
	return result
}

