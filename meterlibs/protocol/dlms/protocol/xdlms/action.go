package xdlms

import (
	"fmt"

	"github.com/yimiliya/idis/meterlibs/protocol/dlms/cosem"
	"github.com/yimiliya/idis/meterlibs/protocol/dlms/enumerations"
)

// ActionRequestNormal represents an Action request normal
const ActionRequestTag = 195

type ActionRequestNormal struct {
	*BaseXDlmsApdu
	CosemMethod         *cosem.CosemMethod
	Data                []byte
	InvokeIdAndPriority *InvokeIdAndPriority
}

// NewActionRequestNormal creates a new ActionRequestNormal
func NewActionRequestNormal(
	cosemMethod *cosem.CosemMethod,
	data []byte,
	invokeIdAndPriority *InvokeIdAndPriority,
) *ActionRequestNormal {
	return &ActionRequestNormal{
		BaseXDlmsApdu: &BaseXDlmsApdu{
			Tag: ActionRequestTag,
		},
		CosemMethod:         cosemMethod,
		Data:                data,
		InvokeIdAndPriority: invokeIdAndPriority,
	}
}

// FromBytes creates ActionRequestNormal from bytes
func (a *ActionRequestNormal) FromBytes(data []byte) (*ActionRequestNormal, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for ActionRequest")
	}
	
	tag := data[0]
	if tag != ActionRequestTag {
		return nil, fmt.Errorf("tag %d is not the correct tag for an ActionRequest, should be %d", tag, ActionRequestTag)
	}
	
	requestType := enumerations.ActionType(data[1])
	if requestType != enumerations.ActionTypeNormal {
		return nil, fmt.Errorf("bytes are not representing a ActionRequestNormal. Action type is %d", requestType)
	}
	
	data = data[2:]
	
	// Parse invoke_id_and_priority
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for invoke_id_and_priority")
	}
	invokeIdAndPriority, err := (&InvokeIdAndPriority{}).FromBytes(data[:1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse invoke_id_and_priority: %w", err)
	}
	data = data[1:]
	
	// Parse cosem_method (9 bytes)
	if len(data) < 9 {
		return nil, fmt.Errorf("insufficient data for cosem_method")
	}
	cosemMethod, err := (&cosem.CosemMethod{}).FromBytes(data[:9])
	if err != nil {
		return nil, fmt.Errorf("failed to parse cosem_method: %w", err)
	}
	data = data[9:]
	
	// Parse has_data flag
	var requestData []byte
	if len(data) > 0 {
		hasData := data[0] != 0
		data = data[1:]
		if hasData {
			requestData = make([]byte, len(data))
			copy(requestData, data)
		}
	}
	
	return NewActionRequestNormal(cosemMethod, requestData, invokeIdAndPriority), nil
}

// ToBytes converts ActionRequestNormal to bytes
func (a *ActionRequestNormal) ToBytes() ([]byte, error) {
	result := []byte{ActionRequestTag}
	result = append(result, byte(enumerations.ActionTypeNormal))
	
	invokeBytes := a.InvokeIdAndPriority.ToBytes()
	result = append(result, invokeBytes...)
	
	cosemBytes := a.CosemMethod.ToBytes()
	result = append(result, cosemBytes...)
	
	if len(a.Data) > 0 {
		result = append(result, 0x01)
		result = append(result, a.Data...)
	} else {
		result = append(result, 0x00)
	}
	
	return result, nil
}

// ActionResponseNormal represents an Action response normal
const ActionResponseTag = 199

type ActionResponseNormal struct {
	*BaseXDlmsApdu
	Status              enumerations.ActionResultStatus
	InvokeIdAndPriority *InvokeIdAndPriority
}

// NewActionResponseNormal creates a new ActionResponseNormal
func NewActionResponseNormal(
	status enumerations.ActionResultStatus,
	invokeIdAndPriority *InvokeIdAndPriority,
) *ActionResponseNormal {
	return &ActionResponseNormal{
		BaseXDlmsApdu: &BaseXDlmsApdu{
			Tag: ActionResponseTag,
		},
		Status:              status,
		InvokeIdAndPriority: invokeIdAndPriority,
	}
}

// FromBytes creates ActionResponseNormal from bytes
func (a *ActionResponseNormal) FromBytes(data []byte) (*ActionResponseNormal, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for ActionResponse")
	}
	
	tag := data[0]
	if tag != ActionResponseTag {
		return nil, fmt.Errorf("tag %d is not correct for ActionResponse. Should be %d", tag, ActionResponseTag)
	}
	
	actionType := enumerations.ActionType(data[1])
	if actionType != enumerations.ActionTypeNormal {
		return nil, fmt.Errorf("bytes are not representing a ActionResponseNormal. Action type is %d", actionType)
	}
	
	data = data[2:]
	
	// Parse invoke_id_and_priority
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for invoke_id_and_priority")
	}
	invokeIdAndPriority, err := (&InvokeIdAndPriority{}).FromBytes(data[:1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse invoke_id_and_priority: %w", err)
	}
	data = data[1:]
	
	// Parse status
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for status")
	}
	status := enumerations.ActionResultStatus(data[0])
	data = data[1:]
	
	// Parse has_data flag (should be 0 for normal response)
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for has_data flag")
	}
		hasData := data[0] != 0
	data = data[1:] // Advance pointer after reading the flag
		if hasData {
		return nil, fmt.Errorf("ActionResponse has data and should not be a ActionResponseNormal")
	}
	
	return NewActionResponseNormal(status, invokeIdAndPriority), nil
}

// ToBytes converts ActionResponseNormal to bytes
func (a *ActionResponseNormal) ToBytes() ([]byte, error) {
	result := []byte{ActionResponseTag}
	result = append(result, byte(enumerations.ActionTypeNormal))
	
	invokeBytes := a.InvokeIdAndPriority.ToBytes()
	result = append(result, invokeBytes...)
	
	result = append(result, byte(a.Status))
	result = append(result, 0x00) // has_data = false
	
	return result, nil
}

// ActionResponseNormalWithData represents an Action response normal with data
type ActionResponseNormalWithData struct {
	*BaseXDlmsApdu
	Status              enumerations.ActionResultStatus
	Data                []byte
	InvokeIdAndPriority *InvokeIdAndPriority
}

// NewActionResponseNormalWithData creates a new ActionResponseNormalWithData
func NewActionResponseNormalWithData(
	status enumerations.ActionResultStatus,
	data []byte,
	invokeIdAndPriority *InvokeIdAndPriority,
) *ActionResponseNormalWithData {
	return &ActionResponseNormalWithData{
		BaseXDlmsApdu: &BaseXDlmsApdu{
			Tag: ActionResponseTag,
		},
		Status:              status,
		Data:                data,
		InvokeIdAndPriority: invokeIdAndPriority,
	}
}

// FromBytes creates ActionResponseNormalWithData from bytes
func (a *ActionResponseNormalWithData) FromBytes(data []byte) (*ActionResponseNormalWithData, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for ActionResponse")
	}
	
	tag := data[0]
	if tag != ActionResponseTag {
		return nil, fmt.Errorf("tag %d is not correct for ActionResponse. Should be %d", tag, ActionResponseTag)
	}
	
	actionType := enumerations.ActionType(data[1])
	if actionType != enumerations.ActionTypeNormal {
		return nil, fmt.Errorf("bytes are not representing a ActionResponseNormal. Action type is %d", actionType)
	}
	
	data = data[2:]
	
	// Parse invoke_id_and_priority
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for invoke_id_and_priority")
	}
	invokeIdAndPriority, err := (&InvokeIdAndPriority{}).FromBytes(data[:1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse invoke_id_and_priority: %w", err)
	}
	data = data[1:]
	
	// Parse status
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for status")
	}
	status := enumerations.ActionResultStatus(data[0])
	data = data[1:]
	
	// Parse has_data flag (should be 1 for response with data)
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for has_data flag")
	}
	hasData := data[0] != 0
	data = data[1:]
	
	if !hasData {
		return nil, fmt.Errorf("ActionResponseNormalWithData should have data")
	}
	
	// Parse data (remaining bytes)
	responseData := make([]byte, len(data))
	copy(responseData, data)
	
	return NewActionResponseNormalWithData(status, responseData, invokeIdAndPriority), nil
}

// ToBytes converts ActionResponseNormalWithData to bytes
func (a *ActionResponseNormalWithData) ToBytes() ([]byte, error) {
	result := []byte{ActionResponseTag}
	result = append(result, byte(enumerations.ActionTypeNormal))
	
	invokeBytes := a.InvokeIdAndPriority.ToBytes()
	result = append(result, invokeBytes...)
	
	result = append(result, byte(a.Status))
	result = append(result, 0x01) // has_data = true
	result = append(result, a.Data...)
	
	return result, nil
}

// ActionResponseNormalWithError represents an Action response normal with error
type ActionResponseNormalWithError struct {
	*BaseXDlmsApdu
	Status              enumerations.ActionResultStatus
	Error                enumerations.DataAccessResult
	InvokeIdAndPriority *InvokeIdAndPriority
}

// NewActionResponseNormalWithError creates a new ActionResponseNormalWithError
func NewActionResponseNormalWithError(
	status enumerations.ActionResultStatus,
	error enumerations.DataAccessResult,
	invokeIdAndPriority *InvokeIdAndPriority,
) *ActionResponseNormalWithError {
	return &ActionResponseNormalWithError{
		BaseXDlmsApdu: &BaseXDlmsApdu{
			Tag: ActionResponseTag,
		},
		Status:              status,
		Error:                error,
		InvokeIdAndPriority: invokeIdAndPriority,
	}
}

// FromBytes creates ActionResponseNormalWithError from bytes
func (a *ActionResponseNormalWithError) FromBytes(data []byte) (*ActionResponseNormalWithError, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("insufficient data for ActionResponse")
	}
	
	tag := data[0]
	if tag != ActionResponseTag {
		return nil, fmt.Errorf("tag %d is not correct for ActionResponse. Should be %d", tag, ActionResponseTag)
	}
	
	actionType := enumerations.ActionType(data[1])
	if actionType != enumerations.ActionTypeNormal {
		return nil, fmt.Errorf("bytes are not representing a ActionResponseNormal. Action type is %d", actionType)
	}
	
	data = data[2:]
	
	// Parse invoke_id_and_priority
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for invoke_id_and_priority")
	}
	invokeIdAndPriority, err := (&InvokeIdAndPriority{}).FromBytes(data[:1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse invoke_id_and_priority: %w", err)
	}
	data = data[1:]
	
	// Parse status
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for status")
	}
	status := enumerations.ActionResultStatus(data[0])
	data = data[1:]
	
	// Parse has_data flag (should be 1 for response with error)
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for has_data flag")
	}
	hasData := data[0] != 0
	data = data[1:]
	
	if !hasData {
		return nil, fmt.Errorf("ActionResponseNormalWithError should have data")
	}
	
	// Parse choice (should be 1 for error)
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for choice")
	}
	choice := data[0]
	if choice != 1 {
		return nil, fmt.Errorf("expected choice=1 for error, got %d", choice)
	}
	data = data[1:]
	
	// Parse error
	if len(data) < 1 {
		return nil, fmt.Errorf("insufficient data for error")
	}
	error := enumerations.DataAccessResult(data[0])
	
	return NewActionResponseNormalWithError(status, error, invokeIdAndPriority), nil
}

// ToBytes converts ActionResponseNormalWithError to bytes
func (a *ActionResponseNormalWithError) ToBytes() ([]byte, error) {
	result := []byte{ActionResponseTag}
	result = append(result, byte(enumerations.ActionTypeNormal))
	
	invokeBytes := a.InvokeIdAndPriority.ToBytes()
	result = append(result, invokeBytes...)
	
	result = append(result, byte(a.Status))
	result = append(result, 0x01) // has_data = true
	result = append(result, 0x01) // choice = 1 (error)
	result = append(result, byte(a.Error))
	
	return result, nil
}
