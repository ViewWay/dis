package dlmsdata

import (
	"encoding/binary"
	"fmt"
	"time"
)

// ClockStatus represents the clock status byte
type ClockStatus struct {
	Invalid              bool
	Doubtful             bool
	DifferentBase        bool
	InvalidStatus        bool
	DaylightSavingActive bool
}

// NewClockStatus creates a new ClockStatus
func NewClockStatus(invalid, doubtful, differentBase, invalidStatus, daylightSavingActive bool) *ClockStatus {
	return &ClockStatus{
		Invalid:              invalid,
		Doubtful:             doubtful,
		DifferentBase:        differentBase,
		InvalidStatus:        invalidStatus,
		DaylightSavingActive: daylightSavingActive,
	}
}

// FromBytes creates ClockStatus from bytes
func (c *ClockStatus) FromBytes(data []byte) (*ClockStatus, error) {
	if len(data) != 1 {
		return nil, fmt.Errorf("ClockStatus is 1 byte, got %d", len(data))
	}
	value := data[0]
	return &ClockStatus{
		Invalid:              (value & 0b00000001) != 0,
		Doubtful:             (value & 0b00000010) != 0,
		DifferentBase:        (value & 0b00000100) != 0,
		InvalidStatus:        (value & 0b00001000) != 0,
		DaylightSavingActive: (value & 0b10000000) != 0,
	}, nil
}

// ToBytes converts ClockStatus to bytes
func (c *ClockStatus) ToBytes() []byte {
	var value byte
	if c.Invalid {
		value |= 0b00000001
	}
	if c.Doubtful {
		value |= 0b00000010
	}
	if c.DifferentBase {
		value |= 0b00000100
	}
	if c.InvalidStatus {
		value |= 0b00001000
	}
	if c.DaylightSavingActive {
		value |= 0b10000000
	}
	return []byte{value}
}

// DateFromBytes parses a date from 5 bytes
// [year highbyte, year lowbyte, month, day of month, day of week]
func DateFromBytes(data []byte) (time.Time, error) {
	if len(data) != 5 {
		return time.Time{}, fmt.Errorf("date is represented by 5 bytes, but got %d", len(data))
	}
	
	year := binary.BigEndian.Uint16(data[:2])
	month := data[2]
	day := data[3]
	// dayOfWeek := data[4] // not used for now
	
	// Handle special cases
	if year == 0xFFFF {
		return time.Time{}, fmt.Errorf("year not specified (0xFFFF)")
	}
	if month == 0xFF {
		return time.Time{}, fmt.Errorf("month not specified (0xFF)")
	}
	if day == 0xFF {
		return time.Time{}, fmt.Errorf("day not specified (0xFF)")
	}
	
	return time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC), nil
}

// TimeFromBytes parses a time from 4 bytes
// [hour, minute, second, hundredths]
func TimeFromBytes(data []byte) (time.Time, error) {
	if len(data) != 4 {
		return time.Time{}, fmt.Errorf("time is represented by 4 bytes, but got %d", len(data))
	}
	
	hour := data[0]
	minute := data[1]
	second := data[2]
	hundredths := data[3]
	
	// Handle special cases (0xFF = not specified, use 0)
	if hour == 0xFF {
		hour = 0
	}
	if minute == 0xFF {
		minute = 0
	}
	if second == 0xFF {
		second = 0
	}
	if hundredths == 0xFF {
		hundredths = 0
	}
	
	// Use a reference date (2000-01-01) for time-only values
	refDate := time.Date(2000, 1, 1, int(hour), int(minute), int(second), int(hundredths)*10000, time.UTC)
	return refDate, nil
}

// DateTimeFromBytes parses a datetime from 12 bytes
// [date[5 bytes], time[4 bytes], deviation[2 bytes], clock_status[1 byte]]
func DateTimeFromBytes(data []byte) (time.Time, *ClockStatus, error) {
	if len(data) != 12 {
		return time.Time{}, nil, fmt.Errorf("datetime is represented by 12 bytes, but got %d", len(data))
	}
	
	dateData := data[:5]
	timeData := data[5:9]
	deviationData := data[9:11]
	statusData := data[11:12]
	
	// Parse date
	year := binary.BigEndian.Uint16(dateData[:2])
	month := dateData[2]
	day := dateData[3]
	
	if year == 0xFFFF || month == 0xFF || day == 0xFF {
		return time.Time{}, nil, fmt.Errorf("date contains unspecified values")
	}
	
	// Parse time
	hour := timeData[0]
	minute := timeData[1]
	second := timeData[2]
	hundredths := timeData[3]
	
	if hour == 0xFF {
		hour = 0
	}
	if minute == 0xFF {
		minute = 0
	}
	if second == 0xFF {
		second = 0
	}
	if hundredths == 0xFF {
		hundredths = 0
	}
	
	// Parse deviation (timezone offset in minutes, signed)
	deviationUint := binary.BigEndian.Uint16(deviationData)
	var tz *time.Location
	if deviationUint == 0x8000 {
		// Not specified, use UTC
		tz = time.UTC
	} else {
		// DLMS uses deviation from local time to UTC (negated)
		// Convert uint16 to int16, handling the sign bit
		deviation := int16(deviationUint)
		offsetSeconds := -int(deviation) * 60
		tz = time.FixedZone("", offsetSeconds)
	}
	
	// Parse clock status
	var status *ClockStatus
	if len(statusData) > 0 {
		status, _ = (&ClockStatus{}).FromBytes(statusData)
	}
	
	dt := time.Date(
		int(year),
		time.Month(month),
		int(day),
		int(hour),
		int(minute),
		int(second),
		int(hundredths)*10000,
		tz,
	)
	
	return dt, status, nil
}

// DateToBytes converts a date to 5 bytes
func DateToBytes(d time.Time) []byte {
	year := uint16(d.Year())
	month := byte(d.Month())
	day := byte(d.Day())
	dayOfWeekUnspecified := byte(0xFF)
	
	result := make([]byte, 5)
	binary.BigEndian.PutUint16(result[:2], year)
	result[2] = month
	result[3] = day
	result[4] = dayOfWeekUnspecified
	
	return result
}

// TimeToBytes converts a time to 4 bytes
func TimeToBytes(t time.Time) []byte {
	result := make([]byte, 4)
	result[0] = byte(t.Hour())
	result[1] = byte(t.Minute())
	result[2] = byte(t.Second())
	result[3] = byte(t.Nanosecond() / 10000000) // hundredths
	return result
}

// DateTimeToBytes converts a datetime to 12 bytes
func DateTimeToBytes(dt time.Time, clockStatus *ClockStatus) []byte {
	dateBytes := DateToBytes(dt)
	timeBytes := TimeToBytes(dt)
	
	// Calculate timezone deviation
	var deviationBytes []byte
	if dt.Location() == nil || dt.Location() == time.UTC {
		deviationBytes = []byte{0x80, 0x00} // not specified
	} else {
		// DLMS uses deviation from local time to UTC (negated)
		_, offset := dt.Zone()
		deviationMinutes := int16(-offset / 60)
		deviationBytes = make([]byte, 2)
		binary.BigEndian.PutUint16(deviationBytes, uint16(deviationMinutes))
	}
	
	// Clock status
	var statusBytes []byte
	if clockStatus != nil {
		statusBytes = clockStatus.ToBytes()
	} else {
		statusBytes = NewClockStatus(false, false, false, false, false).ToBytes()
	}
	
	result := make([]byte, 0, 12)
	result = append(result, dateBytes...)
	result = append(result, timeBytes...)
	result = append(result, deviationBytes...)
	result = append(result, statusBytes...)
	
	return result
}

