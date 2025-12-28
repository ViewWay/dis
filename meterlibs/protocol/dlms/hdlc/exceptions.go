package hdlc

import "fmt"

// LocalProtocolError represents an error in HDLC Protocol
type LocalProtocolError struct {
	Message string
}

func (e *LocalProtocolError) Error() string {
	return fmt.Sprintf("HDLC protocol error: %s", e.Message)
}

// NewLocalProtocolError creates a new LocalProtocolError
func NewLocalProtocolError(message string) *LocalProtocolError {
	return &LocalProtocolError{Message: message}
}

// HdlcException is the base class for HDLC protocol parts
type HdlcException struct {
	Message string
}

func (e *HdlcException) Error() string {
	return fmt.Sprintf("HDLC exception: %s", e.Message)
}

// NewHdlcException creates a new HdlcException
func NewHdlcException(message string) *HdlcException {
	return &HdlcException{Message: message}
}

// HdlcParsingError represents an error that occurred when parsing bytes into HDLC object
type HdlcParsingError struct {
	*HdlcException
}

// NewHdlcParsingError creates a new HdlcParsingError
func NewHdlcParsingError(message string) *HdlcParsingError {
	return &HdlcParsingError{
		HdlcException: NewHdlcException(message),
	}
}

// MissingHdlcFlags represents an error when frame is not enclosed by HDLC flags
type MissingHdlcFlags struct {
	*HdlcParsingError
}

// NewMissingHdlcFlags creates a new MissingHdlcFlags
func NewMissingHdlcFlags() *MissingHdlcFlags {
	return &MissingHdlcFlags{
		HdlcParsingError: NewHdlcParsingError("frame is not enclosed by HDLC flags"),
	}
}

