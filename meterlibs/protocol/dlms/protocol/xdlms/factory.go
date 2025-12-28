package xdlms

import (
	"fmt"

	"github.com/yimiliya/idis/meterlibs/protocol/dlms/protocol/acse"
)

// XDlmsApduFactory is a factory to return the correct APDU depending on the tag
type XDlmsApduFactory struct{}

// APDUFromBytes parses an APDU from bytes based on its tag
func (f *XDlmsApduFactory) APDUFromBytes(apduBytes []byte) (interface{}, error) {
	if len(apduBytes) == 0 {
		return nil, fmt.Errorf("insufficient data for APDU tag")
	}

	tag := apduBytes[0]

	switch tag {
	// xDLMS APDUs
	case 1:
		initReq := &InitiateRequest{}
		return initReq.FromBytes(apduBytes)
	case 8:
		initResp := &InitiateResponse{}
		return initResp.FromBytes(apduBytes)
	case 14:
		// ConfirmedServiceError - TODO: implement when needed
		return nil, fmt.Errorf("ConfirmedServiceError not yet implemented")
	case 15:
		dataNotif := &DataNotification{}
		return dataNotif.FromBytes(apduBytes)
	case 33:
		// GlobalCipherInitiateRequest - TODO: implement when needed
		return nil, fmt.Errorf("GlobalCipherInitiateRequest not yet implemented")
	case 40:
		// GlobalCipherInitiateResponse - TODO: implement when needed
		return nil, fmt.Errorf("GlobalCipherInitiateResponse not yet implemented")
	case 216:
		excResp := &ExceptionResponse{}
		return excResp.FromBytes(apduBytes)
	case 219:
		// GeneralGlobalCipher - TODO: implement when needed
		return nil, fmt.Errorf("GeneralGlobalCipher not yet implemented")
	// ACSE APDUs
	case 96:
		aarq := &acse.ApplicationAssociationRequest{}
		return aarq.FromBytes(apduBytes)
	case 97:
		aare := &acse.ApplicationAssociationResponse{}
		return aare.FromBytes(apduBytes)
	case 98:
		rlrq := &acse.ReleaseRequest{}
		return rlrq.FromBytes(apduBytes)
	case 99:
		rlre := &acse.ReleaseResponse{}
		return rlre.FromBytes(apduBytes)
	// GET requests/responses (use factories)
	case 192:
		return GetRequestFromBytes(apduBytes)
	case 196:
		return GetResponseFromBytes(apduBytes)
	// SET requests/responses (use factories)
	case 193:
		return SetRequestFromBytes(apduBytes)
	case 197:
		return SetResponseFromBytes(apduBytes)
	// ACTION requests/responses (use factories)
	case 195:
		return ActionRequestFromBytes(apduBytes)
	case 199:
		return ActionResponseFromBytes(apduBytes)
	default:
		return nil, fmt.Errorf("tag 0x%02x is not available in DLMS APDU Factory", tag)
	}
}

// GetRequestFromBytes parses a GetRequest from bytes
func GetRequestFromBytes(sourceBytes []byte) (interface{}, error) {
	if len(sourceBytes) < 2 {
		return nil, fmt.Errorf("insufficient data for GetRequest")
	}

	tag := sourceBytes[0]
	if tag != GetRequestTag {
		return nil, fmt.Errorf("tag for GET request is not correct, got %d, should be %d", tag, GetRequestTag)
	}

	requestType := sourceBytes[1]
	switch requestType {
	case 1: // GetRequestNormal
		req := &GetRequestNormal{}
		return req.FromBytes(sourceBytes)
	case 2: // GetRequestNext
		req := &GetRequestNext{}
		return req.FromBytes(sourceBytes)
	case 3: // GetRequestWithList
		req := &GetRequestWithList{}
		return req.FromBytes(sourceBytes)
	default:
		return nil, fmt.Errorf("received an enum request type that is not valid for GetRequest: %d", requestType)
	}
}

// GetResponseFromBytes parses a GetResponse from bytes
func GetResponseFromBytes(sourceBytes []byte) (interface{}, error) {
	if len(sourceBytes) < 2 {
		return nil, fmt.Errorf("insufficient data for GetResponse")
	}

	tag := sourceBytes[0]
	if tag != 196 {
		return nil, fmt.Errorf("tag for GET response is not correct, got %d, should be 196", tag)
	}

	responseType := sourceBytes[1]
	switch responseType {
	case 1: // GetResponseNormal
		// Check if it's an error response by looking at the choice field
		// Format: [tag, type, invoke_id_and_priority(1 byte), choice(1 byte), ...]
		if len(sourceBytes) >= 4 {
			choice := sourceBytes[3]
			if choice == 1 {
				// GetResponseNormalWithError
				respWithError := &GetResponseNormalWithError{}
				return respWithError.FromBytes(sourceBytes)
			}
		}
		// GetResponseNormal
		resp := &GetResponseNormal{}
		return resp.FromBytes(sourceBytes)
	case 2: // GetResponseWithDataBlock
		resp := &GetResponseWithDataBlock{}
		return resp.FromBytes(sourceBytes)
	case 3: // GetResponseWithList
		resp := &GetResponseWithList{}
		return resp.FromBytes(sourceBytes)
	case 4: // GetResponseLastBlock
		resp := &GetResponseLastBlock{}
		return resp.FromBytes(sourceBytes)
	case 5: // GetResponseLastBlockWithError
		resp := &GetResponseLastBlockWithError{}
		return resp.FromBytes(sourceBytes)
	default:
		return nil, fmt.Errorf("received an enum response type that is not valid for GetResponse: %d", responseType)
	}
}

// SetRequestFromBytes parses a SetRequest from bytes
func SetRequestFromBytes(sourceBytes []byte) (interface{}, error) {
	if len(sourceBytes) < 2 {
		return nil, fmt.Errorf("insufficient data for SetRequest")
	}

	tag := sourceBytes[0]
	if tag != 193 {
		return nil, fmt.Errorf("tag for SET request is not correct, got %d, should be 193", tag)
	}

	requestType := sourceBytes[1]
	switch requestType {
	case 1: // SetRequestNormal
		req := &SetRequestNormal{}
		return req.FromBytes(sourceBytes)
	default:
		return nil, fmt.Errorf("received an enum request type that is not valid for SetRequest: %d", requestType)
	}
}

// SetResponseFromBytes parses a SetResponse from bytes
func SetResponseFromBytes(sourceBytes []byte) (interface{}, error) {
	if len(sourceBytes) < 2 {
		return nil, fmt.Errorf("insufficient data for SetResponse")
	}

	tag := sourceBytes[0]
	if tag != 197 {
		return nil, fmt.Errorf("tag for SET response is not correct, got %d, should be 197", tag)
	}

	responseType := sourceBytes[1]
	switch responseType {
	case 1: // SetResponseNormal
		resp := &SetResponseNormal{}
		return resp.FromBytes(sourceBytes)
	default:
		return nil, fmt.Errorf("received an enum response type that is not valid for SetResponse: %d", responseType)
	}
}

// ActionRequestFromBytes parses an ActionRequest from bytes
func ActionRequestFromBytes(sourceBytes []byte) (interface{}, error) {
	if len(sourceBytes) < 2 {
		return nil, fmt.Errorf("insufficient data for ActionRequest")
	}

	tag := sourceBytes[0]
	if tag != 195 {
		return nil, fmt.Errorf("tag for ACTION request is not correct, got %d, should be 195", tag)
	}

	requestType := sourceBytes[1]
	switch requestType {
	case 1: // ActionRequestNormal
		req := &ActionRequestNormal{}
		return req.FromBytes(sourceBytes)
	default:
		return nil, fmt.Errorf("received an enum request type that is not valid for ActionRequest: %d", requestType)
	}
}

// ActionResponseFromBytes parses an ActionResponse from bytes
func ActionResponseFromBytes(sourceBytes []byte) (interface{}, error) {
	if len(sourceBytes) < 4 {
		return nil, fmt.Errorf("insufficient data for ActionResponse")
	}

	tag := sourceBytes[0]
	if tag != 199 {
		return nil, fmt.Errorf("tag for ACTION response is not correct, got %d, should be 199", tag)
	}

	responseType := sourceBytes[1]
	if responseType != 1 {
		return nil, fmt.Errorf("received an enum response type that is not valid for ActionResponse: %d", responseType)
	}

	// Check if it has data by looking at the has_data flag (after invoke_id_and_priority and status)
	// Format: [tag, type, invoke_id_and_priority(1 byte), status(1 byte), has_data(1 byte), ...]
	if len(sourceBytes) >= 5 {
		hasData := sourceBytes[4] != 0
		if hasData {
			// Check the choice field to determine if it's data or error
			if len(sourceBytes) >= 6 {
				choice := sourceBytes[5]
				if choice == 0 {
					// ActionResponseNormalWithData
					respWithData := &ActionResponseNormalWithData{}
					return respWithData.FromBytes(sourceBytes)
				} else if choice == 1 {
					// ActionResponseNormalWithError
					respWithError := &ActionResponseNormalWithError{}
					return respWithError.FromBytes(sourceBytes)
				}
			}
		} else {
			// ActionResponseNormal (no data)
			resp := &ActionResponseNormal{}
			return resp.FromBytes(sourceBytes)
		}
	}

	// Fallback to ActionResponseNormal
	resp := &ActionResponseNormal{}
	return resp.FromBytes(sourceBytes)
}

// NewXDlmsApduFactory creates a new XDlmsApduFactory
func NewXDlmsApduFactory() *XDlmsApduFactory {
	return &XDlmsApduFactory{}
}
