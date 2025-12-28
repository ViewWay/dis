package acse

import (
	"encoding/binary"
	"fmt"

	"github.com/yimiliya/idis/meterlibs/protocol/dlms/enumerations"
	"github.com/yimiliya/idis/meterlibs/protocol/dlms/encoding"
)

// Asn1Integer wraps Integers for BER encoding
type Asn1Integer struct {
	Value int
}

const Asn1IntegerTag = 2 // ASN1 Universal tag 2, Integer

// NewAsn1Integer creates a new Asn1Integer
func NewAsn1Integer(value int) *Asn1Integer {
	return &Asn1Integer{Value: value}
}

// FromBytes creates Asn1Integer from bytes
func (a *Asn1Integer) FromBytes(sourceBytes []byte) (*Asn1Integer, error) {
	ber := encoding.NewBER()
	tag, length, data, err := ber.Decode(sourceBytes, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to decode BER: %w", err)
	}

	if !bytesEqual(tag, []byte{Asn1IntegerTag}) {
		return nil, fmt.Errorf("data provided is not of the correct type, tag is %v but should be %d", tag, Asn1IntegerTag)
	}

	var value int
	if len(data) == 1 {
		value = int(data[0])
	} else if len(data) == 2 {
		value = int(binary.BigEndian.Uint16(data))
	} else if len(data) == 4 {
		value = int(binary.BigEndian.Uint32(data))
	} else {
		return nil, fmt.Errorf("unsupported integer length: %d", len(data))
	}

	return NewAsn1Integer(value), nil
}

// ToBytes converts Asn1Integer to bytes
func (a *Asn1Integer) ToBytes() ([]byte, error) {
	ber := encoding.NewBER()
	// For values that fit in 1 byte
	if a.Value >= 0 && a.Value <= 255 {
		return ber.Encode(Asn1IntegerTag, []byte{byte(a.Value)})
	}
	return nil, fmt.Errorf("integer value %d does not fit in 1 byte", a.Value)
}

// ResultSourceDiagnostics represents result source diagnostics
type ResultSourceDiagnostics struct {
	Name  string // "acse-service-user" or "acse-service-provider"
	Value int
}

// NewResultSourceDiagnostics creates a new ResultSourceDiagnostics
func NewResultSourceDiagnostics(name string, value int) *ResultSourceDiagnostics {
	return &ResultSourceDiagnostics{
		Name:  name,
		Value: value,
	}
}

// FromBytes creates ResultSourceDiagnostics from bytes
func (r *ResultSourceDiagnostics) FromBytes(sourceBytes []byte) (*ResultSourceDiagnostics, error) {
	ber := encoding.NewBER()
	
	// Try to decode as acse-service-user (tag 0x81)
	if len(sourceBytes) >= 2 && sourceBytes[0] == 0x81 {
		tag, _, data, err := ber.Decode(sourceBytes, 1)
		if err == nil && bytesEqual(tag, []byte{0x81}) {
			if len(data) > 0 {
				value := int(data[0])
				return NewResultSourceDiagnostics("acse-service-user", value), nil
			}
		}
	}
	
	// Try to decode as acse-service-provider (tag 0x82)
	if len(sourceBytes) >= 2 && sourceBytes[0] == 0x82 {
		tag, _, data, err := ber.Decode(sourceBytes, 1)
		if err == nil && bytesEqual(tag, []byte{0x82}) {
			if len(data) > 0 {
				value := int(data[0])
				return NewResultSourceDiagnostics("acse-service-provider", value), nil
			}
		}
	}
	
	return nil, fmt.Errorf("failed to parse result source diagnostics")
}

// ToBytes converts ResultSourceDiagnostics to bytes
func (r *ResultSourceDiagnostics) ToBytes() ([]byte, error) {
	ber := encoding.NewBER()
	var tag int
	if r.Name == "acse-service-user" {
		tag = 0x81
	} else if r.Name == "acse-service-provider" {
		tag = 0x82
	} else {
		return nil, fmt.Errorf("invalid result source diagnostics name: %s", r.Name)
	}
	return ber.Encode(tag, []byte{byte(r.Value)})
}

// ApplicationAssociationResponse represents an AARE (Application Association Response)
type ApplicationAssociationResponse struct {
	Result                      enumerations.AssociationResult
	ResultSourceDiagnostics     interface{} // AcseServiceUserDiagnostics or AcseServiceProviderDiagnostics
	Ciphered                    bool
	Authentication              *enumerations.AuthenticationMechanism
	SystemTitle                 []byte
	PublicCert                  []byte
	AuthenticationValue         []byte
	UserInformation             *UserInformation
	ImplementationInformation   []byte
	RespondingAPInvocationID    []byte
	RespondingAEInvocationID    []byte
}

const AARETag = 0x61 // Application 1

// NewApplicationAssociationResponse creates a new ApplicationAssociationResponse
func NewApplicationAssociationResponse(
	result enumerations.AssociationResult,
	resultSourceDiagnostics interface{},
	ciphered bool,
	authentication *enumerations.AuthenticationMechanism,
	systemTitle []byte,
	publicCert []byte,
	authenticationValue []byte,
	userInformation *UserInformation,
) *ApplicationAssociationResponse {
	return &ApplicationAssociationResponse{
		Result:                  result,
		ResultSourceDiagnostics: resultSourceDiagnostics,
		Ciphered:                ciphered,
		Authentication:          authentication,
		SystemTitle:             systemTitle,
		PublicCert:              publicCert,
		AuthenticationValue:     authenticationValue,
		UserInformation:         userInformation,
	}
}

// ResponderACSERequirements returns the AuthFunctionalUnit if authentication is needed
func (a *ApplicationAssociationResponse) ResponderACSERequirements() *AuthFunctionalUnit {
	if aarqShouldSetAuthenticated(a.Authentication) {
		return NewAuthFunctionalUnit(true)
	}
	return nil
}

// MechanismName returns the MechanismName if authentication is used
func (a *ApplicationAssociationResponse) MechanismName() *MechanismName {
	if a.ResponderACSERequirements() != nil && a.Authentication != nil {
		return NewMechanismName(*a.Authentication)
	}
	return nil
}

// ApplicationContextName returns the AppContextName based on ciphered setting
func (a *ApplicationAssociationResponse) ApplicationContextName() *AppContextName {
	if a.Ciphered {
		return NewAppContextName(true, true)
	}
	return NewAppContextName(true, false)
}

// ProtocolVersion returns the protocol version (always 0)
func (a *ApplicationAssociationResponse) ProtocolVersion() int {
	return 0
}

// FromBytes creates ApplicationAssociationResponse from bytes
func (a *ApplicationAssociationResponse) FromBytes(sourceBytes []byte) (*ApplicationAssociationResponse, error) {
	if len(sourceBytes) == 0 {
		return nil, fmt.Errorf("insufficient data for AARE tag")
	}

	aareData := make([]byte, len(sourceBytes))
	copy(aareData, sourceBytes)

	aareTag := aareData[0]
	if aareTag != AARETag {
		return nil, fmt.Errorf("bytes are not an AARE APDU, tag is not 0x61, got 0x%02x", aareTag)
	}

	if len(aareData) < 2 {
		return nil, fmt.Errorf("insufficient data for AARE length")
	}

	aareLength := int(aareData[1])
	aareData = aareData[2:]

	if len(aareData) != aareLength {
		return nil, fmt.Errorf("the APDU data length does not correspond to length byte, expected %d, got %d", aareLength, len(aareData))
	}

	// Parse tags
	objectDict := make(map[string]interface{})
	ber := encoding.NewBER()

	for len(aareData) > 0 {
		if len(aareData) < 2 {
			return nil, fmt.Errorf("insufficient data for tag and length")
		}

		objectTag := aareData[0]
		objectLength := int(aareData[1])
		aareData = aareData[2:]

		if len(aareData) < objectLength {
			return nil, fmt.Errorf("insufficient data for object, need %d bytes, got %d", objectLength, len(aareData))
		}

		objectData := aareData[:objectLength]
		aareData = aareData[objectLength:]

		var objectName string
		var parsedData interface{}
		var err error

		switch objectTag {
		case 128: // protocol_version
			objectName = "protocol_version"
			parsedData = nil // We assume version 1 and don't decode it
		case 161: // application_context_name
			objectName = "application_context_name"
			appCtx := NewAppContextName(false, false)
			parsedData, err = appCtx.FromBytes(objectData)
			if err != nil {
				return nil, fmt.Errorf("failed to parse application_context_name: %w", err)
			}
		case 162: // result
			objectName = "result"
			asn1Int := &Asn1Integer{}
			parsedData, err = asn1Int.FromBytes(objectData)
			if err != nil {
				return nil, fmt.Errorf("failed to parse result: %w", err)
			}
		case 163: // result_source_diagnostics
			objectName = "result_source_diagnostics"
			rsd := &ResultSourceDiagnostics{}
			parsedData, err = rsd.FromBytes(objectData)
			if err != nil {
				return nil, fmt.Errorf("failed to parse result_source_diagnostics: %w", err)
			}
		case 164: // responding_ap_title
			objectName = "responding_ap_title"
			// It is BER encoded universal tag octetstring. Simple handling
			if len(objectData) >= 2 {
				parsedData = objectData[2:] // Skip tag and length
			} else {
				parsedData = objectData
			}
		case 165: // responding_ae_qualifier
			objectName = "responding_ae_qualifier"
			// It is BER encoded universal tag octetstring. Simple handling
			if len(objectData) >= 2 {
				parsedData = objectData[2:] // Skip tag and length
			} else {
				parsedData = objectData
			}
		case 166: // responding_ap_invocation_id
			objectName = "responding_ap_invocation_id"
			parsedData = objectData
		case 167: // responding_ae_invocation_id
			objectName = "responding_ae_invocation_id"
			parsedData = objectData
		case 0x88: // responder_acse_requirements
			objectName = "responder_acse_requirements"
			authFunc := NewAuthFunctionalUnit(false)
			parsedData, err = authFunc.FromBytes(objectData)
			if err != nil {
				return nil, fmt.Errorf("failed to parse responder_acse_requirements: %w", err)
			}
		case 0x89: // mechanism_name
			objectName = "mechanism_name"
			mechName := NewMechanismName(enumerations.AuthenticationMechanismNone)
			parsedData, err = mechName.FromBytes(objectData)
			if err != nil {
				return nil, fmt.Errorf("failed to parse mechanism_name: %w", err)
			}
		case 170: // responding_authentication_value
			objectName = "responding_authentication_value"
			authVal := &AuthenticationValue{}
			parsedData, err = authVal.FromBytes(objectData)
			if err != nil {
				return nil, fmt.Errorf("failed to parse responding_authentication_value: %w", err)
			}
		case 189: // implementation_information
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
			return nil, fmt.Errorf("could not find object with tag 0x%02x in AARE definition", objectTag)
		}

		objectDict[objectName] = parsedData
	}

	// Extract and validate required fields
	applicationContextName, ok := objectDict["application_context_name"].(*AppContextName)
	if !ok {
		return nil, fmt.Errorf("application_context_name is required")
	}

	ciphered := applicationContextName.CipheredAPDUs
	if !applicationContextName.LogicalNameRefs {
		return nil, fmt.Errorf("AARE requests use of Short Name referencing which is not supported")
	}

	// Transform result into enum
	resultInt, ok := objectDict["result"].(*Asn1Integer)
	if !ok {
		return nil, fmt.Errorf("result is required")
	}
	result := enumerations.AssociationResult(resultInt.Value)

	// Transform source diagnostic into enum
	sourceDiagnostic, ok := objectDict["result_source_diagnostics"].(*ResultSourceDiagnostics)
	if !ok {
		return nil, fmt.Errorf("result_source_diagnostics is required")
	}

	var resultSourceDiagnostics interface{}
	if sourceDiagnostic.Name == "acse-service-user" {
		resultSourceDiagnostics = enumerations.AcseServiceUserDiagnostics(sourceDiagnostic.Value)
	} else if sourceDiagnostic.Name == "acse-service-provider" {
		resultSourceDiagnostics = enumerations.AcseServiceProviderDiagnostics(sourceDiagnostic.Value)
	} else {
		return nil, fmt.Errorf("not a valid choice of result_source_diagnostics")
	}

	responderACSERequirements, _ := objectDict["responder_acse_requirements"].(*AuthFunctionalUnit)
	mechanismName, _ := objectDict["mechanism_name"].(*MechanismName)

	var authentication *enumerations.AuthenticationMechanism
	if responderACSERequirements != nil && mechanismName != nil {
		if responderACSERequirements.Authentication {
			auth := mechanismName.Mechanism
			authentication = &auth
		}
	}

	systemTitle, _ := objectDict["responding_ap_title"].([]byte)
	publicCert, _ := objectDict["responding_ae_qualifier"].([]byte)

	authValue, _ := objectDict["responding_authentication_value"].(*AuthenticationValue)
	var authenticationValue []byte
	if authValue != nil {
		authenticationValue = authValue.Password
	}

	userInformation, _ := objectDict["user_information"].(*UserInformation)

	respondingAPInvocationID, _ := objectDict["responding_ap_invocation_id"].([]byte)
	respondingAEInvocationID, _ := objectDict["responding_ae_invocation_id"].([]byte)
	implementationInformation, _ := objectDict["implementation_information"].([]byte)

	return &ApplicationAssociationResponse{
		Result:                      result,
		ResultSourceDiagnostics:     resultSourceDiagnostics,
		Ciphered:                    ciphered,
		Authentication:              authentication,
		SystemTitle:                 systemTitle,
		PublicCert:                  publicCert,
		AuthenticationValue:        authenticationValue,
		UserInformation:            userInformation,
		RespondingAPInvocationID:    respondingAPInvocationID,
		RespondingAEInvocationID:    respondingAEInvocationID,
		ImplementationInformation:  implementationInformation,
	}, nil
}

// ToBytes converts ApplicationAssociationResponse to bytes
func (a *ApplicationAssociationResponse) ToBytes() ([]byte, error) {
	ber := encoding.NewBER()
	aareData := make([]byte, 0)

	// Application context name
	appCtxName := a.ApplicationContextName()
	appCtxBytes, err := appCtxName.ToBytes()
	if err != nil {
		return nil, fmt.Errorf("failed to encode application_context_name: %w", err)
	}
	encodedAppCtx, err := ber.Encode(161, appCtxBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to BER encode application_context_name: %w", err)
	}
	aareData = append(aareData, encodedAppCtx...)

	// Result
	resultInt := NewAsn1Integer(int(a.Result))
	resultBytes, err := resultInt.ToBytes()
	if err != nil {
		return nil, fmt.Errorf("failed to encode result: %w", err)
	}
	encodedResult, err := ber.Encode(162, resultBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to BER encode result: %w", err)
	}
	aareData = append(aareData, encodedResult...)

	// Result source diagnostics
	if a.ResultSourceDiagnostics != nil {
		var rsd *ResultSourceDiagnostics
		switch diag := a.ResultSourceDiagnostics.(type) {
		case enumerations.AcseServiceUserDiagnostics:
			rsd = NewResultSourceDiagnostics("acse-service-user", int(diag))
		case enumerations.AcseServiceProviderDiagnostics:
			rsd = NewResultSourceDiagnostics("acse-service-provider", int(diag))
		default:
			return nil, fmt.Errorf("unsupported result source diagnostics type: %T", a.ResultSourceDiagnostics)
		}
		rsdBytes, err := rsd.ToBytes()
		if err != nil {
			return nil, fmt.Errorf("failed to encode result_source_diagnostics: %w", err)
		}
		encodedRSD, err := ber.Encode(163, rsdBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to BER encode result_source_diagnostics: %w", err)
		}
		aareData = append(aareData, encodedRSD...)
	}

	// Optional fields
	if a.SystemTitle != nil {
		octetStringBytes, err := ber.Encode(4, a.SystemTitle)
		if err != nil {
			return nil, fmt.Errorf("failed to encode system_title as octet string: %w", err)
		}
		encoded, err := ber.Encode(164, octetStringBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to encode responding_ap_title: %w", err)
		}
		aareData = append(aareData, encoded...)
	}

	if a.PublicCert != nil {
		octetStringBytes, err := ber.Encode(4, a.PublicCert)
		if err != nil {
			return nil, fmt.Errorf("failed to encode public_cert as octet string: %w", err)
		}
		encoded, err := ber.Encode(165, octetStringBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to encode responding_ae_qualifier: %w", err)
		}
		aareData = append(aareData, encoded...)
	}

	if a.RespondingAPInvocationID != nil {
		encoded, err := ber.Encode(166, a.RespondingAPInvocationID)
		if err != nil {
			return nil, fmt.Errorf("failed to encode responding_ap_invocation_id: %w", err)
		}
		aareData = append(aareData, encoded...)
	}

	if a.RespondingAEInvocationID != nil {
		encoded, err := ber.Encode(167, a.RespondingAEInvocationID)
		if err != nil {
			return nil, fmt.Errorf("failed to encode responding_ae_invocation_id: %w", err)
		}
		aareData = append(aareData, encoded...)
	}

	if a.ResponderACSERequirements() != nil {
		authFuncBytes, err := a.ResponderACSERequirements().ToBytes()
		if err != nil {
			return nil, fmt.Errorf("failed to encode responder_acse_requirements: %w", err)
		}
		encoded, err := ber.Encode(0x88, authFuncBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to BER encode responder_acse_requirements: %w", err)
		}
		aareData = append(aareData, encoded...)
	}

	if a.MechanismName() != nil {
		mechNameBytes, err := a.MechanismName().ToBytes()
		if err != nil {
			return nil, fmt.Errorf("failed to encode mechanism_name: %w", err)
		}
		encoded, err := ber.Encode(0x89, mechNameBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to BER encode mechanism_name: %w", err)
		}
		aareData = append(aareData, encoded...)
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
		encoded, err := ber.Encode(170, authValBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to BER encode responding_authentication_value: %w", err)
		}
		aareData = append(aareData, encoded...)
	}

	if a.ImplementationInformation != nil {
		encoded, err := ber.Encode(189, a.ImplementationInformation)
		if err != nil {
			return nil, fmt.Errorf("failed to encode implementation_information: %w", err)
		}
		aareData = append(aareData, encoded...)
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
		aareData = append(aareData, encoded...)
	}

	return ber.Encode(AARETag, aareData)
}

