package acse

import (
	"fmt"

	"github.com/yimiliya/idis/meterlibs/protocol/dlms/enumerations"
	"github.com/yimiliya/idis/meterlibs/protocol/dlms/encoding"
)

// ReleaseResponse represents an RLRE (Release Response)
// When closing down an Application Association a ReleaseResponse is sent from the server (meter) after a ReleaseRequest
const RLRETag = 99 // Application 3

type ReleaseResponse struct {
	Reason            *enumerations.ReleaseResponseReason
	UserInformation   *UserInformation
}

// NewReleaseResponse creates a new ReleaseResponse
func NewReleaseResponse(
	reason *enumerations.ReleaseResponseReason,
	userInformation *UserInformation,
) *ReleaseResponse {
	return &ReleaseResponse{
		Reason:          reason,
		UserInformation: userInformation,
	}
}

// FromBytes creates ReleaseResponse from bytes
func (r *ReleaseResponse) FromBytes(sourceBytes []byte) (*ReleaseResponse, error) {
	if len(sourceBytes) == 0 {
		return nil, fmt.Errorf("insufficient data for RLRE tag")
	}

	data := make([]byte, len(sourceBytes))
	copy(data, sourceBytes)

	tag := data[0]
	if tag != RLRETag {
		return nil, fmt.Errorf("bytes are not an RLRE APDU, tag is not %d, got %d", RLRETag, tag)
	}

	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for RLRE length")
	}

	length := int(data[1])
	data = data[2:]

	if len(data) != length {
		return nil, fmt.Errorf("the APDU data length does not correspond to length byte, expected %d, got %d", length, len(data))
	}

	// Parse tags
	objectDict := make(map[string]interface{})
	ber := encoding.NewBER()

	for len(data) > 0 {
		if len(data) < 2 {
			return nil, fmt.Errorf("insufficient data for tag and length")
		}

		objectTag := data[0]
		objectLength := int(data[1])
		data = data[2:]

		if len(data) < objectLength {
			return nil, fmt.Errorf("insufficient data for object, need %d bytes, got %d", objectLength, len(data))
		}

		objectData := data[:objectLength]
		data = data[objectLength:]

		var objectName string
		var parsedData interface{}
		var err error

		switch objectTag {
		case 0x80: // reason
			objectName = "reason"
			if len(objectData) > 0 {
				// Decode BER encoded integer
				tag, length, berData, err := ber.Decode(objectData, 1)
				if err != nil {
					return nil, fmt.Errorf("failed to decode reason: %w", err)
				}
				if !bytesEqual(tag, []byte{2}) { // Integer tag
					return nil, fmt.Errorf("reason is not an integer")
				}
				if len(berData) != int(length) || len(berData) == 0 {
					return nil, fmt.Errorf("invalid reason data length")
				}
				reason := enumerations.ReleaseResponseReason(berData[0])
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
			return nil, fmt.Errorf("could not find object with tag 0x%02x in RLRE definition", objectTag)
		}

		objectDict[objectName] = parsedData
	}

	reason, _ := objectDict["reason"].(*enumerations.ReleaseResponseReason)
	userInformation, _ := objectDict["user_information"].(*UserInformation)

	return NewReleaseResponse(reason, userInformation), nil
}

// ToBytes converts ReleaseResponse to bytes
func (r *ReleaseResponse) ToBytes() ([]byte, error) {
	ber := encoding.NewBER()
	rlreData := make([]byte, 0)

	if r.Reason != nil {
		// First encode the reason as a BER integer (tag 0x02)
		integerBytes, err := ber.Encode(2, []byte{byte(*r.Reason)})
		if err != nil {
			return nil, fmt.Errorf("failed to encode reason as integer: %w", err)
		}
		// Then wrap it with tag 0x80
		reasonBytes, err := ber.Encode(0x80, integerBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to encode reason: %w", err)
		}
		rlreData = append(rlreData, reasonBytes...)
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
		rlreData = append(rlreData, encoded...)
	}

	return ber.Encode(RLRETag, rlreData)
}

