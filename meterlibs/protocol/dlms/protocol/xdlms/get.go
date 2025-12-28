package xdlms

import (
	"encoding/binary"
	"fmt"

	"github.com/yimiliya/idis/meterlibs/protocol/dlms/cosem"
	"github.com/yimiliya/idis/meterlibs/protocol/dlms/dlmsdata"
	"github.com/yimiliya/idis/meterlibs/protocol/dlms/enumerations"
)

// GetRequestNormal represents a Get request normal
// Get requests work in single attributes on interface classes.
// To get a value you would need the interface class, the instance (OBIS) and the attribute id.
// Some attributes allow for selective access to the attributes.
const GetRequestTag = 192

type GetRequestNormal struct {
	*BaseXDlmsApdu
	CosemAttribute      *cosem.CosemAttribute
	InvokeIdAndPriority *InvokeIdAndPriority
	AccessSelection     interface{} // RangeDescriptor or EntryDescriptor
}

// NewGetRequestNormal creates a new GetRequestNormal
func NewGetRequestNormal(
	cosemAttribute *cosem.CosemAttribute,
	invokeIdAndPriority *InvokeIdAndPriority,
	accessSelection interface{},
) *GetRequestNormal {
	return &GetRequestNormal{
		BaseXDlmsApdu: &BaseXDlmsApdu{
			Tag: GetRequestTag,
		},
		CosemAttribute:      cosemAttribute,
		InvokeIdAndPriority: invokeIdAndPriority,
		AccessSelection:     accessSelection,
	}
}

// FromBytes creates GetRequestNormal from bytes
func (g *GetRequestNormal) FromBytes(data []byte) (*GetRequestNormal, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for GetRequest")
	}

	tag := data[0]
	if tag != GetRequestTag {
		return nil, fmt.Errorf("tag for GET request is not correct. Got %d, should be %d", tag, GetRequestTag)
	}

	typeChoice := enumerations.GetRequestType(data[1])
	if typeChoice != enumerations.GetRequestTypeNormal {
		return nil, fmt.Errorf("the data for the GetRequest is not for a GetRequestNormal")
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
			factory := cosem.NewAccessDescriptorFactory()
			var consumed int
			accessSelection, consumed, err = factory.FromBytes(data)
			if err != nil {
				return nil, fmt.Errorf("failed to parse access selection: %w", err)
			}
			data = data[consumed:]
		}
	}

	return NewGetRequestNormal(cosemAttribute, invokeIdAndPriority, accessSelection), nil
}

// ToBytes converts GetRequestNormal to bytes
func (g *GetRequestNormal) ToBytes() ([]byte, error) {
	result := []byte{GetRequestTag}
	result = append(result, byte(enumerations.GetRequestTypeNormal))

	invokeBytes := g.InvokeIdAndPriority.ToBytes()
	result = append(result, invokeBytes...)

	cosemBytes := g.CosemAttribute.ToBytes()
	result = append(result, cosemBytes...)

	if g.AccessSelection != nil {
		result = append(result, 0x01)
		switch sel := g.AccessSelection.(type) {
		case *cosem.RangeDescriptor:
			result = append(result, sel.ToBytes()...)
		case *cosem.EntryDescriptor:
			result = append(result, sel.ToBytes()...)
		default:
			return nil, fmt.Errorf("unknown access selection type: %T", g.AccessSelection)
		}
	} else {
		result = append(result, 0x00)
	}

	return result, nil
}

// GetRequestNext represents a Get request next (for block transfer)
type GetRequestNext struct {
	*BaseXDlmsApdu
	BlockNumber         uint32
	InvokeIdAndPriority *InvokeIdAndPriority
}

// NewGetRequestNext creates a new GetRequestNext
func NewGetRequestNext(blockNumber uint32, invokeIdAndPriority *InvokeIdAndPriority) *GetRequestNext {
	return &GetRequestNext{
		BaseXDlmsApdu: &BaseXDlmsApdu{
			Tag: GetRequestTag,
		},
		BlockNumber:         blockNumber,
		InvokeIdAndPriority: invokeIdAndPriority,
	}
}

// FromBytes creates GetRequestNext from bytes
func (g *GetRequestNext) FromBytes(data []byte) (*GetRequestNext, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for GetRequestNext")
	}

	tag := data[0]
	if tag != GetRequestTag {
		return nil, fmt.Errorf("tag for GET request is not correct. Got %d, should be %d", tag, GetRequestTag)
	}

	typeChoice := enumerations.GetRequestType(data[1])
	if typeChoice != enumerations.GetRequestTypeNext {
		return nil, fmt.Errorf("the data for the GetRequest is not for a GetRequestNext")
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

	// Parse block_number (4 bytes)
	if len(data) < 4 {
		return nil, fmt.Errorf("insufficient data for block_number")
	}
	blockNumber := binary.BigEndian.Uint32(data[:4])

	return NewGetRequestNext(blockNumber, invokeIdAndPriority), nil
}

// ToBytes converts GetRequestNext to bytes
func (g *GetRequestNext) ToBytes() ([]byte, error) {
	result := []byte{GetRequestTag}
	result = append(result, byte(enumerations.GetRequestTypeNext))

	invokeBytes := g.InvokeIdAndPriority.ToBytes()
	result = append(result, invokeBytes...)

	blockBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(blockBytes, g.BlockNumber)
	result = append(result, blockBytes...)

	return result, nil
}

// GetResponseNormal represents a Get response normal
const GetResponseTag = 196

type GetResponseNormal struct {
	*BaseXDlmsApdu
	InvokeIdAndPriority *InvokeIdAndPriority
	Data                []byte
}

// NewGetResponseNormal creates a new GetResponseNormal
func NewGetResponseNormal(
	invokeIdAndPriority *InvokeIdAndPriority,
	data []byte,
) *GetResponseNormal {
	return &GetResponseNormal{
		BaseXDlmsApdu: &BaseXDlmsApdu{
			Tag: GetResponseTag,
		},
		InvokeIdAndPriority: invokeIdAndPriority,
		Data:                data,
	}
}

// FromBytes creates GetResponseNormal from bytes
func (g *GetResponseNormal) FromBytes(data []byte) (*GetResponseNormal, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for GetResponse")
	}

	tag := data[0]
	if tag != GetResponseTag {
		return nil, fmt.Errorf("tag for GET response is not correct. Got %d, should be %d", tag, GetResponseTag)
	}

	typeChoice := enumerations.GetResponseType(data[1])
	if typeChoice != enumerations.GetResponseTypeNormal {
		return nil, fmt.Errorf("the data for the GetResponse is not for a GetResponseNormal")
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

	// Parse choice (0 = data, 1 = error)
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for choice")
	}
	choice := data[0]
	if choice != 0 {
		return nil, fmt.Errorf("the data choice is not 0 to indicate data but: %d", choice)
	}
	data = data[1:]

	// Parse data (remaining bytes)
	responseData := make([]byte, len(data))
	copy(responseData, data)

	return NewGetResponseNormal(invokeIdAndPriority, responseData), nil
}

// ToBytes converts GetResponseNormal to bytes
func (g *GetResponseNormal) ToBytes() ([]byte, error) {
	result := []byte{GetResponseTag}
	result = append(result, byte(enumerations.GetResponseTypeNormal))

	invokeBytes := g.InvokeIdAndPriority.ToBytes()
	result = append(result, invokeBytes...)

	result = append(result, 0) // data result choice = 0 (data)
	result = append(result, g.Data...)

	return result, nil
}

// GetResponseNormalWithError represents a Get response normal with error
type GetResponseNormalWithError struct {
	*BaseXDlmsApdu
	InvokeIdAndPriority *InvokeIdAndPriority
	Error               enumerations.DataAccessResult
}

// NewGetResponseNormalWithError creates a new GetResponseNormalWithError
func NewGetResponseNormalWithError(
	invokeIdAndPriority *InvokeIdAndPriority,
	error enumerations.DataAccessResult,
) *GetResponseNormalWithError {
	return &GetResponseNormalWithError{
		BaseXDlmsApdu: &BaseXDlmsApdu{
			Tag: GetResponseTag,
		},
		InvokeIdAndPriority: invokeIdAndPriority,
		Error:               error,
	}
}

// FromBytes creates GetResponseNormalWithError from bytes
func (g *GetResponseNormalWithError) FromBytes(data []byte) (*GetResponseNormalWithError, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for GetResponse")
	}

	tag := data[0]
	if tag != GetResponseTag {
		return nil, fmt.Errorf("tag for GET response is not correct. Got %d, should be %d", tag, GetResponseTag)
	}

	typeChoice := enumerations.GetResponseType(data[1])
	if typeChoice != enumerations.GetResponseTypeNormal {
		return nil, fmt.Errorf("the data for the GetResponse is not for a GetResponseNormal")
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

	// Parse choice (0 = data, 1 = error)
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for choice")
	}
	choice := data[0]
	if choice != 1 {
		return nil, fmt.Errorf("the data choice is not 1 to indicate error but: %d", choice)
	}
	data = data[1:]

	// Parse error
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for error")
	}
	error := enumerations.DataAccessResult(data[0])

	return NewGetResponseNormalWithError(invokeIdAndPriority, error), nil
}

// ToBytes converts GetResponseNormalWithError to bytes
func (g *GetResponseNormalWithError) ToBytes() ([]byte, error) {
	result := []byte{GetResponseTag}
	result = append(result, byte(enumerations.GetResponseTypeNormal))

	invokeBytes := g.InvokeIdAndPriority.ToBytes()
	result = append(result, invokeBytes...)

	result = append(result, 1) // data error choice = 1 (error)
	result = append(result, byte(g.Error))

	return result, nil
}

// GetResponseWithDataBlock represents a Get response with data block
type GetResponseWithDataBlock struct {
	*BaseXDlmsApdu
	InvokeIdAndPriority *InvokeIdAndPriority
	LastBlock           bool
	BlockNumber         uint32
	RawData             []byte
}

// NewGetResponseWithDataBlock creates a new GetResponseWithDataBlock
func NewGetResponseWithDataBlock(
	invokeIdAndPriority *InvokeIdAndPriority,
	lastBlock bool,
	blockNumber uint32,
	rawData []byte,
) *GetResponseWithDataBlock {
	return &GetResponseWithDataBlock{
		BaseXDlmsApdu: &BaseXDlmsApdu{
			Tag: GetResponseTag,
		},
		InvokeIdAndPriority: invokeIdAndPriority,
		LastBlock:           lastBlock,
		BlockNumber:         blockNumber,
		RawData:             rawData,
	}
}

// FromBytes creates GetResponseWithDataBlock from bytes
func (g *GetResponseWithDataBlock) FromBytes(data []byte) (*GetResponseWithDataBlock, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for GetResponseWithDataBlock")
	}

	tag := data[0]
	if tag != GetResponseTag {
		return nil, fmt.Errorf("tag for GET response is not correct. Got %d, should be %d", tag, GetResponseTag)
	}

	typeChoice := enumerations.GetResponseType(data[1])
	if typeChoice != enumerations.GetResponseWithBlock {
		return nil, fmt.Errorf("the data for the GetResponse is not for a GetResponseWithDataBlock")
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

	// Parse last_block (1 byte boolean)
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for last_block")
	}
	lastBlock := data[0] != 0
	data = data[1:]

	// Parse block_number (4 bytes)
	if len(data) < 4 {
		return nil, fmt.Errorf("insufficient data for block_number")
	}
	blockNumber := binary.BigEndian.Uint32(data[:4])
	data = data[4:]

	// Parse raw_data length and data
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for raw_data length")
	}
	rawDataLength := int(data[0])
	data = data[1:]
	if len(data) < rawDataLength {
		return nil, fmt.Errorf("insufficient data for raw_data")
	}
	rawData := make([]byte, rawDataLength)
	copy(rawData, data[:rawDataLength])

	return NewGetResponseWithDataBlock(invokeIdAndPriority, lastBlock, blockNumber, rawData), nil
}

// ToBytes converts GetResponseWithDataBlock to bytes
func (g *GetResponseWithDataBlock) ToBytes() ([]byte, error) {
	result := []byte{GetResponseTag}
	result = append(result, byte(enumerations.GetResponseWithBlock))

	invokeBytes := g.InvokeIdAndPriority.ToBytes()
	result = append(result, invokeBytes...)

	if g.LastBlock {
		result = append(result, 0x01)
	} else {
		result = append(result, 0x00)
	}

	blockBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(blockBytes, g.BlockNumber)
	result = append(result, blockBytes...)

	result = append(result, byte(len(g.RawData)))
	result = append(result, g.RawData...)

	return result, nil
}

// GetRequestWithList represents a Get request with list
type GetRequestWithList struct {
	*BaseXDlmsApdu
	InvokeIdAndPriority *InvokeIdAndPriority
	Attributes          []*cosem.CosemAttribute
	AccessSelections    []interface{} // Optional access selections for each attribute
}

// NewGetRequestWithList creates a new GetRequestWithList
func NewGetRequestWithList(
	invokeIdAndPriority *InvokeIdAndPriority,
	attributes []*cosem.CosemAttribute,
	accessSelections []interface{},
) *GetRequestWithList {
	return &GetRequestWithList{
		BaseXDlmsApdu: &BaseXDlmsApdu{
			Tag: GetRequestTag,
		},
		InvokeIdAndPriority: invokeIdAndPriority,
		Attributes:          attributes,
		AccessSelections:    accessSelections,
	}
}

// FromBytes creates GetRequestWithList from bytes
func (g *GetRequestWithList) FromBytes(data []byte) (*GetRequestWithList, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for GetRequestWithList")
	}

	tag := data[0]
	if tag != GetRequestTag {
		return nil, fmt.Errorf("tag for GET request is not correct. Got %d, should be %d", tag, GetRequestTag)
	}

	typeChoice := enumerations.GetRequestType(data[1])
	if typeChoice != enumerations.GetRequestTypeWithList {
		return nil, fmt.Errorf("the data for the GetRequest is not for a GetRequestWithList")
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

	// Parse attribute descriptor list count
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for attribute descriptor list count")
	}
	attributeCount := int(data[0])
	data = data[1:]

	attributes := make([]*cosem.CosemAttribute, 0, attributeCount)
	accessSelections := make([]interface{}, 0, attributeCount)

	for i := 0; i < attributeCount; i++ {
		if len(data) < 9 {
			return nil, fmt.Errorf("insufficient data for attribute %d", i)
		}
		cosemAttribute, err := (&cosem.CosemAttribute{}).FromBytes(data[:9])
		if err != nil {
			return nil, fmt.Errorf("failed to parse cosem_attribute %d: %w", i, err)
		}
		attributes = append(attributes, cosemAttribute)
		data = data[9:]

		// Parse access_selection (optional)
		var accessSelection interface{}
		if len(data) > 0 {
			hasAccessSelection := data[0] != 0
			data = data[1:]
			if hasAccessSelection {
				factory := cosem.NewAccessDescriptorFactory()
				var consumed int
				accessSelection, consumed, err = factory.FromBytes(data)
				if err != nil {
					return nil, fmt.Errorf("failed to parse access selection %d: %w", i, err)
				}
				data = data[consumed:]
			}
		}
		accessSelections = append(accessSelections, accessSelection)
	}

	return NewGetRequestWithList(invokeIdAndPriority, attributes, accessSelections), nil
}

// ToBytes converts GetRequestWithList to bytes
func (g *GetRequestWithList) ToBytes() ([]byte, error) {
	result := []byte{GetRequestTag}
	result = append(result, byte(enumerations.GetRequestTypeWithList))

	invokeBytes := g.InvokeIdAndPriority.ToBytes()
	result = append(result, invokeBytes...)

	result = append(result, byte(len(g.Attributes)))

	for i, attr := range g.Attributes {
		cosemBytes := attr.ToBytes()
		result = append(result, cosemBytes...)

		if i < len(g.AccessSelections) && g.AccessSelections[i] != nil {
			result = append(result, 0x01)
			switch sel := g.AccessSelections[i].(type) {
			case *cosem.RangeDescriptor:
				result = append(result, sel.ToBytes()...)
			case *cosem.EntryDescriptor:
				result = append(result, sel.ToBytes()...)
			default:
				return nil, fmt.Errorf("unknown access selection type: %T", sel)
			}
		} else {
			result = append(result, 0x00)
		}
	}

	return result, nil
}

// GetDataResult represents a single result in GetResponseWithList
type GetDataResult struct {
	Data  []byte
	Error enumerations.DataAccessResult
}

// GetResponseWithList represents a Get response with list
type GetResponseWithList struct {
	*BaseXDlmsApdu
	InvokeIdAndPriority *InvokeIdAndPriority
	Results             []*GetDataResult
}

// NewGetResponseWithList creates a new GetResponseWithList
func NewGetResponseWithList(
	invokeIdAndPriority *InvokeIdAndPriority,
	results []*GetDataResult,
) *GetResponseWithList {
	return &GetResponseWithList{
		BaseXDlmsApdu: &BaseXDlmsApdu{
			Tag: GetResponseTag,
		},
		InvokeIdAndPriority: invokeIdAndPriority,
		Results:             results,
	}
}

// FromBytes creates GetResponseWithList from bytes
func (g *GetResponseWithList) FromBytes(data []byte) (*GetResponseWithList, error) {
	// TODO: Implement full parsing
	return nil, fmt.Errorf("GetResponseWithList.FromBytes not yet implemented")
}

// ToBytes converts GetResponseWithList to bytes
func (g *GetResponseWithList) ToBytes() ([]byte, error) {
	// TODO: Implement full encoding
	return nil, fmt.Errorf("GetResponseWithList.ToBytes not yet implemented")
}

// GetResponseLastBlock represents a Get response last block
type GetResponseLastBlock struct {
	*BaseXDlmsApdu
	InvokeIdAndPriority *InvokeIdAndPriority
	BlockNumber         uint32
	RawData             []byte
}

// NewGetResponseLastBlock creates a new GetResponseLastBlock
func NewGetResponseLastBlock(
	invokeIdAndPriority *InvokeIdAndPriority,
	blockNumber uint32,
	rawData []byte,
) *GetResponseLastBlock {
	return &GetResponseLastBlock{
		BaseXDlmsApdu: &BaseXDlmsApdu{
			Tag: GetResponseTag,
		},
		InvokeIdAndPriority: invokeIdAndPriority,
		BlockNumber:         blockNumber,
		RawData:             rawData,
	}
}

// FromBytes creates GetResponseLastBlock from bytes
func (g *GetResponseLastBlock) FromBytes(data []byte) (*GetResponseLastBlock, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for GetResponseLastBlock")
	}

	tag := data[0]
	if tag != GetResponseTag {
		return nil, fmt.Errorf("tag for GET response is not correct. Got %d, should be %d", tag, GetResponseTag)
	}

	typeChoice := enumerations.GetResponseType(data[1])
	if typeChoice != enumerations.GetResponseTypeLastBlock {
		return nil, fmt.Errorf("the data for the GetResponse is not for a GetResponseLastBlock")
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

	// Parse block_number (4 bytes)
	if len(data) < 4 {
		return nil, fmt.Errorf("insufficient data for block_number")
	}
	blockNumber := binary.BigEndian.Uint32(data[:4])
	data = data[4:]

	// Parse raw_data length and data
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for raw_data length")
	}
	rawDataLength := int(data[0])
	data = data[1:]
	if len(data) < rawDataLength {
		return nil, fmt.Errorf("insufficient data for raw_data")
	}
	rawData := make([]byte, rawDataLength)
	copy(rawData, data[:rawDataLength])

	return NewGetResponseLastBlock(invokeIdAndPriority, blockNumber, rawData), nil
}

// ToBytes converts GetResponseLastBlock to bytes
func (g *GetResponseLastBlock) ToBytes() ([]byte, error) {
	result := []byte{GetResponseTag}
	result = append(result, byte(enumerations.GetResponseTypeLastBlock))

	invokeBytes := g.InvokeIdAndPriority.ToBytes()
	result = append(result, invokeBytes...)

	blockBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(blockBytes, g.BlockNumber)
	result = append(result, blockBytes...)

	result = append(result, byte(len(g.RawData)))
	result = append(result, g.RawData...)

	return result, nil
}

// GetResponseLastBlockWithError represents a Get response last block with error
type GetResponseLastBlockWithError struct {
	*BaseXDlmsApdu
	InvokeIdAndPriority *InvokeIdAndPriority
	BlockNumber         uint32
	Error               enumerations.DataAccessResult
}

// NewGetResponseLastBlockWithError creates a new GetResponseLastBlockWithError
func NewGetResponseLastBlockWithError(
	invokeIdAndPriority *InvokeIdAndPriority,
	blockNumber uint32,
	error enumerations.DataAccessResult,
) *GetResponseLastBlockWithError {
	return &GetResponseLastBlockWithError{
		BaseXDlmsApdu: &BaseXDlmsApdu{
			Tag: GetResponseTag,
		},
		InvokeIdAndPriority: invokeIdAndPriority,
		BlockNumber:         blockNumber,
		Error:               error,
	}
}

// FromBytes creates GetResponseLastBlockWithError from bytes
func (g *GetResponseLastBlockWithError) FromBytes(data []byte) (*GetResponseLastBlockWithError, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for GetResponseLastBlockWithError")
	}

	tag := data[0]
	if tag != GetResponseTag {
		return nil, fmt.Errorf("tag for GET response is not correct. Got %d, should be %d", tag, GetResponseTag)
	}

	typeChoice := enumerations.GetResponseType(data[1])
	if typeChoice != enumerations.GetResponseTypeLastBlockWithError {
		return nil, fmt.Errorf("the data for the GetResponse is not for a GetResponseLastBlockWithError")
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

	// Parse block_number (4 bytes)
	if len(data) < 4 {
		return nil, fmt.Errorf("insufficient data for block_number")
	}
	blockNumber := binary.BigEndian.Uint32(data[:4])
	data = data[4:]

	// Parse error
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for error")
	}
	error := enumerations.DataAccessResult(data[0])

	return NewGetResponseLastBlockWithError(invokeIdAndPriority, blockNumber, error), nil
}

// ToBytes converts GetResponseLastBlockWithError to bytes
func (g *GetResponseLastBlockWithError) ToBytes() ([]byte, error) {
	result := []byte{GetResponseTag}
	result = append(result, byte(enumerations.GetResponseTypeLastBlockWithError))

	invokeBytes := g.InvokeIdAndPriority.ToBytes()
	result = append(result, invokeBytes...)

	blockBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(blockBytes, g.BlockNumber)
	result = append(result, blockBytes...)

	result = append(result, byte(g.Error))

	return result, nil
}
