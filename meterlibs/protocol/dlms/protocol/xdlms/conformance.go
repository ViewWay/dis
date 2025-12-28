package xdlms

import (
	"encoding/binary"
	"fmt"
)

// Conformance holds information about the supported services in a DLMS association.
// Is used to send the proposed conformance in AARQ and to send back the negotiated
// conformance in the AARE.
// Only LN referencing is supported.
type Conformance struct {
	GeneralProtection            bool
	GeneralBlockTransfer         bool
	DeltaValueEncoding           bool
	Attribute0SupportedWithSet  bool
	PriorityManagementSupported bool
	Attribute0SupportedWithGet  bool
	BlockTransferWithGetOrRead  bool
	BlockTransferWithSetOrWrite bool
	BlockTransferWithAction     bool
	MultipleReferences          bool
	DataNotification            bool
	Access                      bool
	Get                         bool
	Set                         bool
	SelectiveAccess             bool
	EventNotification           bool
	Action                      bool
}

// ConformanceBitPosition maps attribute names to bit positions
// Bit numbering starts at 0, with bit 0 being LSB
var ConformanceBitPosition = map[string]int{
	"general_protection":             22,
	"general_block_transfer":          21,
	"delta_value_encoding":            17,
	"attribute_0_supported_with_set":  15,
	"priority_management_supported":  14,
	"attribute_0_supported_with_get":  13,
	"block_transfer_with_get_or_read": 12,
	"block_transfer_with_set_or_write": 11,
	"block_transfer_with_action":      10,
	"multiple_references":             9,
	"data_notification":               7,
	"access":                           6,
	"get":                              4,
	"set":                              3,
	"selective_access":                 2,
	"event_notification":               1,
	"action":                            0,
}

// NewConformance creates a new Conformance
func NewConformance(
	generalProtection, generalBlockTransfer, deltaValueEncoding,
	attribute0SupportedWithSet, priorityManagementSupported,
	attribute0SupportedWithGet, blockTransferWithGetOrRead,
	blockTransferWithSetOrWrite, blockTransferWithAction,
	multipleReferences, dataNotification, access, get, set,
	selectiveAccess, eventNotification, action bool,
) *Conformance {
	return &Conformance{
		GeneralProtection:            generalProtection,
		GeneralBlockTransfer:         generalBlockTransfer,
		DeltaValueEncoding:           deltaValueEncoding,
		Attribute0SupportedWithSet:   attribute0SupportedWithSet,
		PriorityManagementSupported:  priorityManagementSupported,
		Attribute0SupportedWithGet:   attribute0SupportedWithGet,
		BlockTransferWithGetOrRead:    blockTransferWithGetOrRead,
		BlockTransferWithSetOrWrite:   blockTransferWithSetOrWrite,
		BlockTransferWithAction:      blockTransferWithAction,
		MultipleReferences:            multipleReferences,
		DataNotification:             dataNotification,
		Access:                       access,
		Get:                          get,
		Set:                          set,
		SelectiveAccess:              selectiveAccess,
		EventNotification:            eventNotification,
		Action:                       action,
	}
}

// FromBytes creates Conformance from bytes
func (c *Conformance) FromBytes(data []byte) (*Conformance, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("insufficient data for Conformance: need at least 4 bytes, got %d", len(data))
	}
	
	// Skip first byte (unused bits indicator) and read 3 bytes
	integerRepresentation := binary.BigEndian.Uint32(append([]byte{0}, data[1:4]...))
	
	conf := &Conformance{}
	
	conf.GeneralProtection = (integerRepresentation & (1 << ConformanceBitPosition["general_protection"])) != 0
	conf.GeneralBlockTransfer = (integerRepresentation & (1 << ConformanceBitPosition["general_block_transfer"])) != 0
	conf.DeltaValueEncoding = (integerRepresentation & (1 << ConformanceBitPosition["delta_value_encoding"])) != 0
	conf.Attribute0SupportedWithSet = (integerRepresentation & (1 << ConformanceBitPosition["attribute_0_supported_with_set"])) != 0
	conf.PriorityManagementSupported = (integerRepresentation & (1 << ConformanceBitPosition["priority_management_supported"])) != 0
	conf.Attribute0SupportedWithGet = (integerRepresentation & (1 << ConformanceBitPosition["attribute_0_supported_with_get"])) != 0
	conf.BlockTransferWithGetOrRead = (integerRepresentation & (1 << ConformanceBitPosition["block_transfer_with_get_or_read"])) != 0
	conf.BlockTransferWithSetOrWrite = (integerRepresentation & (1 << ConformanceBitPosition["block_transfer_with_set_or_write"])) != 0
	conf.BlockTransferWithAction = (integerRepresentation & (1 << ConformanceBitPosition["block_transfer_with_action"])) != 0
	conf.MultipleReferences = (integerRepresentation & (1 << ConformanceBitPosition["multiple_references"])) != 0
	conf.DataNotification = (integerRepresentation & (1 << ConformanceBitPosition["data_notification"])) != 0
	conf.Access = (integerRepresentation & (1 << ConformanceBitPosition["access"])) != 0
	conf.Get = (integerRepresentation & (1 << ConformanceBitPosition["get"])) != 0
	conf.Set = (integerRepresentation & (1 << ConformanceBitPosition["set"])) != 0
	conf.SelectiveAccess = (integerRepresentation & (1 << ConformanceBitPosition["selective_access"])) != 0
	conf.EventNotification = (integerRepresentation & (1 << ConformanceBitPosition["event_notification"])) != 0
	conf.Action = (integerRepresentation & (1 << ConformanceBitPosition["action"])) != 0
	
	return conf, nil
}

// ToBytes converts Conformance to bytes
func (c *Conformance) ToBytes() []byte {
	var out uint32
	
	if c.GeneralProtection {
		out |= 1 << ConformanceBitPosition["general_protection"]
	}
	if c.GeneralBlockTransfer {
		out |= 1 << ConformanceBitPosition["general_block_transfer"]
	}
	if c.DeltaValueEncoding {
		out |= 1 << ConformanceBitPosition["delta_value_encoding"]
	}
	if c.Attribute0SupportedWithSet {
		out |= 1 << ConformanceBitPosition["attribute_0_supported_with_set"]
	}
	if c.PriorityManagementSupported {
		out |= 1 << ConformanceBitPosition["priority_management_supported"]
	}
	if c.Attribute0SupportedWithGet {
		out |= 1 << ConformanceBitPosition["attribute_0_supported_with_get"]
	}
	if c.BlockTransferWithGetOrRead {
		out |= 1 << ConformanceBitPosition["block_transfer_with_get_or_read"]
	}
	if c.BlockTransferWithSetOrWrite {
		out |= 1 << ConformanceBitPosition["block_transfer_with_set_or_write"]
	}
	if c.BlockTransferWithAction {
		out |= 1 << ConformanceBitPosition["block_transfer_with_action"]
	}
	if c.MultipleReferences {
		out |= 1 << ConformanceBitPosition["multiple_references"]
	}
	if c.DataNotification {
		out |= 1 << ConformanceBitPosition["data_notification"]
	}
	if c.Access {
		out |= 1 << ConformanceBitPosition["access"]
	}
	if c.Get {
		out |= 1 << ConformanceBitPosition["get"]
	}
	if c.Set {
		out |= 1 << ConformanceBitPosition["set"]
	}
	if c.SelectiveAccess {
		out |= 1 << ConformanceBitPosition["selective_access"]
	}
	if c.EventNotification {
		out |= 1 << ConformanceBitPosition["event_notification"]
	}
	if c.Action {
		out |= 1 << ConformanceBitPosition["action"]
	}
	
	// It is a bit string so need to encode how many bits that are unused in the
	// last byte. It's none so we can just put 0x00 in front.
	result := make([]byte, 4)
	result[0] = 0x00 // unused bits indicator
	binary.BigEndian.PutUint32(result[1:], out)
	// Only use 3 bytes for the bit string
	return result[:4]
}

