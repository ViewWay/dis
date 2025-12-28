package acse

import (
	"fmt"

	"github.com/yimiliya/idis/meterlibs/protocol/dlms/enumerations"
	"github.com/yimiliya/idis/meterlibs/protocol/dlms/encoding"
)

// AbstractAcseApdu is the base interface for ACSE APDUs
type AbstractAcseApdu interface {
	FromBytes(sourceBytes []byte) (AbstractAcseApdu, error)
	ToBytes() ([]byte, error)
}

// DLMSObjectIdentifier represents a DLMS object identifier
type DLMSObjectIdentifier struct {
	Tag    []byte
	Prefix []byte
}

// NewDLMSObjectIdentifier creates a new DLMS object identifier
func NewDLMSObjectIdentifier() *DLMSObjectIdentifier {
	return &DLMSObjectIdentifier{
		Tag:    []byte{0x06},
		Prefix: []byte{0x60, 0x85, 0x74, 0x05, 0x08},
	}
}

// AppContextName defines how to reference objects in the meter and if ciphered APDUs are allowed
type AppContextName struct {
	*DLMSObjectIdentifier
	AppContext      int
	ValidContextIDs []int
	LogicalNameRefs bool
	CipheredAPDUs   bool
}

// NewAppContextName creates a new AppContextName
func NewAppContextName(logicalNameRefs, cipheredAPDUs bool) *AppContextName {
	return &AppContextName{
		DLMSObjectIdentifier: NewDLMSObjectIdentifier(),
		AppContext:           1,
		ValidContextIDs:      []int{1, 2, 3, 4},
		LogicalNameRefs:      logicalNameRefs,
		CipheredAPDUs:        cipheredAPDUs,
	}
}

// ContextID returns the context ID based on logical name refs and ciphered APDUs
func (a *AppContextName) ContextID() int {
	if a.LogicalNameRefs && !a.CipheredAPDUs {
		return 1
	} else if !a.LogicalNameRefs && !a.CipheredAPDUs {
		return 2
	} else if a.LogicalNameRefs && a.CipheredAPDUs {
		return 3
	} else if !a.LogicalNameRefs && a.CipheredAPDUs {
		return 4
	}
	panic("Combination of logical name ref and ciphered apdus not possible")
}

// FromBytes creates AppContextName from bytes
func (a *AppContextName) FromBytes(data []byte) (*AppContextName, error) {
	ber := encoding.NewBER()
	tag, length, berData, err := ber.Decode(data, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to decode BER: %w", err)
	}

	if !bytesEqual(tag, a.DLMSObjectIdentifier.Tag) {
		return nil, fmt.Errorf("tag %v is not a valid tag for ObjectIdentifiers", tag)
	}

	if len(berData) < 1 {
		return nil, fmt.Errorf("insufficient data for context ID")
	}

	contextID := int(berData[len(berData)-1])
	valid := false
	for _, id := range a.ValidContextIDs {
		if id == contextID {
			valid = true
			break
		}
	}
	if !valid {
		return nil, fmt.Errorf("context_id of %d is not valid", contextID)
	}

	totalPrefix := berData[:len(berData)-1]
	expectedPrefix := append(a.Prefix, byte(a.AppContext))
	if !bytesEqual(totalPrefix, expectedPrefix) {
		return nil, fmt.Errorf("static part of object id is not correct according to DLMS: %v", totalPrefix)
	}

	settings := getSettingsByContextID(contextID)
	return NewAppContextName(settings["logical_name_refs"].(bool), settings["ciphered_apdus"].(bool)), nil
}

// ToBytes converts AppContextName to bytes
func (a *AppContextName) ToBytes() ([]byte, error) {
	ber := encoding.NewBER()
	totalData := append(a.Prefix, byte(a.AppContext), byte(a.ContextID()))
	return ber.Encode(a.Tag, totalData)
}

func getSettingsByContextID(contextID int) map[string]interface{} {
	settings := map[int]map[string]interface{}{
		1: {"logical_name_refs": true, "ciphered_apdus": false},
		2: {"logical_name_refs": false, "ciphered_apdus": false},
		3: {"logical_name_refs": true, "ciphered_apdus": true},
		4: {"logical_name_refs": false, "ciphered_apdus": true},
	}
	return settings[contextID]
}

// MechanismName represents an authentication mechanism name
type MechanismName struct {
	*DLMSObjectIdentifier
	AppContext int
	Mechanism  enumerations.AuthenticationMechanism
}

// NewMechanismName creates a new MechanismName
func NewMechanismName(mechanism enumerations.AuthenticationMechanism) *MechanismName {
	return &MechanismName{
		DLMSObjectIdentifier: NewDLMSObjectIdentifier(),
		AppContext:           2,
		Mechanism:            mechanism,
	}
}

// FromBytes creates MechanismName from bytes
func (m *MechanismName) FromBytes(data []byte) (*MechanismName, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for mechanism ID")
	}

	mechanismID := int(data[len(data)-1])
	totalPrefix := data[:len(data)-1]
	expectedPrefix := append(m.Prefix, byte(m.AppContext))
	if !bytesEqual(totalPrefix, expectedPrefix) {
		return nil, fmt.Errorf("static part of object id is not correct according to DLMS: %v", totalPrefix)
	}

	return NewMechanismName(enumerations.AuthenticationMechanism(mechanismID)), nil
}

// ToBytes converts MechanismName to bytes
func (m *MechanismName) ToBytes() ([]byte, error) {
	totalData := append(m.Prefix, byte(m.AppContext), byte(m.Mechanism))
	return totalData, nil
}

// AuthenticationValue holds "password" in the AARQ and AARE
type AuthenticationValue struct {
	Password     []byte
	PasswordType string
}

var allowedPasswordTypes = []string{"chars", "bits"}

// NewAuthenticationValue creates a new AuthenticationValue
func NewAuthenticationValue(password []byte, passwordType string) (*AuthenticationValue, error) {
	valid := false
	for _, pt := range allowedPasswordTypes {
		if pt == passwordType {
			valid = true
			break
		}
	}
	if !valid {
		return nil, fmt.Errorf("%s is not a valid auth value type", passwordType)
	}

	return &AuthenticationValue{
		Password:     password,
		PasswordType: passwordType,
	}, nil
}

// FromBytes creates AuthenticationValue from bytes
func (a *AuthenticationValue) FromBytes(data []byte) (*AuthenticationValue, error) {
	ber := encoding.NewBER()
	tag, length, berData, err := ber.Decode(data, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to decode BER: %w", err)
	}

	var passwordType string
	if bytesEqual(tag, []byte{0x80}) {
		passwordType = "chars"
	} else if bytesEqual(tag, []byte{0x81}) {
		passwordType = "bits"
	} else {
		return nil, fmt.Errorf("tag %v is not valid for password", tag)
	}

	return &AuthenticationValue{
		Password:     berData,
		PasswordType: passwordType,
	}, nil
}

// ToBytes converts AuthenticationValue to bytes
func (a *AuthenticationValue) ToBytes() ([]byte, error) {
	ber := encoding.NewBER()
	var tag int
	if a.PasswordType == "chars" {
		tag = 0x80
	} else if a.PasswordType == "bits" {
		tag = 0x81
	} else {
		return nil, fmt.Errorf("invalid password type: %s", a.PasswordType)
	}
	return ber.Encode(tag, a.Password)
}

// AuthFunctionalUnit consists of 2 bytes
type AuthFunctionalUnit struct {
	Authentication bool
}

// NewAuthFunctionalUnit creates a new AuthFunctionalUnit
func NewAuthFunctionalUnit(authentication bool) *AuthFunctionalUnit {
	return &AuthFunctionalUnit{
		Authentication: authentication,
	}
}

// FromBytes creates AuthFunctionalUnit from bytes
func (a *AuthFunctionalUnit) FromBytes(data []byte) (*AuthFunctionalUnit, error) {
	if len(data) != 2 {
		return nil, fmt.Errorf("authentication functional unit data should be 2 bytes, got: %v", data)
	}
	lastByte := data[1] != 0
	return NewAuthFunctionalUnit(lastByte), nil
}

// ToBytes converts AuthFunctionalUnit to bytes
func (a *AuthFunctionalUnit) ToBytes() ([]byte, error) {
	if a.Authentication {
		return []byte{0x07, 0x80}, nil
	}
	// when not using authentication this the sender-acse-requirements should not be in the data
	return nil, nil
}

// bytesEqual compares two byte slices
func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

