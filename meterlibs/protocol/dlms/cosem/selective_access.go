package cosem

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/yimiliya/idis/meterlibs/protocol/dlms/dlmsdata"
	"github.com/yimiliya/idis/meterlibs/protocol/dlms/enumerations"
)

// AccessDescriptorType represents the type of access descriptor
type AccessDescriptorType uint8

const (
	AccessDescriptorTypeRange AccessDescriptorType = 1
	AccessDescriptorTypeEntry AccessDescriptorType = 2
)

// RangeDescriptor can be used to read buffers of Profile Generic.
// Only buffer element that corresponds to the descriptor shall be returned in a get request.
type RangeDescriptor struct {
	RestrictingObject *CaptureObject
	FromValue         time.Time
	ToValue           time.Time
	SelectedValues    []*CaptureObject // nil means all columns
}

// NewRangeDescriptor creates a new RangeDescriptor
func NewRangeDescriptor(
	restrictingObject *CaptureObject,
	fromValue time.Time,
	toValue time.Time,
	selectedValues []*CaptureObject,
) *RangeDescriptor {
	return &RangeDescriptor{
		RestrictingObject: restrictingObject,
		FromValue:         fromValue,
		ToValue:           toValue,
		SelectedValues:    selectedValues,
	}
}

// ToBytes converts RangeDescriptor to bytes
func (r *RangeDescriptor) ToBytes() []byte {
	result := []byte{byte(AccessDescriptorTypeRange)}
	
	// Structure of 4 elements
	result = append(result, 0x02, 0x04)
	
	// Restricting object
	result = append(result, r.RestrictingObject.ToBytes()...)
	
	// From value (datetime as OctetString)
	// TODO: Implement datetime_to_bytes when time package is ready
	fromBytes := datetimeToBytes(r.FromValue)
	result = append(result, 0x09) // OctetString tag
	result = append(result, byte(len(fromBytes)))
	result = append(result, fromBytes...)
	
	// To value (datetime as OctetString)
	toBytes := datetimeToBytes(r.ToValue)
	result = append(result, 0x09) // OctetString tag
	result = append(result, byte(len(toBytes)))
	result = append(result, toBytes...)
	
	// Selected values
	if r.SelectedValues == nil || len(r.SelectedValues) == 0 {
		// Empty array means all columns
		result = append(result, 0x01, 0x00) // Array tag + length 0
	} else {
		// TODO: Implement selected values
		panic("selected values not yet implemented")
	}
	
	return result
}

// FromBytes creates RangeDescriptor from bytes and returns the number of bytes consumed
func (r *RangeDescriptor) FromBytes(sourceBytes []byte) (*RangeDescriptor, int, error) {
	if len(sourceBytes) < 1 {
		return nil, 0, fmt.Errorf("insufficient data for RangeDescriptor")
	}
	
	// Check access descriptor type
	if sourceBytes[0] != byte(AccessDescriptorTypeRange) {
		return nil, 0, fmt.Errorf("access descriptor type %d is not valid for RangeDescriptor. It should be %d", sourceBytes[0], AccessDescriptorTypeRange)
	}
	
	// Check structure tag and length (should be 0x02 0x04 for structure of 4 elements)
	if len(sourceBytes) < 3 || sourceBytes[1] != 0x02 || sourceBytes[2] != 0x04 {
		return nil, 0, fmt.Errorf("invalid structure tag or length for RangeDescriptor")
	}
	
	offset := 3
	
	// Parse restricting object (CaptureObject - structure of 4 elements)
	// CaptureObject structure: interface (UnsignedLong), instance (OctetString), attribute (Integer), data_index (UnsignedLong)
	if len(sourceBytes) < offset+2 || sourceBytes[offset] != 0x02 || sourceBytes[offset+1] != 0x04 {
		return nil, 0, fmt.Errorf("invalid restricting object structure")
	}
	offset += 2
	
	// Parse interface (UnsignedLong - 2 bytes)
	// Need to check offset+4 because we access sourceBytes[offset+3] after offset += 2
	if len(sourceBytes) < offset+4 || sourceBytes[offset] != 0x11 || sourceBytes[offset+1] != 0x02 {
		return nil, 0, fmt.Errorf("invalid interface tag or length")
	}
	offset += 2
	interfaceValue := uint16(sourceBytes[offset])<<8 | uint16(sourceBytes[offset+1])
	offset += 2
	
	// Parse instance (OctetString - 6 bytes OBIS)
	if len(sourceBytes) < offset+2 || sourceBytes[offset] != 0x09 || sourceBytes[offset+1] != 0x06 {
		return nil, 0, fmt.Errorf("invalid instance tag or length")
	}
	offset += 2
	obisBytes := sourceBytes[offset : offset+6]
	offset += 6
	obis, err := FromBytes(obisBytes)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse OBIS: %w", err)
	}
	
	// Parse attribute (Integer - 1 byte)
	if len(sourceBytes) < offset+2 || sourceBytes[offset] != 0x0F || sourceBytes[offset+1] != 0x01 {
		return nil, 0, fmt.Errorf("invalid attribute tag or length")
	}
	offset += 2
	attribute := sourceBytes[offset]
	offset++
	
	// Parse data_index (UnsignedLong - 2 bytes, tag 0x12)
	if len(sourceBytes) < offset+2 || sourceBytes[offset] != 0x12 || sourceBytes[offset+1] != 0x02 {
		return nil, 0, fmt.Errorf("invalid data_index tag or length: expected 0x12 0x02 (UnsignedLong), got 0x%02x 0x%02x", sourceBytes[offset], sourceBytes[offset+1])
	}
	offset += 2
	dataIndex := uint16(sourceBytes[offset])<<8 | uint16(sourceBytes[offset+1])
	offset += 2
	
	cosemAttribute := NewCosemAttribute(enumerations.CosemInterface(interfaceValue), obis, attribute)
	restrictingObject := NewCaptureObject(cosemAttribute, dataIndex)
	
	// Parse from_value (OctetString containing datetime)
	if len(sourceBytes) < offset+2 || sourceBytes[offset] != 0x09 {
		return nil, 0, fmt.Errorf("invalid from_value tag: expected 0x09 (OctetString), got 0x%02x", sourceBytes[offset])
	}
	offset++
	fromValueLength := int(sourceBytes[offset])
	offset++
	if len(sourceBytes) < offset+fromValueLength {
		return nil, 0, fmt.Errorf("insufficient data for from_value")
	}
	fromValueBytes := sourceBytes[offset : offset+fromValueLength]
	offset += fromValueLength
	fromValue, _, err := dlmsdata.DateTimeFromBytes(fromValueBytes)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse from_value datetime: %w", err)
	}
	
	// Parse to_value (OctetString containing datetime)
	if len(sourceBytes) < offset+2 || sourceBytes[offset] != 0x09 {
		return nil, 0, fmt.Errorf("invalid to_value tag: expected 0x09 (OctetString), got 0x%02x", sourceBytes[offset])
	}
	offset++
	toValueLength := int(sourceBytes[offset])
	offset++
	if len(sourceBytes) < offset+toValueLength {
		return nil, 0, fmt.Errorf("insufficient data for to_value")
	}
	toValueBytes := sourceBytes[offset : offset+toValueLength]
	offset += toValueLength
	toValue, _, err := dlmsdata.DateTimeFromBytes(toValueBytes)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse to_value datetime: %w", err)
	}
	
	// Parse selected_values (Array - can be empty)
	var selectedValues []*CaptureObject
	if len(sourceBytes) > offset {
		if sourceBytes[offset] == 0x01 { // Array tag
			offset++
			if len(sourceBytes) <= offset {
				return nil, 0, fmt.Errorf("insufficient data for selected_values array length")
			}
			arrayLength := int(sourceBytes[offset])
			offset++
			if arrayLength > 0 {
				// TODO: Parse selected values when needed
				// For now, we'll skip them as they're not commonly used
				return nil, 0, fmt.Errorf("selected values parsing not yet implemented")
			}
			// Empty array (arrayLength == 0) means all columns, which is valid
			selectedValues = nil
		}
	}
	
	return NewRangeDescriptor(restrictingObject, fromValue, toValue, selectedValues), offset, nil
}

// EntryDescriptor limits response data by entries.
// It is possible to limit the entries and also the columns returned.
// The from/to_selected_value limits the columns returned from/to_entry limits the entries.
// Numbering of selected values and entries start from 1.
// Setting to_entry=0 or to_selected_value=0 requests the highest possible value.
type EntryDescriptor struct {
	FromEntry        uint32
	ToEntry          uint32 // 0 means highest possible
	FromSelectedValue uint16 // default 1
	ToSelectedValue  uint16 // 0 means highest possible
}

// NewEntryDescriptor creates a new EntryDescriptor
func NewEntryDescriptor(
	fromEntry uint32,
	toEntry uint32,
	fromSelectedValue uint16,
	toSelectedValue uint16,
) (*EntryDescriptor, error) {
	// Note: uint32 and uint16 types already enforce their maximum values at compile time.
	// The previous checks comparing against 0xFFFFFFFF and 0xFFFF were dead code since
	// these values can never exceed their type's maximum. Validation removed as it was
	// impossible to trigger.
	
	return &EntryDescriptor{
		FromEntry:        fromEntry,
		ToEntry:          toEntry,
		FromSelectedValue: fromSelectedValue,
		ToSelectedValue:  toSelectedValue,
	}, nil
}

// ToBytes converts EntryDescriptor to bytes
func (e *EntryDescriptor) ToBytes() []byte {
	result := []byte{byte(AccessDescriptorTypeEntry)}
	
	// Structure of 4 elements
	result = append(result, 0x02, 0x04)
	
	// From entry (DoubleLongUnsigned)
	result = append(result, 0x06) // DoubleLongUnsigned tag
	fromEntryBytes := make([]byte, 4)
	fromEntryBytes[0] = byte(e.FromEntry >> 24)
	fromEntryBytes[1] = byte(e.FromEntry >> 16)
	fromEntryBytes[2] = byte(e.FromEntry >> 8)
	fromEntryBytes[3] = byte(e.FromEntry)
	result = append(result, fromEntryBytes...)
	
	// To entry (DoubleLongUnsigned)
	result = append(result, 0x06) // DoubleLongUnsigned tag
	toEntryBytes := make([]byte, 4)
	toEntryBytes[0] = byte(e.ToEntry >> 24)
	toEntryBytes[1] = byte(e.ToEntry >> 16)
	toEntryBytes[2] = byte(e.ToEntry >> 8)
	toEntryBytes[3] = byte(e.ToEntry)
	result = append(result, toEntryBytes...)
	
	// From selected value (LongUnsigned)
	result = append(result, 0x12) // LongUnsigned tag
	fromSelectedBytes := make([]byte, 2)
	fromSelectedBytes[0] = byte(e.FromSelectedValue >> 8)
	fromSelectedBytes[1] = byte(e.FromSelectedValue)
	result = append(result, fromSelectedBytes...)
	
	// To selected value (LongUnsigned)
	result = append(result, 0x12) // LongUnsigned tag
	toSelectedBytes := make([]byte, 2)
	toSelectedBytes[0] = byte(e.ToSelectedValue >> 8)
	toSelectedBytes[1] = byte(e.ToSelectedValue)
	result = append(result, toSelectedBytes...)
	
	return result
}

// FromBytes creates EntryDescriptor from bytes and returns the number of bytes consumed
func (e *EntryDescriptor) FromBytes(sourceBytes []byte) (*EntryDescriptor, int, error) {
	if len(sourceBytes) < 19 {
		return nil, 0, fmt.Errorf("insufficient data for EntryDescriptor: need at least 19 bytes, got %d", len(sourceBytes))
	}
	
	// Check access descriptor type
	if sourceBytes[0] != byte(AccessDescriptorTypeEntry) {
		return nil, 0, fmt.Errorf("access descriptor type %d is not valid for EntryDescriptor. It should be %d", sourceBytes[0], AccessDescriptorTypeEntry)
	}
	
	// Check structure tag and length
	if sourceBytes[1] != 0x02 || sourceBytes[2] != 0x04 {
		return nil, 0, fmt.Errorf("invalid structure tag or length for EntryDescriptor")
	}
	
	offset := 3
	
	// Parse from entry (DoubleLongUnsigned tag 0x06 + 4 bytes)
	if sourceBytes[offset] != 0x06 {
		return nil, 0, fmt.Errorf("invalid tag for from_entry: expected 0x06 (DoubleLongUnsigned), got 0x%02x", sourceBytes[offset])
	}
	offset++
	fromEntry := uint32(sourceBytes[offset])<<24 | uint32(sourceBytes[offset+1])<<16 | uint32(sourceBytes[offset+2])<<8 | uint32(sourceBytes[offset+3])
	offset += 4
	
	// Parse to entry (DoubleLongUnsigned tag 0x06 + 4 bytes)
	if sourceBytes[offset] != 0x06 {
		return nil, 0, fmt.Errorf("invalid tag for to_entry: expected 0x06 (DoubleLongUnsigned), got 0x%02x", sourceBytes[offset])
	}
	offset++
	toEntry := uint32(sourceBytes[offset])<<24 | uint32(sourceBytes[offset+1])<<16 | uint32(sourceBytes[offset+2])<<8 | uint32(sourceBytes[offset+3])
	offset += 4
	
	// Parse from selected value (LongUnsigned tag 0x12 + 2 bytes)
	if sourceBytes[offset] != 0x12 {
		return nil, 0, fmt.Errorf("invalid tag for from_selected_value: expected 0x12 (LongUnsigned), got 0x%02x", sourceBytes[offset])
	}
	offset++
	fromSelectedValue := uint16(sourceBytes[offset])<<8 | uint16(sourceBytes[offset+1])
	offset += 2
	
	// Parse to selected value (LongUnsigned tag 0x12 + 2 bytes)
	if sourceBytes[offset] != 0x12 {
		return nil, 0, fmt.Errorf("invalid tag for to_selected_value: expected 0x12 (LongUnsigned), got 0x%02x", sourceBytes[offset])
	}
	offset++
	toSelectedValue := uint16(sourceBytes[offset])<<8 | uint16(sourceBytes[offset+1])
	offset += 2
	
	entry, err := NewEntryDescriptor(fromEntry, toEntry, fromSelectedValue, toSelectedValue)
	if err != nil {
		return nil, 0, err
	}
	
	return entry, offset, nil
}

// AccessDescriptorFactory handles the selection of parsing the first byte
// to find what kind of access descriptor it is and returns the object.
type AccessDescriptorFactory struct{}

// FromBytes creates an access descriptor from bytes and returns the number of bytes consumed
func (f *AccessDescriptorFactory) FromBytes(sourceBytes []byte) (interface{}, int, error) {
	if len(sourceBytes) == 0 {
		return nil, 0, fmt.Errorf("insufficient data for access descriptor")
	}
	
	accessDescriptor := AccessDescriptorType(sourceBytes[0])
	
	switch accessDescriptor {
	case AccessDescriptorTypeRange:
		return (&RangeDescriptor{}).FromBytes(sourceBytes)
	case AccessDescriptorTypeEntry:
		return (&EntryDescriptor{}).FromBytes(sourceBytes)
	default:
		return nil, 0, fmt.Errorf("%d is not a valid access descriptor", accessDescriptor)
	}
}

// NewAccessDescriptorFactory creates a new AccessDescriptorFactory
func NewAccessDescriptorFactory() *AccessDescriptorFactory {
	return &AccessDescriptorFactory{}
}

// datetimeToBytes converts time.Time to DLMS datetime bytes (12 bytes)
// This is a placeholder - will be properly implemented when time package is ready
func datetimeToBytes(t time.Time) []byte {
	// DLMS datetime format: year(2), month(1), day_of_month(1), day_of_week(1),
	// hour(1), minute(1), second(1), hundredths(1), deviation(2), clock_status(1)
	result := make([]byte, 12)
	
	year := t.Year()
	result[0] = byte(year >> 8)
	result[1] = byte(year & 0xFF)
	result[2] = byte(t.Month())
	result[3] = byte(t.Day())
	result[4] = byte(t.Weekday())
	result[5] = byte(t.Hour())
	result[6] = byte(t.Minute())
	result[7] = byte(t.Second())
	result[8] = byte(t.Nanosecond() / 10000000) // hundredths
	result[9] = 0 // deviation (not set)
	result[10] = 0 // deviation (not set)
	result[11] = 0 // clock_status
	
	return result
}

