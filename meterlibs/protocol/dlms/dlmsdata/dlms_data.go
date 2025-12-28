package dlmsdata

import (
	"encoding/binary"
	"fmt"
)

const VariableLength = -1

// DlmsDataTag represents DLMS data type tags
type DlmsDataTag uint8

const (
	TagNull               DlmsDataTag = 0
	TagArray              DlmsDataTag = 1
	TagStructure          DlmsDataTag = 2
	TagBoolean            DlmsDataTag = 3
	TagBitString          DlmsDataTag = 4
	TagDoubleLong         DlmsDataTag = 5
	TagDoubleLongUnsigned DlmsDataTag = 6
	TagOctetString        DlmsDataTag = 9
	TagVisibleString      DlmsDataTag = 10
	TagUTF8String         DlmsDataTag = 12
	TagBCD                DlmsDataTag = 13
	TagInteger            DlmsDataTag = 15
	TagLong               DlmsDataTag = 16
	TagUnsigned           DlmsDataTag = 17
	TagLongUnsigned       DlmsDataTag = 18
	TagCompactArray       DlmsDataTag = 19
	TagLong64             DlmsDataTag = 20
	TagLong64Unsigned     DlmsDataTag = 21
	TagEnum               DlmsDataTag = 22
	TagFloat32            DlmsDataTag = 23
	TagFloat64            DlmsDataTag = 24
	TagDateTime           DlmsDataTag = 25
	TagDate               DlmsDataTag = 26
	TagTime               DlmsDataTag = 27
	TagDontCare           DlmsDataTag = 255
)

// DlmsData is the interface for all DLMS data types
type DlmsData interface {
	ToPython() interface{}
	ToBytes() ([]byte, error)  // 统一返回error
	FromBytes(data []byte) (DlmsData, error)
	GetTag() DlmsDataTag
	GetLength() int
	String() string  // 添加String方法用于调试
}

// BaseDlmsData is the base struct for DLMS data types
type BaseDlmsData struct {
	Tag    DlmsDataTag
	Length int
	Value  interface{}
}

// GetTag returns the tag
func (b *BaseDlmsData) GetTag() DlmsDataTag {
	return b.Tag
}

// GetLength returns the length
func (b *BaseDlmsData) GetLength() int {
	return b.Length
}

// ToPython converts to Python-like value
func (b *BaseDlmsData) ToPython() interface{} {
	return b.Value
}

// ValueToBytes converts value to bytes (to be implemented by subclasses)
func (b *BaseDlmsData) ValueToBytes() ([]byte, error) {
	return nil, fmt.Errorf("value_to_bytes must be implemented in subclass")
}

// ToBytes converts the data to bytes
func (b *BaseDlmsData) ToBytes() ([]byte, error) {
	result := []byte{byte(b.Tag)}
	valueBytes, err := b.ValueToBytes()
	if err != nil {
		return nil, err
	}
	if b.Length == VariableLength {
		result = append(result, byte(len(valueBytes)))
	}
	result = append(result, valueBytes...)
	return result, nil
}

// NullData represents null data
type NullData struct {
	*BaseDlmsData
}

// NewNullData creates a new NullData
func NewNullData() *NullData {
	return &NullData{
		BaseDlmsData: &BaseDlmsData{
			Tag:    TagNull,
			Length: 0,
			Value:  nil,
		},
	}
}

// FromBytes creates NullData from bytes
func (n *NullData) FromBytes(data []byte) (DlmsData, error) {
	return NewNullData(), nil
}

// ToPython returns nil
func (n *NullData) ToPython() interface{} {
	return nil
}

// ToBytes returns empty bytes for null
func (n *NullData) ToBytes() ([]byte, error) {
	return []byte{byte(TagNull)}, nil
}

// ValueToBytes returns empty bytes
func (n *NullData) ValueToBytes() ([]byte, error) {
	return []byte{}, nil
}

// String returns string representation
func (n *NullData) String() string {
	return "null"
}

// BooleanData represents boolean data
type BooleanData struct {
	*BaseDlmsData
}

// NewBooleanData creates a new BooleanData
func NewBooleanData(value bool) *BooleanData {
	return &BooleanData{
		BaseDlmsData: &BaseDlmsData{
			Tag:    TagBoolean,
			Length: 1,
			Value:  value,
		},
	}
}

// FromBytes creates BooleanData from bytes
func (b *BooleanData) FromBytes(data []byte) (DlmsData, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for BooleanData")
	}
	value := data[0] != 0
	return NewBooleanData(value), nil
}

// ValueToBytes converts boolean to bytes
func (b *BooleanData) ValueToBytes() ([]byte, error) {
	if b.Value.(bool) {
		return []byte{1}, nil
	}
	return []byte{0}, nil
}

// String returns string representation
func (b *BooleanData) String() string {
	if b.Value.(bool) {
		return "true"
	}
	return "false"
}

// IntegerData represents 8-bit signed integer
type IntegerData struct {
	*BaseDlmsData
}

// NewIntegerData creates a new IntegerData
func NewIntegerData(value int8) *IntegerData {
	return &IntegerData{
		BaseDlmsData: &BaseDlmsData{
			Tag:    TagInteger,
			Length: 1,
			Value:  value,
		},
	}
}

// FromBytes creates IntegerData from bytes
func (i *IntegerData) FromBytes(data []byte) (DlmsData, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for IntegerData")
	}
	return NewIntegerData(int8(data[0])), nil
}

// ValueToBytes converts int8 to bytes
func (i *IntegerData) ValueToBytes() ([]byte, error) {
	return []byte{byte(i.Value.(int8))}, nil
}

// String returns string representation
func (i *IntegerData) String() string {
	return fmt.Sprintf("%d", i.Value.(int8))
}

// UnsignedIntegerData represents 8-bit unsigned integer
type UnsignedIntegerData struct {
	*BaseDlmsData
}

// NewUnsignedIntegerData creates a new UnsignedIntegerData
func NewUnsignedIntegerData(value uint8) *UnsignedIntegerData {
	return &UnsignedIntegerData{
		BaseDlmsData: &BaseDlmsData{
			Tag:    TagUnsigned,
			Length: 1,
			Value:  value,
		},
	}
}

// FromBytes creates UnsignedIntegerData from bytes
func (u *UnsignedIntegerData) FromBytes(data []byte) (DlmsData, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for UnsignedIntegerData")
	}
	return NewUnsignedIntegerData(data[0]), nil
}

// ValueToBytes converts uint8 to bytes
func (u *UnsignedIntegerData) ValueToBytes() ([]byte, error) {
	return []byte{u.Value.(uint8)}, nil
}

// String returns string representation
func (u *UnsignedIntegerData) String() string {
	return fmt.Sprintf("%d", u.Value.(uint8))
}

// LongData represents 16-bit signed integer
type LongData struct {
	*BaseDlmsData
}

// NewLongData creates a new LongData
func NewLongData(value int16) *LongData {
	return &LongData{
		BaseDlmsData: &BaseDlmsData{
			Tag:    TagLong,
			Length: 2,
			Value:  value,
		},
	}
}

// FromBytes creates LongData from bytes
func (l *LongData) FromBytes(data []byte) (DlmsData, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for LongData")
	}
	value := int16(binary.BigEndian.Uint16(data))
	return NewLongData(value), nil
}

// ValueToBytes converts int16 to bytes
func (l *LongData) ValueToBytes() ([]byte, error) {
	result := make([]byte, 2)
	binary.BigEndian.PutUint16(result, uint16(l.Value.(int16)))
	return result, nil
}

// String returns string representation
func (l *LongData) String() string {
	return fmt.Sprintf("%d", l.Value.(int16))
}

// UnsignedLongData represents 16-bit unsigned integer
type UnsignedLongData struct {
	*BaseDlmsData
}

// NewUnsignedLongData creates a new UnsignedLongData
func NewUnsignedLongData(value uint16) *UnsignedLongData {
	return &UnsignedLongData{
		BaseDlmsData: &BaseDlmsData{
			Tag:    TagLongUnsigned,
			Length: 2,
			Value:  value,
		},
	}
}

// FromBytes creates UnsignedLongData from bytes
func (u *UnsignedLongData) FromBytes(data []byte) (DlmsData, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for UnsignedLongData")
	}
	value := binary.BigEndian.Uint16(data)
	return NewUnsignedLongData(value), nil
}

// ValueToBytes converts uint16 to bytes
func (u *UnsignedLongData) ValueToBytes() ([]byte, error) {
	result := make([]byte, 2)
	binary.BigEndian.PutUint16(result, u.Value.(uint16))
	return result, nil
}

// String returns string representation
func (u *UnsignedLongData) String() string {
	return fmt.Sprintf("%d", u.Value.(uint16))
}

// DoubleLongData represents 32-bit signed integer
type DoubleLongData struct {
	*BaseDlmsData
}

// NewDoubleLongData creates a new DoubleLongData
func NewDoubleLongData(value int32) *DoubleLongData {
	return &DoubleLongData{
		BaseDlmsData: &BaseDlmsData{
			Tag:    TagDoubleLong,
			Length: 4,
			Value:  value,
		},
	}
}

// FromBytes creates DoubleLongData from bytes
func (d *DoubleLongData) FromBytes(data []byte) (DlmsData, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("insufficient data for DoubleLongData")
	}
	value := int32(binary.BigEndian.Uint32(data))
	return NewDoubleLongData(value), nil
}

// ValueToBytes converts int32 to bytes
func (d *DoubleLongData) ValueToBytes() ([]byte, error) {
	result := make([]byte, 4)
	binary.BigEndian.PutUint32(result, uint32(d.Value.(int32)))
	return result, nil
}

// String returns string representation
func (d *DoubleLongData) String() string {
	return fmt.Sprintf("%d", d.Value.(int32))
}

// DoubleLongUnsignedData represents 32-bit unsigned integer
type DoubleLongUnsignedData struct {
	*BaseDlmsData
}

// NewDoubleLongUnsignedData creates a new DoubleLongUnsignedData
func NewDoubleLongUnsignedData(value uint32) *DoubleLongUnsignedData {
	return &DoubleLongUnsignedData{
		BaseDlmsData: &BaseDlmsData{
			Tag:    TagDoubleLongUnsigned,
			Length: 4,
			Value:  value,
		},
	}
}

// FromBytes creates DoubleLongUnsignedData from bytes
func (d *DoubleLongUnsignedData) FromBytes(data []byte) (DlmsData, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("insufficient data for DoubleLongUnsignedData")
	}
	value := binary.BigEndian.Uint32(data)
	return NewDoubleLongUnsignedData(value), nil
}

// ValueToBytes converts uint32 to bytes
func (d *DoubleLongUnsignedData) ValueToBytes() ([]byte, error) {
	result := make([]byte, 4)
	binary.BigEndian.PutUint32(result, d.Value.(uint32))
	return result, nil
}

// String returns string representation
func (d *DoubleLongUnsignedData) String() string {
	return fmt.Sprintf("%d", d.Value.(uint32))
}

// OctetStringData represents octet string data
type OctetStringData struct {
	*BaseDlmsData
}

// NewOctetStringData creates a new OctetStringData
func NewOctetStringData(value []byte) *OctetStringData {
	return &OctetStringData{
		BaseDlmsData: &BaseDlmsData{
			Tag:    TagOctetString,
			Length: VariableLength,
			Value:  value,
		},
	}
}

// FromBytes creates OctetStringData from bytes
func (o *OctetStringData) FromBytes(data []byte) (DlmsData, error) {
	value := make([]byte, len(data))
	copy(value, data)
	return NewOctetStringData(value), nil
}

// ToPython returns the bytes value
func (o *OctetStringData) ToPython() interface{} {
	return o.Value.([]byte)
}

// ValueToBytes returns the bytes value
func (o *OctetStringData) ValueToBytes() ([]byte, error) {
	return o.Value.([]byte), nil
}

// String returns string representation
func (o *OctetStringData) String() string {
	return fmt.Sprintf("0x%x", o.Value.([]byte))
}

// VisibleStringData represents visible string data (ASCII)
type VisibleStringData struct {
	*BaseDlmsData
}

// NewVisibleStringData creates a new VisibleStringData
func NewVisibleStringData(value string) *VisibleStringData {
	return &VisibleStringData{
		BaseDlmsData: &BaseDlmsData{
			Tag:    TagVisibleString,
			Length: VariableLength,
			Value:  value,
		},
	}
}

// FromBytes creates VisibleStringData from bytes
func (v *VisibleStringData) FromBytes(data []byte) (DlmsData, error) {
	value := string(data)
	return NewVisibleStringData(value), nil
}

// ValueToBytes converts string to ASCII bytes
func (v *VisibleStringData) ValueToBytes() ([]byte, error) {
	return []byte(v.Value.(string)), nil
}

// String returns string representation
func (v *VisibleStringData) String() string {
	return fmt.Sprintf("\"%s\"", v.Value.(string))
}

// DataArray represents an array of DLMS data
type DataArray struct {
	*BaseDlmsData
}

// NewDataArray creates a new DataArray
func NewDataArray(value []DlmsData) *DataArray {
	return &DataArray{
		BaseDlmsData: &BaseDlmsData{
			Tag:    TagArray,
			Length: VariableLength,
			Value:  value,
		},
	}
}

// ToBytes converts DataArray to bytes
func (d *DataArray) ToBytes() ([]byte, error) {
	result := []byte{byte(TagArray)}
	items := d.Value.([]DlmsData)
	lengthBytes := EncodeVariableInteger(len(items))
	result = append(result, lengthBytes...)
	for _, item := range items {
		itemBytes, err := item.ToBytes()
		if err != nil {
			return nil, err
		}
		result = append(result, itemBytes...)
	}
	return result, nil
}

// ToPython converts to Python-like list
func (d *DataArray) ToPython() interface{} {
	items := d.Value.([]DlmsData)
	result := make([]interface{}, len(items))
	for i, item := range items {
		result[i] = item.ToPython()
	}
	return result
}

// FromBytes creates DataArray from bytes
// TODO: This is a simplified implementation. For full AXDR parsing, use AXdrDecoder
func (d *DataArray) FromBytes(data []byte) (DlmsData, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for DataArray tag")
	}
	if data[0] != byte(TagArray) {
		return nil, fmt.Errorf("invalid tag for DataArray: %d", data[0])
	}
	
	// Skip tag byte
	data = data[1:]
	
	// Decode length
	length, remaining, err := DecodeVariableInteger(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode array length: %w", err)
	}
	
	items := make([]DlmsData, 0, length)
	pos := 0
	
	for i := 0; i < length && pos < len(remaining); i++ {
		if pos >= len(remaining) {
			return nil, fmt.Errorf("insufficient data for array item %d", i)
		}
		
		tag := DlmsDataTag(remaining[pos])
		factory := NewDlmsDataFactory()
		itemFactory, err := factory.GetDataClass(tag)
		if err != nil {
			return nil, fmt.Errorf("unknown data tag in array: %d", tag)
		}
		
		item := itemFactory()
		
		// For variable length items, we need to parse them properly
		// For now, use a simple approach: try to parse the item
		itemData, err := item.FromBytes(remaining[pos:])
		if err != nil {
			return nil, fmt.Errorf("failed to parse array item %d: %w", i, err)
		}
		
		// Calculate consumed bytes
		itemBytes, err := itemData.ToBytes()
		if err != nil {
			return nil, fmt.Errorf("failed to encode array item %d: %w", i, err)
		}
		pos += len(itemBytes)
		items = append(items, itemData)
	}
	
	return NewDataArray(items), nil
}

// String returns string representation
func (d *DataArray) String() string {
	items := d.Value.([]DlmsData)
	result := "["
	for i, item := range items {
		if i > 0 {
			result += ", "
		}
		result += item.String()
	}
	result += "]"
	return result
}

// DataStructure represents a structure of DLMS data
type DataStructure struct {
	*BaseDlmsData
}

// NewDataStructure creates a new DataStructure
func NewDataStructure(value []DlmsData) *DataStructure {
	return &DataStructure{
		BaseDlmsData: &BaseDlmsData{
			Tag:    TagStructure,
			Length: VariableLength,
			Value:  value,
		},
	}
}

// ToBytes converts DataStructure to bytes
func (d *DataStructure) ToBytes() ([]byte, error) {
	result := []byte{byte(TagStructure)}
	items := d.Value.([]DlmsData)
	lengthBytes := EncodeVariableInteger(len(items))
	result = append(result, lengthBytes...)
	for _, item := range items {
		itemBytes, err := item.ToBytes()
		if err != nil {
			return nil, err
		}
		result = append(result, itemBytes...)
	}
	return result, nil
}

// ToPython converts to Python-like list
func (d *DataStructure) ToPython() interface{} {
	items := d.Value.([]DlmsData)
	result := make([]interface{}, len(items))
	for i, item := range items {
		result[i] = item.ToPython()
	}
	return result
}

// FromBytes creates DataStructure from bytes
// TODO: This is a simplified implementation. For full AXDR parsing, use AXdrDecoder
func (d *DataStructure) FromBytes(data []byte) (DlmsData, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for DataStructure tag")
	}
	if data[0] != byte(TagStructure) {
		return nil, fmt.Errorf("invalid tag for DataStructure: %d", data[0])
	}
	
	// Skip tag byte
	data = data[1:]
	
	// Decode length
	length, remaining, err := DecodeVariableInteger(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode structure length: %w", err)
	}
	
	items := make([]DlmsData, 0, length)
	pos := 0
	
	for i := 0; i < length && pos < len(remaining); i++ {
		if pos >= len(remaining) {
			return nil, fmt.Errorf("insufficient data for structure item %d", i)
		}
		
		tag := DlmsDataTag(remaining[pos])
		factory := NewDlmsDataFactory()
		itemFactory, err := factory.GetDataClass(tag)
		if err != nil {
			return nil, fmt.Errorf("unknown data tag in structure: %d", tag)
		}
		
		item := itemFactory()
		
		// Try to parse the item
		itemData, err := item.FromBytes(remaining[pos:])
		if err != nil {
			return nil, fmt.Errorf("failed to parse structure item %d: %w", i, err)
		}
		
		// Calculate consumed bytes
		itemBytes, err := itemData.ToBytes()
		if err != nil {
			return nil, fmt.Errorf("failed to encode structure item %d: %w", i, err)
		}
		pos += len(itemBytes)
		items = append(items, itemData)
	}
	
	return NewDataStructure(items), nil
}

// String returns string representation
func (d *DataStructure) String() string {
	items := d.Value.([]DlmsData)
	result := "{"
	for i, item := range items {
		if i > 0 {
			result += ", "
		}
		result += item.String()
	}
	result += "}"
	return result
}

// EncodeVariableInteger encodes a variable length integer
// If the length fits in 7 bits, it can be encoded in 1 byte.
// If it is larger, the last bit of the first byte indicates
// that the length of the length is encoded in the first byte
// and the length is encoded in the following bytes.
func EncodeVariableInteger(length int) []byte {
	if length <= 0x7F {
		return []byte{byte(length)}
	}
	
	encodedLength := 1
	for {
		maxValue := (1 << (8 * encodedLength)) - 1
		if length <= maxValue {
			break
		}
		encodedLength++
	}
	
	lengthByte := byte(0x80 | encodedLength)
	result := []byte{lengthByte}
	
	lengthBytes := make([]byte, encodedLength)
	for i := encodedLength - 1; i >= 0; i-- {
		lengthBytes[i] = byte(length & 0xFF)
		length >>= 8
	}
	result = append(result, lengthBytes...)
	
	return result
}

// DecodeVariableInteger decodes a variable length integer
func DecodeVariableInteger(data []byte) (int, []byte, error) {
	if len(data) == 0 {
		return 0, nil, fmt.Errorf("insufficient data for variable integer")
	}
	
	firstByte := data[0]
	isMultipleBytes := (firstByte & 0x80) != 0
	
	if !isMultipleBytes {
		length := int(firstByte & 0x7F)
		return length, data[1:], nil
	}
	
	lengthLength := int(firstByte & 0x7F)
	if len(data) < lengthLength+1 {
		return 0, nil, fmt.Errorf("insufficient data for variable integer length")
	}
	
	lengthBytes := data[1 : lengthLength+1]
	length := 0
	for _, b := range lengthBytes {
		length = (length << 8) | int(b)
	}
	
	return length, data[lengthLength+1:], nil
}

// DlmsDataFactory creates DLMS data instances from tags
type DlmsDataFactory struct{}

var dataFactoryMap = map[DlmsDataTag]func() DlmsData{
	TagNull:               func() DlmsData { return NewNullData() },
	TagArray:              func() DlmsData { return NewDataArray(nil) },
	TagStructure:          func() DlmsData { return NewDataStructure(nil) },
	TagBoolean:            func() DlmsData { return NewBooleanData(false) },
	TagInteger:            func() DlmsData { return NewIntegerData(0) },
	TagUnsigned:           func() DlmsData { return NewUnsignedIntegerData(0) },
	TagLong:               func() DlmsData { return NewLongData(0) },
	TagLongUnsigned:       func() DlmsData { return NewUnsignedLongData(0) },
	TagDoubleLong:         func() DlmsData { return NewDoubleLongData(0) },
	TagDoubleLongUnsigned: func() DlmsData { return NewDoubleLongUnsignedData(0) },
	TagOctetString:        func() DlmsData { return NewOctetStringData(nil) },
	TagVisibleString:      func() DlmsData { return NewVisibleStringData("") },
}

// GetDataClass returns a factory function for the given tag
func (f *DlmsDataFactory) GetDataClass(tag DlmsDataTag) (func() DlmsData, error) {
	factory, ok := dataFactoryMap[tag]
	if !ok {
		return nil, fmt.Errorf("unknown DLMS data tag: %d", tag)
	}
	return factory, nil
}

// NewDlmsDataFactory creates a new DlmsDataFactory
func NewDlmsDataFactory() *DlmsDataFactory {
	return &DlmsDataFactory{}
}

