package acse

import (
	"fmt"

	"github.com/yimiliya/idis/meterlibs/protocol/dlms/enumerations"
	"github.com/yimiliya/idis/meterlibs/protocol/dlms/encoding"
)

// ReleaseRequest represents an RLRQ (Release Request)
// When closing down an Application Association a ReleaseRequest is sent
const RLRQTag = 98 // Application 2

type ReleaseRequest struct {
	Reason            *enumerations.ReleaseRequestReason
	UserInformation   *UserInformation
}

// NewReleaseRequest creates a new ReleaseRequest
func NewReleaseRequest(
	reason *enumerations.ReleaseRequestReason,
	userInformation *UserInformation,
) *ReleaseRequest {
	return &ReleaseRequest{
		Reason:          reason,
		UserInformation: userInformation,
	}
}

// FromBytes creates ReleaseRequest from bytes
func (r *ReleaseRequest) FromBytes(sourceBytes []byte) (*ReleaseRequest, error) {
	if len(sourceBytes) == 0 {
		return nil, fmt.Errorf("insufficient data for RLRQ tag")
	}

	rlrqData := make([]byte, len(sourceBytes))
	copy(rlrqData, sourceBytes)

	rlrqTag := rlrqData[0]
	if rlrqTag != RLRQTag {
		return nil, fmt.Errorf("bytes are not an RLRQ APDU, tag is not %d, got %d", RLRQTag, rlrqTag)
	}

	if len(rlrqData) < 2 {
		return nil, fmt.Errorf("insufficient data for RLRQ length")
	}

	rlrqLength := int(rlrqData[1])
	rlrqData = rlrqData[2:]

	if len(rlrqData) != rlrqLength {
		return nil, fmt.Errorf("the APDU data length does not correspond to length byte, expected %d, got %d", rlrqLength, len(rlrqData))
	}

	// Parse tags
	objectDict := make(map[string]interface{})
	ber := encoding.NewBER()

	for len(rlrqData) > 0 {
		if len(rlrqData) < 2 {
			return nil, fmt.Errorf("insufficient data for tag and length")
		}

		objectTag := rlrqData[0]
		objectLength := int(rlrqData[1])
		rlrqData = rlrqData[2:]

		if len(rlrqData) < objectLength {
			return nil, fmt.Errorf("insufficient data for object, need %d bytes, got %d", objectLength, len(rlrqData))
		}

		objectData := rlrqData[:objectLength]
		rlrqData = rlrqData[objectLength:]

		var objectName string
		var parsedData interface{}
		var err error

		switch objectTag {
		case 0x80: // reason
			objectName = "reason"
			if len(objectData) > 0 {
				// Decode BER encoded integer
				tag, length, data, err := ber.Decode(objectData, 1)
				if err != nil {
					return nil, fmt.Errorf("failed to decode reason: %w", err)
				}
				if !bytesEqual(tag, []byte{2}) { // Integer tag
					return nil, fmt.Errorf("reason is not an integer")
				}
				if len(data) != int(length) || len(data) == 0 {
					return nil, fmt.Errorf("invalid reason data length")
				}
				reason := enumerations.ReleaseRequestReason(data[0])
				parsedData = &reason
			} else {
				parsedData = nil
			}
		case 0xBE: // user_information
			objectName = "user_information"
			userInfo := &UserInformation{}
			parsedData, err = userInfo.FromBytes(objectData)
			if err != nil {
				return nil, fmt.Errorf("failed to parse user_information: %w", err)
			}
		default:
			return nil, fmt.Errorf("could not find object with tag 0x%02x in RLRQ definition", objectTag)
		}

		objectDict[objectName] = parsedData
	}

	reason, _ := objectDict["reason"].(*enumerations.ReleaseRequestReason)
	userInformation, _ := objectDict["user_information"].(*UserInformation)

	return NewReleaseRequest(reason, userInformation), nil
}

// ToBytes converts ReleaseRequest to bytes
func (r *ReleaseRequest) ToBytes() ([]byte, error) {
	ber := encoding.NewBER()
	rlrqData := make([]byte, 0)

	if r.Reason != nil {
		reasonBytes, err := ber.Encode(0x80, []byte{byte(*r.Reason)})
		if err != nil {
			return nil, fmt.Errorf("failed to encode reason: %w", err)
		}
		rlrqData = append(rlrqData, reasonBytes...)
	}

	if r.UserInformation != nil {
		userInfoBytes, err := r.UserInformation.ToBytes()
		if err != nil {
			return nil, fmt.Errorf("failed to encode user_information: %w", err)
		}
		encoded, err := ber.Encode(0xBE, userInfoBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to BER encode user_information: %w", err)
		}
		rlrqData = append(rlrqData, encoded...)
	}

	return ber.Encode(RLRQTag, rlrqData)
}
