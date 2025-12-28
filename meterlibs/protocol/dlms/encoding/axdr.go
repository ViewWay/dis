package encoding

import (
	"fmt"

	"github.com/yimiliya/idis/meterlibs/protocol/dlms/dlmsdata"
)

const VariableLength = -1

// GetAXdrLength finds the length of an XADR element assuming the length is the first bytes
// Works with bytearray and will remove element from the array as it finds the variable length.
func GetAXdrLength(data []byte) (int, []byte, error) {
	if len(data) == 0 {
		return 0, nil, fmt.Errorf("insufficient data for AXDR length")
	}
	
	firstByte := data[0]
	lengthIsMultipleBytes := (firstByte & 0b10000000) != 0
	
	if !lengthIsMultipleBytes {
		return int(firstByte), data[1:], nil
	}
	
	numberOfBytes := int(firstByte & 0b01111111)
	if len(data) < numberOfBytes+1 {
		return 0, nil, fmt.Errorf("insufficient data for AXDR length: need %d bytes, got %d", numberOfBytes+1, len(data))
	}
	
	lengthBytes := data[1 : numberOfBytes+1]
	length := 0
	for _, b := range lengthBytes {
		length = (length << 8) | int(b)
	}
	
	return length, data[numberOfBytes+1:], nil
}

// Attribute represents an attribute in encoding configuration
type Attribute struct {
	AttributeName string
	CreateInstance func([]byte) (interface{}, error)
	Length        int
	ReturnValue   bool
	WrapEnd       bool
	Default       interface{}
	Optional      bool
}

// Sequence represents a sequence in encoding configuration
type Sequence struct {
	AttributeName string
	InstanceFactory interface{} // DlmsDataFactory or similar
}

// Choice represents a choice in encoding configuration
type Choice struct {
	Choices map[byte]interface{} // byte -> Attribute or Sequence
}

// EncodingConf represents encoding configuration
type EncodingConf struct {
	Attributes []interface{} // Attribute, Sequence, or Choice
}

// AXdrDecoder decodes A-XDR encoded data
type AXdrDecoder struct {
	EncodingConf *EncodingConf
	Buffer       []byte
	Pointer      int
	Result       map[string]interface{}
}

// NewAXdrDecoder creates a new AXdrDecoder
func NewAXdrDecoder(encodingConf *EncodingConf) *AXdrDecoder {
	return &AXdrDecoder{
		EncodingConf: encodingConf,
		Buffer:       make([]byte, 0),
		Pointer:      0,
		Result:       make(map[string]interface{}),
	}
}

// BufferEmpty returns true if buffer is empty
func (a *AXdrDecoder) BufferEmpty() bool {
	return a.Pointer >= len(a.Buffer)
}

// Decode decodes data according to encoding configuration
func (a *AXdrDecoder) Decode(data []byte) (map[string]interface{}, error) {
	// Clear previous results
	a.Result = make(map[string]interface{})
	// Fill the buffer
	a.Buffer = make([]byte, len(data))
	copy(a.Buffer, data)
	a.Pointer = 0
	
	for index, dataAttribute := range a.EncodingConf.Attributes {
		result, err := a.DecodeSingle(dataAttribute, index)
		if err != nil {
			return nil, err
		}
		for k, v := range result {
			a.Result[k] = v
		}
	}
	
	return a.Result, nil
}

// IsLastEncodingElement checks if this is the last element
func (a *AXdrDecoder) IsLastEncodingElement(index int) bool {
	return index == len(a.EncodingConf.Attributes)-1
}

// GetBufferTail returns remaining buffer
func (a *AXdrDecoder) GetBufferTail() []byte {
	return a.Buffer[a.Pointer:]
}

// DecodeSingle decodes a single element
func (a *AXdrDecoder) DecodeSingle(dataType interface{}, index int) (map[string]interface{}, error) {
	switch t := dataType.(type) {
	case *Attribute:
		value, err := a.DecodeAttribute(t, index)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{t.AttributeName: value}, nil
	case *Choice:
		choiceByte, err := a.GetBytes(1)
		if err != nil {
			return nil, err
		}
		choice, ok := t.Choices[choiceByte[0]]
		if !ok {
			return nil, fmt.Errorf("unknown choice value: %d", choiceByte[0])
		}
		return a.DecodeSingle(choice, index)
	case *Sequence:
		return a.DecodeSequence(t)
	default:
		return nil, fmt.Errorf("no valid class type")
	}
}

// DecodeAttribute decodes an attribute
func (a *AXdrDecoder) DecodeAttribute(attribute *Attribute, index int) (interface{}, error) {
	if attribute.Optional {
		indicator, err := a.GetBytes(1)
		if err != nil {
			return nil, err
		}
		if indicator[0] == 0x00 {
			// Not used
			return nil, nil
		}
	}
	
	if attribute.Default != nil {
		indicator, err := a.GetBytes(1)
		if err != nil {
			return nil, err
		}
		if indicator[0] == 0x00 {
			// Use the default
			return attribute.Default, nil
		}
	}
	
	// Fixed length?
	if attribute.Length != VariableLength {
		data, err := a.GetBytes(attribute.Length)
		if err != nil {
			return nil, err
		}
		return attribute.CreateInstance(data)
	}
	
	// Check if last element
	if a.IsLastEncodingElement(index) {
		// Use all remaining data
		remaining := a.GetBufferTail()
		return attribute.CreateInstance(remaining)
	}
	
	// We know how to create the instance (just not how long it is)
	length, _, err := GetAXdrLength(a.GetBufferTail())
	if err != nil {
		return nil, err
	}
	data, err := a.GetBytes(length)
	if err != nil {
		return nil, err
	}
	return attribute.CreateInstance(data)
}

// DecodeSequence decodes a sequence
func (a *AXdrDecoder) DecodeSequence(seq *Sequence) (map[string]interface{}, error) {
	parsedData := make([]interface{}, 0)
	
	for !a.BufferEmpty() {
		tag, err := a.GetBytes(1)
		if err != nil {
			return nil, err
		}
		
		dataClass, err := dlmsdata.NewDlmsDataFactory().GetDataClass(dlmsdata.DlmsDataTag(tag[0]))
		if err != nil {
			return nil, err
		}
		
		instance := dataClass()
		
		switch instance.GetTag() {
		case dlmsdata.TagArray:
			arrayData, err := a.DecodeArray()
			if err != nil {
				return nil, err
			}
			parsedData = append(parsedData, arrayData)
			continue
		case dlmsdata.TagStructure:
			structureData, err := a.DecodeStructure()
			if err != nil {
				return nil, err
			}
			parsedData = append(parsedData, structureData)
			continue
		}
		
		if instance.GetLength() != VariableLength {
			data, err := a.GetBytes(instance.GetLength())
			if err != nil {
				return nil, err
			}
			decoded, err := instance.FromBytes(data)
			if err != nil {
				return nil, err
			}
			parsedData = append(parsedData, decoded.ToPython())
			continue
		}
		
		// Variable length
		length, _, err := GetAXdrLength(a.GetBufferTail())
		if err != nil {
			return nil, err
		}
		data, err := a.GetBytes(length)
		if err != nil {
			return nil, err
		}
		decoded, err := instance.FromBytes(data)
		if err != nil {
			return nil, err
		}
		parsedData = append(parsedData, decoded.ToPython())
	}
	
	if len(parsedData) == 1 {
		return map[string]interface{}{seq.AttributeName: parsedData[0]}, nil
	}
	
	return map[string]interface{}{seq.AttributeName: parsedData}, nil
}

// DecodeArray decodes an array
func (a *AXdrDecoder) DecodeArray() ([]interface{}, error) {
	itemCount, _, err := GetAXdrLength(a.GetBufferTail())
	if err != nil {
		return nil, err
	}
	
	elements := make([]interface{}, 0, itemCount)
	for i := 0; i < itemCount; i++ {
		element, err := a.DecodeSequenceOf()
		if err != nil {
			return nil, err
		}
		elements = append(elements, element)
	}
	return elements, nil
}

// DecodeStructure decodes a structure
func (a *AXdrDecoder) DecodeStructure() ([]interface{}, error) {
	itemCount, _, err := GetAXdrLength(a.GetBufferTail())
	if err != nil {
		return nil, err
	}
	
	elements := make([]interface{}, 0, itemCount)
	for i := 0; i < itemCount; i++ {
		element, err := a.DecodeSequenceOf()
		if err != nil {
			return nil, err
		}
		elements = append(elements, element)
	}
	return elements, nil
}

// DecodeSequenceOf decodes a sequence of elements
func (a *AXdrDecoder) DecodeSequenceOf() (interface{}, error) {
	tag, err := a.GetBytes(1)
	if err != nil {
		return nil, err
	}
	
	dataClass, err := dlmsdata.NewDlmsDataFactory().GetDataClass(dlmsdata.DlmsDataTag(tag[0]))
	if err != nil {
		return nil, err
	}
	
	instance := dataClass()
	
	switch instance.GetTag() {
	case dlmsdata.TagArray:
		return a.DecodeArray()
	case dlmsdata.TagStructure:
		return a.DecodeStructure()
	default:
		return a.DecodeData(instance)
	}
}

// DecodeData decodes a single data element
func (a *AXdrDecoder) DecodeData(dataClass func() dlmsdata.DlmsData) (interface{}, error) {
	instance := dataClass()
	
	if instance.GetLength() == VariableLength {
		length, _, err := GetAXdrLength(a.GetBufferTail())
		if err != nil {
			return nil, err
		}
		data, err := a.GetBytes(length)
		if err != nil {
			return nil, err
		}
		decoded, err := instance.FromBytes(data)
		if err != nil {
			return nil, err
		}
		return decoded.ToPython(), nil
	}
	
	data, err := a.GetBytes(instance.GetLength())
	if err != nil {
		return nil, err
	}
	decoded, err := instance.FromBytes(data)
	if err != nil {
		return nil, err
	}
	return decoded.ToPython(), nil
}

// GetBytes gets some bytes from the buffer and moves the pointer forward
func (a *AXdrDecoder) GetBytes(length int) ([]byte, error) {
	if a.Pointer+length > len(a.Buffer) {
		return nil, fmt.Errorf("insufficient data: need %d bytes, have %d", length, len(a.Buffer)-a.Pointer)
	}
	part := a.Buffer[a.Pointer : a.Pointer+length]
	a.Pointer += length
	return part, nil
}

// RemainingBuffer returns remaining buffer
func (a *AXdrDecoder) RemainingBuffer() []byte {
	return a.Buffer[a.Pointer:]
}

// GetAXdrLength gets the AXDR length from buffer
func (a *AXdrDecoder) GetAXdrLength() (int, error) {
	length, remaining, err := GetAXdrLength(a.GetBufferTail())
	if err != nil {
		return 0, err
	}
	// Update pointer
	a.Pointer += len(a.GetBufferTail()) - len(remaining)
	return length, nil
}

