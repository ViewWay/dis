package acse

import (
	"fmt"

	"github.com/yimiliya/idis/meterlibs/protocol/dlms/encoding"
	"github.com/yimiliya/idis/meterlibs/protocol/dlms/protocol/xdlms"
)

// UserInformation holds InitiateRequests for AARQ and InitiateResponse for AARE
type UserInformation struct {
	Tag     []byte
	Content interface{} // Can be InitiateRequest, InitiateResponse, ConfirmedServiceError, etc.
}

// NewUserInformation creates a new UserInformation
func NewUserInformation(content interface{}) *UserInformation {
	return &UserInformation{
		Tag:     []byte{0x04}, // encoded as an octetstring
		Content: content,
	}
}

// FromBytes creates UserInformation from bytes
func (u *UserInformation) FromBytes(data []byte) (*UserInformation, error) {
	ber := encoding.NewBER()
	tag, length, berData, err := ber.Decode(data, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to decode BER: %w", err)
	}

	if !bytesEqual(tag, u.Tag) {
		return nil, fmt.Errorf("the tag for UserInformation data should be 0x04, not %v", tag)
	}

	if len(berData) == 0 {
		return nil, fmt.Errorf("insufficient data for user information content")
	}

	var content interface{}
	switch berData[0] {
	case 1:
		initReq := &xdlms.InitiateRequest{}
		parsedReq, err := initReq.FromBytes(berData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse InitiateRequest: %w", err)
		}
		content = parsedReq
	case 8:
		initResp := &xdlms.InitiateResponse{}
		parsedResp, err := initResp.FromBytes(berData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse InitiateResponse: %w", err)
		}
		content = parsedResp
	case 14:
		// ConfirmedServiceError - TODO: implement when needed
		return nil, fmt.Errorf("ConfirmedServiceError not yet implemented")
	case 33:
		// GlobalCipherInitiateRequest - TODO: implement when needed
		return nil, fmt.Errorf("GlobalCipherInitiateRequest not yet implemented")
	case 40:
		// GlobalCipherInitiateResponse - TODO: implement when needed
		return nil, fmt.Errorf("GlobalCipherInitiateResponse not yet implemented")
	default:
		return nil, fmt.Errorf("not able to find a proper data tag in UserInformation, got %d", berData[0])
	}

	return NewUserInformation(content), nil
}

// ToBytes converts UserInformation to bytes
func (u *UserInformation) ToBytes() ([]byte, error) {
	ber := encoding.NewBER()
	var contentBytes []byte
	var err error

	switch c := u.Content.(type) {
	case *xdlms.InitiateRequest:
		contentBytes, err = c.ToBytes()
		if err != nil {
			return nil, fmt.Errorf("failed to encode InitiateRequest: %w", err)
		}
	case *xdlms.InitiateResponse:
		contentBytes, err = c.ToBytes()
		if err != nil {
			return nil, fmt.Errorf("failed to encode InitiateResponse: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported content type: %T", u.Content)
	}

	return ber.Encode(u.Tag, contentBytes)
}

