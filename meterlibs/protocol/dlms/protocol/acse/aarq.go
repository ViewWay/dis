package acse

import (
	"fmt"

	"github.com/yimiliya/idis/meterlibs/protocol/dlms/enumerations"
	"github.com/yimiliya/idis/meterlibs/protocol/dlms/encoding"
)

// aarqShouldSetAuthenticated determines if authentication should be set based on mechanism
func aarqShouldSetAuthenticated(mechanism *enumerations.AuthenticationMechanism) bool {
	if mechanism == nil {
		return false
	}
	if *mechanism == enumerations.AuthenticationMechanismNone {
		return false
	}
	return true
}

// ApplicationAssociationRequest represents an AARQ (Application Association Request)
// It is used for starting an Application Association with a DLMS server (meter)
type ApplicationAssociationRequest struct {
	UserInformation                  *UserInformation
	SystemTitle                      []byte
	PublicCert                       []byte
	Authentication                   *enumerations.AuthenticationMechanism
	Ciphered                         bool
	AuthenticationValue              []byte
	CallingAEInvocationIdentifier   []byte
	CalledAPTitle                    []byte
	CalledAEQualifier                []byte
	CalledAPInvocationIdentifier     []byte
	CalledAEInvocationIdentifier    []byte
	CallingAPInvocationIdentifier    []byte
	ImplementationInformation        []byte
}

const AARQTag = 0x60 // Application 0 = 60H = 96

// NewApplicationAssociationRequest creates a new ApplicationAssociationRequest
func NewApplicationAssociationRequest(
	userInformation *UserInformation,
	systemTitle []byte,
	publicCert []byte,
	authentication *enumerations.AuthenticationMechanism,
	ciphered bool,
	authenticationValue []byte,
	callingAEInvocationIdentifier []byte,
) *ApplicationAssociationRequest {
	return &ApplicationAssociationRequest{
		UserInformation:                userInformation,
		SystemTitle:                    systemTitle,
		PublicCert:                     publicCert,
		Authentication:                 authentication,
		Ciphered:                       ciphered,
		AuthenticationValue:            authenticationValue,
		CallingAEInvocationIdentifier:  callingAEInvocationIdentifier,
	}
}

// SenderACSERequirements returns the AuthFunctionalUnit if authentication is needed
func (a *ApplicationAssociationRequest) SenderACSERequirements() *AuthFunctionalUnit {
	if aarqShouldSetAuthenticated(a.Authentication) {
		return NewAuthFunctionalUnit(true)
	}
	return nil
}

// MechanismName returns the MechanismName if authentication is used
func (a *ApplicationAssociationRequest) MechanismName() *MechanismName {
	if a.SenderACSERequirements() != nil && a.Authentication != nil {
		return NewMechanismName(*a.Authentication)
	}
	return nil
}

// ApplicationContextName returns the AppContextName based on ciphered setting
func (a *ApplicationAssociationRequest) ApplicationContextName() *AppContextName {
	if a.Ciphered {
		return NewAppContextName(true, true)
	}
	return NewAppContextName(true, false)
}

// ProtocolVersion returns the protocol version (always 0)
func (a *ApplicationAssociationRequest) ProtocolVersion() int {
	return 0
}

// FromBytes creates ApplicationAssociationRequest from bytes
func (a *ApplicationAssociationRequest) FromBytes(sourceBytes []byte) (*ApplicationAssociationRequest, error) {
	if len(sourceBytes) == 0 {
		return nil, fmt.Errorf("insufficient data for AARQ tag")
	}

	aarqData := make([]byte, len(sourceBytes))
	copy(aarqData, sourceBytes)

	aarqTag := aarqData[0]
	if aarqTag != AARQTag {
		return nil, fmt.Errorf("bytes are not an AARQ APDU, tag is not 0x60, got 0x%02x", aarqTag)
	}

	if len(aarqData) < 2 {
		return nil, fmt.Errorf("insufficient data for AARQ length")
	}

	aarqLength := int(aarqData[1])
	aarqData = aarqData[2:]

	if len(aarqData) != aarqLength {
		return nil, fmt.Errorf("the APDU data length does not correspond to length byte, expected %d, got %d", aarqLength, len(aarqData))
	}

	// Parse tags
	objectDict := make(map[string]interface{})
	ber := encoding.NewBER()

	for len(aarqData) > 0 {
		if len(aarqData) < 2 {
			return nil, fmt.Errorf("insufficient data for tag and length")
		}

		objectTag := aarqData[0]
		objectLength := int(aarqData[1])
		aarqData = aarqData[2:]

		if len(aarqData) < objectLength {
			return nil, fmt.Errorf("insufficient data for object, need %d bytes, got %d", objectLength, len(aarqData))
		}

		objectData := aarqData[:objectLength]
		aarqData = aarqData[objectLength:]

		var objectName string
		var parsedData interface{}
		var err error

		switch objectTag {
		case 0x80: // protocol_version
			objectName = "protocol_version"
			parsedData = nil // We assume version 1 and don't decode it
		case 0xA1: // application_context_name
			objectName = "application_context_name"
			appCtx := NewAppContextName(false, false)
			parsedData, err = appCtx.FromBytes(objectData)
			if err != nil {
				return nil, fmt.Errorf("failed to parse application_context_name: %w", err)
			}
		case 162: // called_ap_title
			objectName = "called_ap_title"
			parsedData = objectData
		case 163: // called_ae_qualifier
			objectName = "called_ae_qualifier"
			parsedData = objectData
		case 164: // called_ap_invocation_identifier
			objectName = "called_ap_invocation_identifier"
			parsedData = objectData
		case 165: // called_ae_invocation_identifier
			objectName = "called_ae_invocation_identifier"
			parsedData = objectData
		case 166: // calling_ap_title
			objectName = "calling_ap_title"
			// It is BER encoded universal tag octetstring. Simple handling
			if len(objectData) >= 2 {
				parsedData = objectData[2:] // Skip tag and length
			} else {
				parsedData = objectData
			}
		case 167: // calling_ae_qualifier
			objectName = "calling_ae_qualifier"
			// It is BER encoded universal tag octetstring. Simple handling
			if len(objectData) >= 2 {
				parsedData = objectData[2:] // Skip tag and length
			} else {
				parsedData = objectData
			}
		case 168: // calling_ap_invocation_identifier
			objectName = "calling_ap_invocation_identifier"
			parsedData = objectData
		case 169: // calling_ae_invocation_identifier
			objectName = "calling_ae_invocation_identifier"
			parsedData = objectData
		case 0x8A: // sender_acse_requirements
			objectName = "sender_acse_requirements"
			authFunc := NewAuthFunctionalUnit(false)
			parsedData, err = authFunc.FromBytes(objectData)
			if err != nil {
				return nil, fmt.Errorf("failed to parse sender_acse_requirements: %w", err)
			}
		case 0x8B: // mechanism_name
			objectName = "mechanism_name"
			mechName := NewMechanismName(enumerations.AuthenticationMechanismNone)
			parsedData, err = mechName.FromBytes(objectData)
			if err != nil {
				return nil, fmt.Errorf("failed to parse mechanism_name: %w", err)
			}
		case 0xAC: // calling_authentication_value
			objectName = "calling_authentication_value"
			authVal := &AuthenticationValue{}
			parsedData, err = authVal.FromBytes(objectData)
			if err != nil {
				return nil, fmt.Errorf("failed to parse calling_authentication_value: %w", err)
			}
		case 0xBD: // implementation_information
			objectName = "implementation_information"
			parsedData = objectData
		case 0xBE: // user_information
			objectName = "user_information"
			userInfo := &UserInformation{}
			parsedData, err = userInfo.FromBytes(objectData)
			if err != nil {
				return nil, fmt.Errorf("failed to parse user_information: %w", err)
			}
		default:
			return nil, fmt.Errorf("could not find object with tag 0x%02x in AARQ definition", objectTag)
		}

		objectDict[objectName] = parsedData
	}

	// Extract and validate required fields
	protocolVersion, _ := objectDict["protocol_version"].(int)
	if protocolVersion != 0 && protocolVersion != 1 {
		// We assume version 1 (0) and don't decode it
	}

	applicationContextName, ok := objectDict["application_context_name"].(*AppContextName)
	if !ok {
		return nil, fmt.Errorf("application_context_name is required")
	}

	ciphered := applicationContextName.CipheredAPDUs
	if !applicationContextName.LogicalNameRefs {
		return nil, fmt.Errorf("parsed an AARQ that uses Short Name Referencing")
	}

	senderACSERequirements, _ := objectDict["sender_acse_requirements"].(*AuthFunctionalUnit)
	mechanismName, _ := objectDict["mechanism_name"].(*MechanismName)

	var authentication *enumerations.AuthenticationMechanism
	if senderACSERequirements != nil && mechanismName != nil {
		if senderACSERequirements.Authentication {
			auth := mechanismName.Mechanism
			authentication = &auth
		}
	}

	systemTitle, _ := objectDict["calling_ap_title"].([]byte)
	publicCert, _ := objectDict["calling_ae_qualifier"].([]byte)

	authValue, _ := objectDict["calling_authentication_value"].(*AuthenticationValue)
	var authenticationValue []byte
	if authValue != nil {
		authenticationValue = authValue.Password
	}

	userInformation, ok := objectDict["user_information"].(*UserInformation)
	if !ok {
		return nil, fmt.Errorf("user_information is required")
	}

	calledAPTitle, _ := objectDict["called_ap_title"].([]byte)
	calledAEQualifier, _ := objectDict["called_ae_qualifier"].([]byte)
	calledAPInvocationIdentifier, _ := objectDict["called_ap_invocation_identifier"].([]byte)
	calledAEInvocationIdentifier, _ := objectDict["called_ae_invocation_identifier"].([]byte)
	callingAPInvocationIdentifier, _ := objectDict["calling_ap_invocation_identifier"].([]byte)
	callingAEInvocationIdentifier, _ := objectDict["calling_ae_invocation_identifier"].([]byte)
	implementationInformation, _ := objectDict["implementation_information"].([]byte)

	return &ApplicationAssociationRequest{
		UserInformation:                userInformation,
		SystemTitle:                    systemTitle,
		PublicCert:                     publicCert,
		Authentication:                 authentication,
		Ciphered:                       ciphered,
		AuthenticationValue:            authenticationValue,
		CalledAPTitle:                   calledAPTitle,
		CalledAEQualifier:               calledAEQualifier,
		CalledAPInvocationIdentifier:    calledAPInvocationIdentifier,
		CalledAEInvocationIdentifier:    calledAEInvocationIdentifier,
		CallingAPInvocationIdentifier:  callingAPInvocationIdentifier,
		CallingAEInvocationIdentifier:  callingAEInvocationIdentifier,
		ImplementationInformation:      implementationInformation,
	}, nil
}

// ToBytes converts ApplicationAssociationRequest to bytes
func (a *ApplicationAssociationRequest) ToBytes() ([]byte, error) {
	ber := encoding.NewBER()
	aarqData := make([]byte, 0)

	// Application context name
	appCtxName := a.ApplicationContextName()
	appCtxBytes, err := appCtxName.ToBytes()
	if err != nil {
		return nil, fmt.Errorf("failed to encode application_context_name: %w", err)
	}
	encodedAppCtx, err := ber.Encode(0xA1, appCtxBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to BER encode application_context_name: %w", err)
	}
	aarqData = append(aarqData, encodedAppCtx...)

	// Optional fields
	if a.CalledAPTitle != nil {
		encoded, err := ber.Encode(162, a.CalledAPTitle)
		if err != nil {
			return nil, fmt.Errorf("failed to encode called_ap_title: %w", err)
		}
		aarqData = append(aarqData, encoded...)
	}

	if a.CalledAEQualifier != nil {
		encoded, err := ber.Encode(163, a.CalledAEQualifier)
		if err != nil {
			return nil, fmt.Errorf("failed to encode called_ae_qualifier: %w", err)
		}
		aarqData = append(aarqData, encoded...)
	}

	if a.CalledAPInvocationIdentifier != nil {
		encoded, err := ber.Encode(164, a.CalledAPInvocationIdentifier)
		if err != nil {
			return nil, fmt.Errorf("failed to encode called_ap_invocation_identifier: %w", err)
		}
		aarqData = append(aarqData, encoded...)
	}

	if a.CalledAEInvocationIdentifier != nil {
		encoded, err := ber.Encode(165, a.CalledAEInvocationIdentifier)
		if err != nil {
			return nil, fmt.Errorf("failed to encode called_ae_invocation_identifier: %w", err)
		}
		aarqData = append(aarqData, encoded...)
	}

	if a.SystemTitle != nil {
		// Encode as BER octet string
		octetStringBytes, err := ber.Encode(4, a.SystemTitle)
		if err != nil {
			return nil, fmt.Errorf("failed to encode system_title as octet string: %w", err)
		}
		encoded, err := ber.Encode(166, octetStringBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to encode calling_ap_title: %w", err)
		}
		aarqData = append(aarqData, encoded...)
	}

	if a.PublicCert != nil {
		// Encode as BER octet string
		octetStringBytes, err := ber.Encode(4, a.PublicCert)
		if err != nil {
			return nil, fmt.Errorf("failed to encode public_cert as octet string: %w", err)
		}
		encoded, err := ber.Encode(167, octetStringBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to encode calling_ae_qualifier: %w", err)
		}
		aarqData = append(aarqData, encoded...)
	}

	if a.CallingAPInvocationIdentifier != nil {
		encoded, err := ber.Encode(168, a.CallingAPInvocationIdentifier)
		if err != nil {
			return nil, fmt.Errorf("failed to encode calling_ap_invocation_identifier: %w", err)
		}
		aarqData = append(aarqData, encoded...)
	}

	if a.CallingAEInvocationIdentifier != nil {
		encoded, err := ber.Encode(169, a.CallingAEInvocationIdentifier)
		if err != nil {
			return nil, fmt.Errorf("failed to encode calling_ae_invocation_identifier: %w", err)
		}
		aarqData = append(aarqData, encoded...)
	}

	if a.SenderACSERequirements() != nil {
		authFuncBytes, err := a.SenderACSERequirements().ToBytes()
		if err != nil {
			return nil, fmt.Errorf("failed to encode sender_acse_requirements: %w", err)
		}
		encoded, err := ber.Encode(0x8A, authFuncBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to BER encode sender_acse_requirements: %w", err)
		}
		aarqData = append(aarqData, encoded...)
	}

	if a.MechanismName() != nil {
		mechNameBytes, err := a.MechanismName().ToBytes()
		if err != nil {
			return nil, fmt.Errorf("failed to encode mechanism_name: %w", err)
		}
		encoded, err := ber.Encode(0x8B, mechNameBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to BER encode mechanism_name: %w", err)
		}
		aarqData = append(aarqData, encoded...)
	}

	if a.AuthenticationValue != nil {
		authVal, err := NewAuthenticationValue(a.AuthenticationValue, "chars")
		if err != nil {
			return nil, fmt.Errorf("failed to create authentication value: %w", err)
		}
		authValBytes, err := authVal.ToBytes()
		if err != nil {
			return nil, fmt.Errorf("failed to encode authentication_value: %w", err)
		}
		encoded, err := ber.Encode(0xAC, authValBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to BER encode calling_authentication_value: %w", err)
		}
		aarqData = append(aarqData, encoded...)
	}

	if a.ImplementationInformation != nil {
		encoded, err := ber.Encode(0xBD, a.ImplementationInformation)
		if err != nil {
			return nil, fmt.Errorf("failed to encode implementation_information: %w", err)
		}
		aarqData = append(aarqData, encoded...)
	}

	if a.UserInformation != nil {
		userInfoBytes, err := a.UserInformation.ToBytes()
		if err != nil {
			return nil, fmt.Errorf("failed to encode user_information: %w", err)
		}
		encoded, err := ber.Encode(0xBE, userInfoBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to BER encode user_information: %w", err)
		}
		aarqData = append(aarqData, encoded...)
	}

	return ber.Encode(AARQTag, aarqData)
}

