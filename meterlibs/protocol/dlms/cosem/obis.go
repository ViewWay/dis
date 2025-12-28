package cosem

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	sixPartRegex  = regexp.MustCompile(`^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$`)
	fivePartRegex = regexp.MustCompile(`^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$`)
)

// Obis represents an OBject Identification System code
// OBIS defines codes for identification of commonly used
// data items in metering equipment.
type Obis struct {
	A int
	B int
	C int
	D int
	E int
	F int
}

// NewObis creates a new Obis from individual components
func NewObis(a, b, c, d, e int, f int) (*Obis, error) {
	if err := validateOBISComponent(a); err != nil {
		return nil, fmt.Errorf("invalid A component: %w", err)
	}
	if err := validateOBISComponent(b); err != nil {
		return nil, fmt.Errorf("invalid B component: %w", err)
	}
	if err := validateOBISComponent(c); err != nil {
		return nil, fmt.Errorf("invalid C component: %w", err)
	}
	if err := validateOBISComponent(d); err != nil {
		return nil, fmt.Errorf("invalid D component: %w", err)
	}
	if err := validateOBISComponent(e); err != nil {
		return nil, fmt.Errorf("invalid E component: %w", err)
	}
	if err := validateOBISComponent(f); err != nil {
		return nil, fmt.Errorf("invalid F component: %w", err)
	}

	return &Obis{
		A: a,
		B: b,
		C: c,
		D: d,
		E: e,
		F: f,
	}, nil
}

// validateOBISComponent validates that a component is between 0 and 255
func validateOBISComponent(value int) error {
	if value < 0 || value > 255 {
		return fmt.Errorf("OBIS component must be between 0 and 255, got %d", value)
	}
	return nil
}

// FromBytes creates an Obis from bytes
func FromBytes(sourceBytes []byte) (*Obis, error) {
	if len(sourceBytes) != 6 {
		return nil, fmt.Errorf("not enough data to parse OBIS. Need 6 bytes but got %d", len(sourceBytes))
	}

	return NewObis(
		int(sourceBytes[0]),
		int(sourceBytes[1]),
		int(sourceBytes[2]),
		int(sourceBytes[3]),
		int(sourceBytes[4]),
		int(sourceBytes[5]),
	)
}

// FromString parses a string as an OBIS code
// Will accept with both the optional 255 at the end and not. Any separator is allowed.
func FromString(obisString string) (*Obis, error) {
	// Try six part match first
	if matches := sixPartRegex.FindStringSubmatch(obisString); matches != nil {
		parts := matches[1:]
		values := make([]int, 6)
		for i, part := range parts {
			val, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("failed to parse component %d: %w", i, err)
			}
			values[i] = val
		}
		return NewObis(values[0], values[1], values[2], values[3], values[4], values[5])
	}

	// Try five part match
	if matches := fivePartRegex.FindStringSubmatch(obisString); matches != nil {
		parts := matches[1:]
		values := make([]int, 5)
		for i, part := range parts {
			val, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("failed to parse component %d: %w", i, err)
			}
			values[i] = val
		}
		// Default F to 255 if not provided
		return NewObis(values[0], values[1], values[2], values[3], values[4], 255)
	}

	return nil, fmt.Errorf("%s is not a parsable OBIS string", obisString)
}

// ToString converts Obis to string representation
// separator is optional, default format is "A-B:C.D.E.F"
func (o *Obis) ToString(separator string) string {
	if separator != "" {
		return fmt.Sprintf("%d%s%d%s%d%s%d%s%d%s%d",
			o.A, separator,
			o.B, separator,
			o.C, separator,
			o.D, separator,
			o.E, separator,
			o.F)
	}
	return fmt.Sprintf("%d-%d:%d.%d.%d.%d", o.A, o.B, o.C, o.D, o.E, o.F)
}

// String implements fmt.Stringer interface
func (o *Obis) String() string {
	return o.ToString("")
}

// ToBytes converts Obis to bytes
func (o *Obis) ToBytes() []byte {
	return []byte{
		byte(o.A),
		byte(o.B),
		byte(o.C),
		byte(o.D),
		byte(o.E),
		byte(o.F),
	}
}

