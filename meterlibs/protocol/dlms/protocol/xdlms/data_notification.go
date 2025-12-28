package xdlms

import (
	"fmt"
	"time"

	"github.com/yimiliya/idis/meterlibs/protocol/dlms/dlmsdata"
)

// LongInvokeIdAndPriority represents a long invoke ID and priority
type LongInvokeIdAndPriority struct {
	LongInvokeID    uint32
	Prioritized     bool
	Confirmed       bool
	SelfDescriptive bool
	BreakOnError    bool
}

// NewLongInvokeIdAndPriority creates a new LongInvokeIdAndPriority
func NewLongInvokeIdAndPriority(
	longInvokeID uint32,
	prioritized bool,
	confirmed bool,
	selfDescriptive bool,
	breakOnError bool,
) *LongInvokeIdAndPriority {
	return &LongInvokeIdAndPriority{
		LongInvokeID:    longInvokeID,
		Prioritized:     prioritized,
		Confirmed:       confirmed,
		SelfDescriptive: selfDescriptive,
		BreakOnError:    breakOnError,
	}
}

// FromBytes creates LongInvokeIdAndPriority from bytes
func (l *LongInvokeIdAndPriority) FromBytes(bytesData []byte) (*LongInvokeIdAndPriority, error) {
	if len(bytesData) != 4 {
		return nil, fmt.Errorf("LongInvokeIdAndPriority is 4 bytes long, received: %d", len(bytesData))
	}

	longInvokeID := uint32(bytesData[1])<<16 | uint32(bytesData[2])<<8 | uint32(bytesData[3])
	statusByte := bytesData[0]
	prioritized := (statusByte & 0b10000000) != 0
	confirmed := (statusByte & 0b01000000) != 0
	breakOnError := (statusByte & 0b00100000) != 0
	selfDescriptive := (statusByte & 0b00010000) != 0

	return NewLongInvokeIdAndPriority(longInvokeID, prioritized, confirmed, selfDescriptive, breakOnError), nil
}

// ToBytes converts LongInvokeIdAndPriority to bytes
func (l *LongInvokeIdAndPriority) ToBytes() []byte {
	status := byte(0)
	if l.Prioritized {
		status |= 0b10000000
	}
	if l.Confirmed {
		status |= 0b01000000
	}
	if l.BreakOnError {
		status |= 0b00100000
	}
	if l.SelfDescriptive {
		status |= 0b00010000
	}

	result := []byte{status}
	result = append(result, byte(l.LongInvokeID>>16), byte(l.LongInvokeID>>8), byte(l.LongInvokeID))
	return result
}

// DataNotification represents a Data Notification APDU
const DataNotificationTag = 15

type DataNotification struct {
	*BaseXDlmsApdu
	LongInvokeIDAndPriority *LongInvokeIdAndPriority
	DateTime                *time.Time
	Body                    []byte
}

// NewDataNotification creates a new DataNotification
func NewDataNotification(
	longInvokeIDAndPriority *LongInvokeIdAndPriority,
	dateTime *time.Time,
	body []byte,
) *DataNotification {
	return &DataNotification{
		BaseXDlmsApdu: &BaseXDlmsApdu{
			Tag: DataNotificationTag,
		},
		LongInvokeIDAndPriority: longInvokeIDAndPriority,
		DateTime:                dateTime,
		Body:                    body,
	}
}

// FromBytes creates DataNotification from bytes
func (d *DataNotification) FromBytes(sourceBytes []byte) (*DataNotification, error) {
	if len(sourceBytes) < 5 {
		return nil, fmt.Errorf("insufficient data for DataNotification, need at least 5 bytes")
	}

	data := make([]byte, len(sourceBytes))
	copy(data, sourceBytes)

	tag := data[0]
	if tag != DataNotificationTag {
		return nil, fmt.Errorf("data is not a DataNotification APDU, expected tag=%d but got %d", DataNotificationTag, tag)
	}

	data = data[1:]

	longInvokeIDData := data[:4]
	longInvokeID, err := (&LongInvokeIdAndPriority{}).FromBytes(longInvokeIDData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse LongInvokeIdAndPriority: %w", err)
	}
	data = data[4:]

	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for has_datetime flag")
	}

	hasDateTime := data[0] != 0
	data = data[1:]

	var dateTime *time.Time
	if hasDateTime {
		if len(data) < 12 {
			return nil, fmt.Errorf("insufficient data for datetime, need 12 bytes")
		}
		dnDateTimeData := data[:12]
		parsedDateTime, _, err := dlmsdata.DateTimeFromBytes(dnDateTimeData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse datetime: %w", err)
		}
		dateTime = &parsedDateTime
		data = data[12:]
	}

	return NewDataNotification(longInvokeID, dateTime, data), nil
}

// ToBytes converts DataNotification to bytes
func (d *DataNotification) ToBytes() ([]byte, error) { {
	result := []byte{DataNotificationTag}
	result = append(result, d.LongInvokeIDAndPriority.ToBytes()...)

	if d.DateTime != nil {
		result = append(result, 0x01)
		// Use default clock status (all false)
		clockStatus := dlmsdata.NewClockStatus(false, false, false, false, false)
		dateTimeBytes := dlmsdata.DateTimeToBytes(*d.DateTime, clockStatus)
		result = append(result, dateTimeBytes...)
	} else {
		result = append(result, 0x00)
	}

	result = append(result, d.Body...)
	return result, nil
}

