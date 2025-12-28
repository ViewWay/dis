package xdlms

import (
	"encoding/binary"
	"fmt"

	"github.com/yimiliya/idis/meterlibs/protocol/dlms/enumerations"
)

// ExceptionResponse represents an Exception Response APDU
const ExceptionResponseTag = 216

type ExceptionResponse struct {
	*BaseXDlmsApdu
	StateError            enumerations.StateException
	ServiceError          enumerations.ServiceException
	InvocationCounterData *uint32
}

// NewExceptionResponse creates a new ExceptionResponse
func NewExceptionResponse(
	stateError enumerations.StateException,
	serviceError enumerations.ServiceException,
	invocationCounterData *uint32,
) *ExceptionResponse {
	return &ExceptionResponse{
		BaseXDlmsApdu: &BaseXDlmsApdu{
			Tag: ExceptionResponseTag,
		},
		StateError:            stateError,
		ServiceError:          serviceError,
		InvocationCounterData: invocationCounterData,
	}
}

// FromBytes creates ExceptionResponse from bytes
func (e *ExceptionResponse) FromBytes(sourceBytes []byte) (*ExceptionResponse, error) {
	if len(sourceBytes) < 3 {
		return nil, fmt.Errorf("insufficient data for ExceptionResponse, need at least 3 bytes")
	}

	data := make([]byte, len(sourceBytes))
	copy(data, sourceBytes)

	tag := data[0]
	if tag != ExceptionResponseTag {
		return nil, fmt.Errorf("tag for ExceptionResponse is not %d, got %d instead", ExceptionResponseTag, tag)
	}

	stateError := enumerations.StateException(data[1])
	serviceError := enumerations.ServiceException(data[2])
	data = data[3:]

	var invocationCounterData *uint32
	if serviceError == enumerations.ServiceExceptionInvocationCounterError {
		if len(data) < 4 {
			return nil, fmt.Errorf("insufficient data for invocation counter, need 4 bytes")
		}
		counter := binary.BigEndian.Uint32(data[:4])
		invocationCounterData = &counter
	}

	return NewExceptionResponse(stateError, serviceError, invocationCounterData), nil
}

// ToBytes converts ExceptionResponse to bytes
func (e *ExceptionResponse) ToBytes() ([]byte, error) {
	result := []byte{ExceptionResponseTag, byte(e.StateError), byte(e.ServiceError)}

	if e.InvocationCounterData != nil {
		counterBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(counterBytes, *e.InvocationCounterData)
		result = append(result, counterBytes...)
	}

	return result, nil
}

