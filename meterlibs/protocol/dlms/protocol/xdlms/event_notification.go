package xdlms

import (
	"fmt"
	"time"

	"github.com/yimiliya/idis/meterlibs/protocol/dlms/dlmsdata"
)

// EventNotification represents an Event Notification APDU
const EventNotificationTag = 15

type EventNotification struct {
	*BaseXDlmsApdu
	LongInvokeIDAndPriority *LongInvokeIdAndPriority
	DateTime                *time.Time
	Body                    []byte
}

// NewEventNotification creates a new EventNotification
func NewEventNotification(
	longInvokeIDAndPriority *LongInvokeIdAndPriority,
	dateTime *time.Time,
	body []byte,
) *EventNotification {
	return &EventNotification{
		BaseXDlmsApdu: &BaseXDlmsApdu{
			Tag: EventNotificationTag,
		},
		LongInvokeIDAndPriority: longInvokeIDAndPriority,
		DateTime:                dateTime,
		Body:                    body,
	}
}

// FromBytes creates EventNotification from bytes
func (e *EventNotification) FromBytes(sourceBytes []byte) (*EventNotification, error) {
	if len(sourceBytes) < 5 {
		return nil, fmt.Errorf("insufficient data for EventNotification, need at least 5 bytes")
	}

	data := make([]byte, len(sourceBytes))
	copy(data, sourceBytes)

	tag := data[0]
	if tag != EventNotificationTag {
		return nil, fmt.Errorf("data is not an EventNotification APDU, expected tag=%d but got %d", EventNotificationTag, tag)
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
		enDateTimeData := data[:12]
		parsedDateTime, _, err := dlmsdata.DateTimeFromBytes(enDateTimeData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse datetime: %w", err)
		}
		dateTime = &parsedDateTime
		data = data[12:]
	}

	return NewEventNotification(longInvokeID, dateTime, data), nil
}

// ToBytes converts EventNotification to bytes
func (e *EventNotification) ToBytes() ([]byte, error) {
	result := []byte{EventNotificationTag}
	result = append(result, e.LongInvokeIDAndPriority.ToBytes()...)

	if e.DateTime != nil {
		result = append(result, 0x01)
		// Use default clock status (all false)
		clockStatus := dlmsdata.NewClockStatus(false, false, false, false, false)
		dateTimeBytes := dlmsdata.DateTimeToBytes(*e.DateTime, clockStatus)
		result = append(result, dateTimeBytes...)
	} else {
		result = append(result, 0x00)
	}

	result = append(result, e.Body...)
	return result, nil
}
